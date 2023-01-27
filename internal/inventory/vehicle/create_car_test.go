package vehicle_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"vehicle-sharing-go/internal/inventory/vehicle"
	"vehicle-sharing-go/internal/inventory/vehicle/mock"
)

type vehicleUnitSuite struct {
	suite.Suite
	ctx         context.Context
	mockCtrl    *gomock.Controller
	now         time.Time
	mockCarRepo *mock.MockCarRepository
	sut         *vehicle.CreateCarHandler
}

func (s *vehicleUnitSuite) SetupTest() {
	s.ctx = context.Background()
	s.now = time.Now()
	s.mockCtrl = gomock.NewController(s.T())
	s.mockCarRepo = mock.NewMockCarRepository(s.mockCtrl)
	s.sut = vehicle.NewCreateCarHandler(func() time.Time { return s.now }, s.mockCarRepo)
}

func (s *vehicleUnitSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func TestVehicleUnitSuite(t *testing.T) {
	suite.Run(t, new(vehicleUnitSuite))
}

func (s *vehicleUnitSuite) TestCreateCar() {
	const (
		vin              = "4Y1SL65848Z411439"
		brandName        = "Mercedes"
		brandModel       = "C Class"
		color            = "Blue Espectral"
		engineType       = "Mild Hybrid"
		transmissionType = "Automatic"
	)

	id := uuid.New()

	expectedCar := vehicle.HydrateCar(id, s.now, s.now, vin, brandName, brandModel, color, engineType, transmissionType)

	s.mockCarRepo.EXPECT().Create(s.ctx, expectedCar).Return(nil)

	err := s.sut.Handle(s.ctx, &vehicle.CreateCarCommand{
		ID:               id,
		VIN:              vin,
		BrandName:        brandName,
		BrandModel:       brandModel,
		Color:            color,
		EngineType:       engineType,
		TransmissionType: transmissionType,
	})
	s.Require().NoError(err)
}
