package command_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"vehicle-sharing-go/internal/inventory/vehicle/application/command"
	"vehicle-sharing-go/internal/inventory/vehicle/domain"
	"vehicle-sharing-go/internal/inventory/vehicle/domain/mock"
)

type createCarUnitSuite struct {
	suite.Suite
	ctx              context.Context
	mockCtrl         *gomock.Controller
	now              time.Time
	mockVinValidator *mock.MockVinValidator
	mockCarRepo      *mock.MockCarRepository
	sut              *command.CreateCarHandler
}

func (s *createCarUnitSuite) SetupTest() {
	s.ctx = context.Background()
	s.now = time.Now()
	s.mockCtrl = gomock.NewController(s.T())
	s.mockVinValidator = mock.NewMockVinValidator(s.mockCtrl)
	s.mockCarRepo = mock.NewMockCarRepository(s.mockCtrl)
	s.sut = command.NewCreateCarHandler(func() time.Time { return s.now }, s.mockVinValidator, s.mockCarRepo)
}

func (s *createCarUnitSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func TestVehicleUnitSuite(t *testing.T) {
	suite.Run(t, new(createCarUnitSuite))
}

func (s *createCarUnitSuite) TestCreateCar() {
	const (
		vin   = "4Y1SL65848Z411439"
		color = "Spectral Blue"
	)

	id := uuid.New()

	expectedCar := domain.HydrateCar(&domain.CarDTO{
		VIN:     vin,
		Color:   color,
		BaseDTO: &domain.BaseDTO{ID: id, CreatedAt: s.now, UpdatedAt: s.now},
	})

	s.mockVinValidator.EXPECT().Validate(vin).Return(nil)
	s.mockCarRepo.EXPECT().Create(s.ctx, expectedCar).Return(nil)

	err := s.sut.Handle(s.ctx, &command.CreateCar{ID: id, VIN: vin, Color: color})
	s.Require().NoError(err)
}
