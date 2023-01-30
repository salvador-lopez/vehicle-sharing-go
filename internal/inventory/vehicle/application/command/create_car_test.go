package command_test

import (
	"context"
	"fmt"
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

	s.mockCarRepo.EXPECT().Create(s.ctx, expectedCar).Return(nil)

	err := s.handleSut(id, vin, color)
	s.Require().NoError(err)
}

func (s *createCarUnitSuite) handleSut(id uuid.UUID, vin, color string) error {
	return s.sut.Handle(s.ctx, &command.CreateCar{ID: id, VIN: vin, Color: color})
}

func (s *createCarUnitSuite) TestCreateCarReturnErrInvalidVin() {
	tests := []struct {
		name string
		vin  string
	}{
		{
			name: "Less than 17 characters",
			vin:  "1FUYDMDB2YPF8709",
		},
		{
			name: "More than 17 characters",
			vin:  "4Y1SL65848Z4114399",
		},
		{
			name: "Letter in lower case",
			vin:  "4y1SL65848Z411439",
		},
		{
			name: "Letter 'I' not allowed in position 1",
			vin:  "IY1SL65848Z411439",
		},
		{
			name: "Letter 'O' not allowed in position 1",
			vin:  "OY1SL65848Z411439",
		},
		{
			name: "Letter 'I' not allowed in position 17",
			vin:  "3Y1SL65848Z41143I",
		},
		{
			name: "Letter 'O' not allowed in position 17",
			vin:  "3Y1SL65848Z41143O",
		},
		{
			name: "only letter X or number allowed in position 9",
			vin:  "3Y1SL658P8Z411439",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			defer s.TearDownTest()

			s.mockCarRepo.EXPECT().Create(gomock.Any(), gomock.Any()).AnyTimes()

			invalidVinErr := fmt.Errorf("%v: %s", domain.ErrInvalidVin, tt.vin)
			err := s.sut.Handle(s.ctx, &command.CreateCar{ID: uuid.New(), VIN: tt.vin, Color: "Blue Spectral"})
			s.Require().EqualError(err, invalidVinErr.Error())
		})
	}
}
