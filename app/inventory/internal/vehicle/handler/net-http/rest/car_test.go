//go:build unit

package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"vehicle-sharing-go/app/inventory/internal/vehicle/command"
	"vehicle-sharing-go/app/inventory/internal/vehicle/handler/net-http/rest"
	"vehicle-sharing-go/app/inventory/internal/vehicle/handler/net-http/rest/mock"
	"vehicle-sharing-go/app/inventory/internal/vehicle/projection"
	"vehicle-sharing-go/pkg/domain"
)

type carUnitSuite struct {
	suite.Suite
	ctx                   context.Context
	mockCtrl              *gomock.Controller
	mockCreateCarCHandler *mock.MockCreateCarCommandHandler
	mockCarQueryService   *mock.MockFindCarQueryService
	sut                   *rest.CarHandler
}

func (s *carUnitSuite) SetupTest() {
	s.ctx = context.Background()

	s.mockCtrl = gomock.NewController(s.T())
	s.mockCreateCarCHandler = mock.NewMockCreateCarCommandHandler(s.mockCtrl)
	s.mockCarQueryService = mock.NewMockFindCarQueryService(s.mockCtrl)

	s.sut = rest.NewCarHandler(s.mockCreateCarCHandler, s.mockCarQueryService)
}

func (s *carUnitSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func TestCarUnitSuite(t *testing.T) {
	suite.Run(t, new(carUnitSuite))
}

func (s *carUnitSuite) TestGetNoErr() {
	tests := []struct {
		name          string
		carID         uuid.UUID
		carCreatedAt  time.Time
		carUpdatedAt  time.Time
		vinNumber     string
		country       string
		manufacturer  string
		brand         string
		engineSize    string
		fuelType      string
		model         string
		year          string
		assemblyPlant string
		sn            string
		color         string
	}{
		{
			name:          "Return car resource with all optional data set",
			carID:         uuid.New(),
			carCreatedAt:  time.Now().Add(-time.Hour),
			carUpdatedAt:  time.Now(),
			vinNumber:     "4Y1SL65848Z411439",
			country:       "country",
			manufacturer:  "manufacturer",
			brand:         "brand",
			engineSize:    "2000",
			fuelType:      "Gasoline",
			model:         "Jazz",
			year:          "2023",
			assemblyPlant: "Barcelona",
			sn:            "411439",
			color:         "Spectral Blue",
		},
		{
			name:         "Return car resource with only mandatory data set",
			carID:        uuid.New(),
			carCreatedAt: time.Now().Add(-time.Minute * 2),
			carUpdatedAt: time.Now().Add(-time.Minute),
			vinNumber:    "4Z1SL65848Z411440",
			color:        "Black Bullet",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			defer s.TearDownTest()

			expectedProjection := &projection.Car{
				ID:        tt.carID,
				CreatedAt: tt.carCreatedAt.UTC(),
				UpdatedAt: tt.carUpdatedAt.UTC(),
				VINData: &projection.VINData{
					VIN:           tt.vinNumber,
					Country:       &tt.country,
					Manufacturer:  &tt.manufacturer,
					Brand:         &tt.brand,
					EngineSize:    &tt.engineSize,
					FuelType:      &tt.fuelType,
					Model:         &tt.model,
					Year:          &tt.year,
					AssemblyPlant: &tt.assemblyPlant,
					SN:            &tt.sn,
				},
				Color: tt.color,
			}
			s.mockCarQueryService.EXPECT().Find(s.ctx, tt.carID).Return(expectedProjection, nil)

			req := httptest.NewRequest(http.MethodGet, "/cars/{id}", nil)
			req.SetPathValue("id", tt.carID.String())
			rr := httptest.NewRecorder()

			s.sut.Get(s.ctx, rr, req)

			s.Require().Equal(http.StatusOK, rr.Code)

			var actualProjection projection.Car
			err := json.NewDecoder(rr.Body).Decode(&actualProjection)
			s.Require().NoError(err)
			s.Require().Equal(expectedProjection, &actualProjection)
		})
	}
}

func (s *carUnitSuite) TestGetErr() {
	tests := []struct {
		name            string
		carID           string
		code            int
		queryServiceErr error
		sutErrMsg       string
	}{
		{
			name:      "Car Not Found Err",
			carID:     "b1e3580a-acd5-4081-9d2c-74366a580f36",
			code:      http.StatusNotFound,
			sutErrMsg: "{\"error\":\"not found: b1e3580a-acd5-4081-9d2c-74366a580f36\"}\n",
		},
		{
			name:            "Query Service Error",
			carID:           "2279e813-d3ec-4be4-9c41-02315873fc34",
			code:            http.StatusInternalServerError,
			queryServiceErr: errors.New("queryService error"),
			sutErrMsg:       "{\"error\":\"internal error\"}\n",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			defer s.TearDownTest()

			carID, err := uuid.Parse(tt.carID)
			s.Require().NoError(err)

			s.mockCarQueryService.EXPECT().Find(s.ctx, carID).Return(nil, tt.queryServiceErr)

			req := httptest.NewRequest(http.MethodGet, "/cars/{id}", nil)
			req.SetPathValue("id", carID.String())
			rr := httptest.NewRecorder()

			s.sut.Get(s.ctx, rr, req)

			s.Require().Equal(tt.code, rr.Code)
			s.Require().Equal(tt.sutErrMsg, rr.Body.String())
		})
	}

	s.Run("Invalid Car ID provided in path", func() {
		s.SetupTest()
		defer s.TearDownTest()
		req := httptest.NewRequest(http.MethodGet, "/cars/{id}", nil)
		req.SetPathValue("id", "invalid-car-id")
		rr := httptest.NewRecorder()

		s.sut.Get(s.ctx, rr, req)

		s.Require().Equal(http.StatusBadRequest, rr.Code)
		s.Require().Equal("{\"error\":\"bad request: invalid UUID length: 14\"}\n", rr.Body.String())
	})
}

func (s *carUnitSuite) TestCreate() {
	tests := []struct {
		name         string
		carID        string
		vinNumber    string
		color        string
		cHandlerErr  error
		code         int
		responseBody string
	}{
		{
			name:      "Created with no error",
			carID:     uuid.NewString(),
			vinNumber: "4Y1SL65848Z411439",
			color:     "Spectral Blue",
			code:      http.StatusOK,
		},
		{
			name:         "Domain conflict 409 response",
			carID:        uuid.NewString(),
			vinNumber:    "4Z1SL65848Z411440",
			color:        "Wolf Gray",
			cHandlerErr:  domain.WrapErrConflict(errors.New("chandler domain conflict err")),
			code:         http.StatusConflict,
			responseBody: "{\"error\":\"domain conflict: chandler domain conflict err\"}\n",
		},
		{
			name:         "Internal server error response 500",
			carID:        uuid.NewString(),
			vinNumber:    "4Z1SL65848Z411440",
			color:        "Red Storm",
			cHandlerErr:  errors.New("command handler err"),
			code:         http.StatusInternalServerError,
			responseBody: "{\"error\":\"internal error\"}\n",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			defer s.TearDownTest()

			carID, err := uuid.Parse(tt.carID)
			s.Require().NoError(err)

			cmd := &command.CreateCar{
				VIN:   tt.vinNumber,
				ID:    carID,
				Color: tt.color,
			}
			s.mockCreateCarCHandler.EXPECT().Handle(s.ctx, cmd).Return(tt.cHandlerErr)

			reqBody := &command.CreateCar{ID: carID, VIN: tt.vinNumber, Color: tt.color}
			jsonReqBody, err := json.Marshal(reqBody)
			s.Require().NoError(err)
			req := httptest.NewRequest(http.MethodPost, "/cars", bytes.NewReader(jsonReqBody))
			rr := httptest.NewRecorder()
			s.sut.Create(s.ctx, rr, req)

			s.Require().Equal(tt.code, rr.Code)
			s.Require().Equal(tt.responseBody, rr.Body.String())
		})
	}

	s.Run("Invalid Request Body param (id is not an UUID)", func() {
		s.SetupTest()
		defer s.TearDownTest()

		type invalidReqBody struct {
			ID    int   `json:"id"`
			VIN   string `json:"vin"`
			Color string `json:"color"`
		}

		reqBody := invalidReqBody{ID: 27, VIN: "4Z1SL65848Z411440", Color: "Spectral Blue Portimao"}
		jsonReqBody, err := json.Marshal(reqBody)
		s.Require().NoError(err)
		req := httptest.NewRequest(http.MethodPost, "/cars", bytes.NewReader(jsonReqBody))
		rr := httptest.NewRecorder()
		s.sut.Create(s.ctx, rr, req)

		s.Require().Equal(http.StatusBadRequest, rr.Code)
		s.Require().Equal(
			"{\"error\":\"bad request: json: cannot unmarshal number into Go struct field CreateCar.id of type uuid.UUID\"}\n",
			rr.Body.String(),
		)
	})
}
