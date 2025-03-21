//go:build unit

package rest_test

import (
	"context"
	"errors"
	"testing"
	"time"
	"vehicle-sharing-go/pkg/domain"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	goa "goa.design/goa/v3/pkg"

	"vehicle-sharing-go/app/inventory/internal/vehicle/command"
	"vehicle-sharing-go/app/inventory/internal/vehicle/handler/rest"
	"vehicle-sharing-go/app/inventory/internal/vehicle/handler/rest/gen/car"
	"vehicle-sharing-go/app/inventory/internal/vehicle/handler/rest/mock"
	"vehicle-sharing-go/app/inventory/internal/vehicle/projection"
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
	s.mockCarQueryService = mock.NewMockFindCarQueryService(s.mockCtrl)
	s.mockCreateCarCHandler = mock.NewMockCreateCarCommandHandler(s.mockCtrl)

	s.sut = rest.NewCarHandler(s.mockCreateCarCHandler, s.mockCarQueryService)
}

func (s *carUnitSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func TestCarUnitSuite(t *testing.T) {
	suite.Run(t, new(carUnitSuite))
}

func (s *carUnitSuite) TestGet() {
	tests := []struct {
		name               string
		carID              uuid.UUID
		carProjectionFound bool
		carCreatedAt       time.Time
		carUpdatedAt       time.Time
		vinNumber          string
		country            string
		manufacturer       string
		brand              string
		engineSize         string
		fuelType           string
		model              string
		year               string
		assemblyPlant      string
		sn                 string
		color              string
		querySvcErr        error
		sutErrMsg          string
		goaErrName         string
	}{
		{
			name:               "Return car resource with all optional data set",
			carID:              uuid.New(),
			carProjectionFound: true,
			carCreatedAt:       time.Now().Add(-time.Hour),
			carUpdatedAt:       time.Now(),
			vinNumber:          "4Y1SL65848Z411439",
			country:            "country",
			manufacturer:       "manufacturer",
			brand:              "brand",
			engineSize:         "2000",
			fuelType:           "Gasoline",
			model:              "Jazz",
			year:               "2023",
			assemblyPlant:      "Barcelona",
			sn:                 "411439",
			color:              "Spectral Blue",
		},
		{
			name:               "Return car resource with only mandatory data set",
			carID:              uuid.New(),
			carProjectionFound: true,
			carCreatedAt:       time.Now().Add(-time.Minute * 2),
			carUpdatedAt:       time.Now().Add(-time.Minute),
			vinNumber:          "4Z1SL65848Z411440",
			color:              "Black Bullet",
		},
		{
			name:        "Return internal goa.ServiceError when query service Find() return error",
			carID:       uuid.New(),
			querySvcErr: errors.New("query service Find() err"),
			sutErrMsg:   rest.ErrInternal.Error(),
			goaErrName:  "internal",
		},
		{
			name:       "Return notFound goa.ServiceError when query service Find() return nil projection",
			carID:      uuid.New(),
			sutErrMsg:  rest.ErrNotFound.Error(),
			goaErrName: "notFound",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			defer s.TearDownTest()

			var carProjection *projection.Car

			if tt.carProjectionFound {
				carProjection = &projection.Car{
					ID:        tt.carID,
					CreatedAt: tt.carCreatedAt,
					UpdatedAt: tt.carUpdatedAt,
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
			}
			s.mockCarQueryService.EXPECT().Find(s.ctx, tt.carID).Return(carProjection, tt.querySvcErr)

			carResource, err := s.sut.Get(s.ctx, &car.GetPayload{ID: tt.carID.String()})

			var expectedCarResource *car.CarResource

			if tt.sutErrMsg == "" {
				s.Require().NoError(err)
				expectedCarResource = &car.CarResource{
					ID:        tt.carID.String(),
					CreatedAt: carProjection.CreatedAt.String(),
					UpdatedAt: carProjection.UpdatedAt.String(),
					Color:     carProjection.Color,
					VinData: &car.VinData{
						Vin:           car.Vin(carProjection.VIN),
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
				}
				s.Require().Equal(expectedCarResource, carResource)

				return
			}
			sErr, ok := err.(*goa.ServiceError)
			s.Require().True(ok)
			s.Require().Equal(tt.goaErrName, sErr.Name)
			s.Require().EqualError(sErr, tt.sutErrMsg)
			s.Require().Nil(carResource)
		})
	}
}

func (s *carUnitSuite) TestCreate() {
	tests := []struct {
		name        string
		carID       uuid.UUID
		vinNumber   string
		color       string
		cHandlerErr error
		sutErr      error
		goaErrName  string
	}{
		{
			name:      "Created with no error",
			carID:     uuid.New(),
			vinNumber: "4Y1SL65848Z411439",
			color:     "Spectral Blue",
		},
		{
			name:        "Return conflict goa.ServiceError when command handler return domain conflict error",
			carID:       uuid.New(),
			vinNumber:   "4Z1SL65848Z411440",
			color:       "Black Bullet",
			cHandlerErr: domain.ErrConflict,
			sutErr:      domain.ErrConflict,
			goaErrName:  "conflict",
		},
		{
			name:        "Return internal goa.ServiceError when command handler return error",
			carID:       uuid.New(),
			vinNumber:   "4Z1SL65848Z411440",
			color:       "Black Bullet",
			cHandlerErr: errors.New("command handler err"),
			sutErr:      rest.ErrInternal,
			goaErrName:  "internal",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			defer s.TearDownTest()

			cmd := &command.CreateCar{
				VIN:   tt.vinNumber,
				ID:    tt.carID,
				Color: tt.color,
			}
			s.mockCreateCarCHandler.EXPECT().Handle(s.ctx, cmd).Return(tt.cHandlerErr)

			payload := &car.CreatePayload{ID: tt.carID.String(), Vin: car.Vin(tt.vinNumber), Color: tt.color}
			err := s.sut.Create(s.ctx, payload)

			if tt.sutErr == nil {
				s.Require().NoError(err)
				return
			}
			var sErr *goa.ServiceError
			ok := errors.As(err, &sErr)
			s.Require().True(ok)
			s.Require().Equal(tt.goaErrName, sErr.Name)
			s.Require().EqualError(sErr, tt.sutErr.Error())
		})
	}
}
