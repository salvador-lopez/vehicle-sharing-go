//go:build unit

package rest_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"vehicle-sharing-go/app/inventory/internal/vehicle/handler/http-net/rest"
	"vehicle-sharing-go/app/inventory/internal/vehicle/handler/http-net/rest/mock"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"vehicle-sharing-go/app/inventory/internal/vehicle/projection"
)

type carUnitSuite struct {
	suite.Suite
	ctx                 context.Context
	mockCtrl            *gomock.Controller
	mockCarQueryService *mock.MockFindCarQueryService
	sut                 *rest.CarHandler
}

func (s *carUnitSuite) SetupTest() {
	s.ctx = context.Background()

	s.mockCtrl = gomock.NewController(s.T())
	s.mockCarQueryService = mock.NewMockFindCarQueryService(s.mockCtrl)

	s.sut = rest.NewCarHandler(s.mockCarQueryService)
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

			req := httptest.NewRequest(http.MethodGet, "/cars/"+tt.carID.String(), nil)
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
		carID           uuid.UUID
		code            int
		queryServiceErr error
		sutErrMsg       string
	}{
		{
			name:            "Car Not Found Err",
			code:            http.StatusNotFound,
			queryServiceErr: nil,
			sutErrMsg:       "not found\n",
		},
		{
			name:            "Query Service Error",
			code:            http.StatusInternalServerError,
			queryServiceErr: errors.New("queryService error"),
			sutErrMsg:       "internal error\n",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			defer s.TearDownTest()

			s.mockCarQueryService.EXPECT().Find(s.ctx, tt.carID).Return(nil, tt.queryServiceErr)

			req := httptest.NewRequest(http.MethodGet, "/cars/"+tt.carID.String(), nil)
			rr := httptest.NewRecorder()

			s.sut.Get(s.ctx, rr, req)

			s.Require().Equal(tt.code, rr.Code)
			s.Require().Equal(tt.sutErrMsg, rr.Body.String())
		})
	}
}
