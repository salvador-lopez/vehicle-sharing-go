package rest_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"vehicle-sharing-go/internal/inventory/vehicle/application/projection"
	"vehicle-sharing-go/internal/inventory/vehicle/infrastructure/controller/gen/car"
	"vehicle-sharing-go/internal/inventory/vehicle/infrastructure/controller/rest"
	"vehicle-sharing-go/internal/inventory/vehicle/infrastructure/controller/rest/mock"
)

type carUnitSuite struct {
	suite.Suite
	ctx                   context.Context
	mockCtrl              *gomock.Controller
	mockCreateCarCHandler *mock.MockCreateCarCommandHandler
	mockCarQueryService   *mock.MockCarQueryService
	sut                   *rest.CarController
}

func (s *carUnitSuite) SetupTest() {
	s.ctx = context.Background()

	s.mockCtrl = gomock.NewController(s.T())
	s.mockCarQueryService = mock.NewMockCarQueryService(s.mockCtrl)

	s.sut = rest.NewCarController(s.mockCreateCarCHandler, s.mockCarQueryService)
}

func (s *carUnitSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func TestCarUnitSuite(t *testing.T) {
	suite.Run(t, new(carUnitSuite))
}

func (s *carUnitSuite) TestGet() {
	carID := uuid.New()

	projectionVinData := s.buildVinDataProjection(
		"4Y1SL65848Z411439",
		"country",
		"manufacturer",
		"brand",
		"2000",
		"Gasoline",
		"Jazz",
		"2023",
		"Barcelona",
		"411439",
	)

	carProjection := &projection.Car{
		ID:        carID,
		CreatedAt: time.Now().Add(-time.Hour),
		UpdatedAt: time.Now(),
		VINData:   projectionVinData,
		Color:     "Spectral Blue",
	}

	s.mockCarQueryService.EXPECT().Find(s.ctx, carID).Return(carProjection, nil)

	carResource, err := s.sut.Get(s.ctx, &car.GetPayload{ID: carID.String()})
	s.Require().NoError(err)

	expectedCarResource := &car.CarResource{
		ID:        carID.String(),
		CreatedAt: carProjection.CreatedAt.String(),
		UpdatedAt: carProjection.UpdatedAt.String(),
		Color:     carProjection.Color,
		VinData: &car.VinData{
			Vin:           car.Vin(carProjection.VIN),
			Country:       projectionVinData.Country,
			Manufacturer:  projectionVinData.Manufacturer,
			Brand:         projectionVinData.Brand,
			EngineSize:    projectionVinData.EngineSize,
			FuelType:      projectionVinData.FuelType,
			Model:         projectionVinData.Model,
			Year:          projectionVinData.Year,
			AssemblyPlant: projectionVinData.AssemblyPlant,
			SN:            projectionVinData.SN,
		},
	}
	s.Require().Equal(expectedCarResource, carResource)
}

func (s *carUnitSuite) buildVinDataProjection(
	vinNumber,
	country,
	manufacturer,
	brand,
	engineSize,
	fuelType,
	model,
	year,
	assemblyPlant,
	sn string,
) *projection.VINData {
	return &projection.VINData{
		VIN:           vinNumber,
		Country:       &country,
		Manufacturer:  &manufacturer,
		Brand:         &brand,
		EngineSize:    &engineSize,
		FuelType:      &fuelType,
		Model:         &model,
		Year:          &year,
		AssemblyPlant: &assemblyPlant,
		SN:            &sn,
	}
}
