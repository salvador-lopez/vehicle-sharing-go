//go:build unit_inventory

package command_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"vehicle-sharing-go/internal/inventory/vehicle/application/command"
	"vehicle-sharing-go/internal/inventory/vehicle/application/command/mock"
	"vehicle-sharing-go/internal/inventory/vehicle/domain"
	"vehicle-sharing-go/internal/inventory/vehicle/domain/event"
	"vehicle-sharing-go/internal/inventory/vehicle/domain/model"
	event2 "vehicle-sharing-go/pkg/domain/event"
	modelpkg "vehicle-sharing-go/pkg/domain/model"
)

const (
	validVinNumber = "4Y1SL65848Z411439"
	color          = "Spectral Blue"
)

type createCarUnitSuite struct {
	suite.Suite
	ctx              context.Context
	mockCtrl         *gomock.Controller
	idGen            func() uuid.UUID
	now              func() time.Time
	mockCarRepo      *mock.MockCarRepository
	mockEvtPublisher *mock.MockEventPublisher
	sut              *command.CreateCarHandler
}

func (s *createCarUnitSuite) SetupTest() {
	s.ctx = context.Background()

	carCreatedEvtID := uuid.New()
	s.idGen = func() uuid.UUID { return carCreatedEvtID }
	now := time.Now()
	s.now = func() time.Time { return now }

	s.mockCtrl = gomock.NewController(s.T())
	s.mockCarRepo = mock.NewMockCarRepository(s.mockCtrl)
	s.mockEvtPublisher = mock.NewMockEventPublisher(s.mockCtrl)

	s.sut = command.NewCreateCarHandler(s.idGen, s.now, s.mockCarRepo, s.mockEvtPublisher)
}

func (s *createCarUnitSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func TestCreateCarUnitSuite(t *testing.T) {
	suite.Run(t, new(createCarUnitSuite))
}

func (s *createCarUnitSuite) TestCreateCar() {
	id := uuid.New()
	now := s.now()

	recordedEvent := &event2.Event{
		ID:            s.idGen(),
		AggregateID:   id,
		AggregateType: "Car",
		EventType:     "CarCreatedEvent",
		Payload: &event.CarCreatedPayload{
			VinNumber: validVinNumber,
			Color:     color,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Timestamp: now,
	}

	expectedCar := domain.CarFromModel(&model.Car{
		VinNumber: validVinNumber,
		Color:     color,
		AggregateRoot: &modelpkg.AggregateRoot{
			ID:             id,
			CreatedAt:      now,
			UpdatedAt:      now,
			RecordedEvents: []*event2.Event{recordedEvent},
		},
	})

	s.mockCarRepo.EXPECT().Create(s.ctx, expectedCar).Return(nil)

	s.mockEvtPublisher.EXPECT().Publish(s.ctx, "inventory", expectedCar).Return(nil)

	err := s.handleSut(id, validVinNumber)
	s.Require().NoError(err)
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
			err := s.handleSut(uuid.New(), tt.vin)
			s.Require().EqualError(err, invalidVinErr.Error())
		})
	}
}

func (s *createCarUnitSuite) TestReturnRepositoryErr() {
	repoErr := errors.New("repository error")
	s.mockCarRepo.EXPECT().Create(s.ctx, gomock.Any()).Return(repoErr)

	err := s.handleSut(uuid.New(), validVinNumber)
	s.Require().EqualError(err, repoErr.Error())
}

func (s *createCarUnitSuite) handleSut(id uuid.UUID, vin string) error {
	return s.sut.Handle(s.ctx, &command.CreateCar{ID: id, VIN: vin, Color: color})
}
