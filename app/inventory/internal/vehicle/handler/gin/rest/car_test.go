//go:build unit

package rest_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"vehicle-sharing-go/app/inventory/internal/vehicle/command"
	"vehicle-sharing-go/app/inventory/internal/vehicle/handler/gin/rest"
	"vehicle-sharing-go/app/inventory/internal/vehicle/handler/gin/rest/mock"
	"vehicle-sharing-go/app/inventory/internal/vehicle/projection"
	"vehicle-sharing-go/pkg/domain"
)

type carUnitSuite struct {
	suite.Suite
	rr                    *httptest.ResponseRecorder
	c                     *gin.Context
	mockCtrl              *gomock.Controller
	mockCreateCarCHandler *mock.MockCreateCarCommandHandler
	mockCarQueryService   *mock.MockFindCarQueryService
	sut                   *rest.CarHandler
}

func (s *carUnitSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (s *carUnitSuite) SetupTest() {
	s.rr = httptest.NewRecorder()
	c, _ := gin.CreateTestContext(s.rr)
	s.c = c
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
			s.mockCarQueryService.EXPECT().Find(s.c, tt.carID).Return(expectedProjection, nil)

			s.c.Request = s.getCarReq()
			s.c.AddParam("id", tt.carID.String())

			s.sut.Get(s.c)

			s.Require().Equal(http.StatusOK, s.rr.Code)

			var actualProjection projection.Car
			err := json.NewDecoder(s.rr.Body).Decode(&actualProjection)
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
			sutErrMsg: "{\"error\":\"not found: b1e3580a-acd5-4081-9d2c-74366a580f36\"}",
		},
		{
			name:            "Query Service Error",
			carID:           "2279e813-d3ec-4be4-9c41-02315873fc34",
			code:            http.StatusInternalServerError,
			queryServiceErr: errors.New("queryService error"),
			sutErrMsg:       "{\"error\":\"internal error\"}",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			defer s.TearDownTest()

			carID, err := uuid.Parse(tt.carID)
			s.Require().NoError(err)

			s.mockCarQueryService.EXPECT().Find(s.c, carID).Return(nil, tt.queryServiceErr)

			s.c.Request = s.getCarReq()
			s.c.AddParam("id", tt.carID)

			s.sut.Get(s.c)

			s.Require().Equal(tt.code, s.rr.Code)
			s.Require().Equal(tt.sutErrMsg, s.rr.Body.String())
		})
	}

	s.Run("Invalid Car ID provided in path", func() {
		s.SetupTest()
		defer s.TearDownTest()

		req := httptest.NewRequest(http.MethodGet, "/i-dont-care-about-the-endpoint", nil)
		s.c.Request = req
		s.c.AddParam("id", "invalid-car-id")

		s.sut.Get(s.c)

		s.Require().Equal(http.StatusBadRequest, s.rr.Code)
		s.Require().Equal("{\"error\":\"bad request: invalid UUID length: 14\"}", s.rr.Body.String())
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
			code:      http.StatusCreated,
		},
		{
			name:         "Domain conflict 409 response",
			carID:        uuid.NewString(),
			vinNumber:    "4Z1SL65848Z411440",
			color:        "Wolf Gray",
			cHandlerErr:  domain.WrapErrConflict(errors.New("chandler domain conflict err")),
			code:         http.StatusConflict,
			responseBody: "{\"error\":\"domain conflict: chandler domain conflict err\"}",
		},
		{
			name:         "Internal server error response 500",
			carID:        uuid.NewString(),
			vinNumber:    "4Z1SL65848Z411440",
			color:        "Red Storm",
			cHandlerErr:  errors.New("command handler err"),
			code:         http.StatusInternalServerError,
			responseBody: "{\"error\":\"internal error\"}",
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
			s.mockCreateCarCHandler.EXPECT().Handle(s.c, cmd).Return(tt.cHandlerErr)

			jsonReqBody, err := json.Marshal(cmd)
			s.Require().NoError(err)
			s.c.Request = s.createCarReq(jsonReqBody)
			s.sut.Create(s.c)

			s.Require().Equal(tt.code, s.rr.Code)
			s.Require().Equal(tt.responseBody, s.rr.Body.String())
		})
	}

	s.Run("Invalid Request Body param (id is not an UUID)", func() {
		s.SetupTest()
		defer s.TearDownTest()

		type invalidReqBody struct {
			ID    int    `json:"id"`
			VIN   string `json:"vin"`
			Color string `json:"color"`
		}

		reqBody := invalidReqBody{ID: 27, VIN: "4Z1SL65848Z411440", Color: "Spectral Blue Portimao"}
		jsonReqBody, err := json.Marshal(reqBody)
		s.Require().NoError(err)
		s.c.Request = s.createCarReq(jsonReqBody)
		s.sut.Create(s.c)

		s.Require().Equal(http.StatusBadRequest, s.rr.Code)
		s.Require().Equal(
			"{\"error\":\"bad request: json: cannot unmarshal number into Go struct field CreateCar.id of type uuid.UUID\"}",
			s.rr.Body.String(),
		)
	})
}

func (s *carUnitSuite) getCarReq() *http.Request {
	return httptest.NewRequest(http.MethodGet, "/i-dont-care-about-the-endpoint", nil)
}

func (s *carUnitSuite) createCarReq(jsonReqBody []byte) *http.Request {
	return httptest.NewRequest(http.MethodPost, "/i-dont-care-about-the-endpoint", bytes.NewReader(jsonReqBody))
}
