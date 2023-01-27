package command_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"vehicle-sharing-go/internal/inventory/vehicle/command"
	"vehicle-sharing-go/internal/inventory/vehicle/domain"
	"vehicle-sharing-go/internal/inventory/vehicle/domain/mock"
)

type createCarUnitSuite struct {
	suite.Suite
	ctx         context.Context
	mockCtrl    *gomock.Controller
	now         time.Time
	mockCarRepo *mock.MockCarRepository
	sut         *command.CreateCarHandler
}

func (s *createCarUnitSuite) SetupTest() {
	s.ctx = context.Background()
	s.now = time.Now()
	s.mockCtrl = gomock.NewController(s.T())
	s.mockCarRepo = mock.NewMockCarRepository(s.mockCtrl)
	s.sut = command.NewCreateCarHandler(func() time.Time { return s.now }, s.mockCarRepo)
}

func (s *createCarUnitSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func TestVehicleUnitSuite(t *testing.T) {
	suite.Run(t, new(createCarUnitSuite))
}

func (s *createCarUnitSuite) TestCreateCar() {
	const vin = "4Y1SL65848Z411439"

	id := uuid.New()

	expectedCar := domain.HydrateCar(id, s.now, s.now, vin)

	s.mockCarRepo.EXPECT().Create(s.ctx, expectedCar).Return(nil)

	err := s.sut.Handle(s.ctx, &command.CreateCar{ID: id, VIN: vin})
	s.Require().NoError(err)
}
