//go:build unit

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

	"vehicle-sharing-go/app/inventory/internal/vehicle/command"
	"vehicle-sharing-go/app/inventory/internal/vehicle/command/mock"
	domain2 "vehicle-sharing-go/app/inventory/internal/vehicle/domain"
	"vehicle-sharing-go/app/inventory/internal/vehicle/domain/event"
	"vehicle-sharing-go/app/inventory/internal/vehicle/domain/model"

	eventpkg "vehicle-sharing-go/pkg/domain/event"
	mockeventpkg "vehicle-sharing-go/pkg/domain/event/mock"
	modelpkg "vehicle-sharing-go/pkg/domain/model"
)

const (
	validVinNumber = "4Y1SL65848Z411439"
	color          = "Spectral Blue"
)

type createCarUnitSuite struct {
	suite.Suite
	ctx                 context.Context
	mockCtrl            *gomock.Controller
	idGen               func() uuid.UUID
	now                 func() time.Time
	mockCarRepo         *mock.MockCarRepository
	mockTxSession       *mock.MockTransactionalSession
	mockPublisher       *mockeventpkg.MockPublisher
	aggRootEvtPublisher *eventpkg.AgRootEventPublisher
	sut                 *command.CreateCarHandler
}

func (s *createCarUnitSuite) SetupTest() {
	s.ctx = context.Background()

	carCreatedEvtID := uuid.New()
	s.idGen = func() uuid.UUID { return carCreatedEvtID }
	now := time.Now()
	s.now = func() time.Time { return now }

	s.mockCtrl = gomock.NewController(s.T())

	s.mockCarRepo = mock.NewMockCarRepository(s.mockCtrl)
	s.mockTxSession = mock.NewMockTransactionalSession(s.mockCtrl)
	s.mockPublisher = mockeventpkg.NewMockPublisher(s.mockCtrl)
	s.aggRootEvtPublisher = eventpkg.NewAgRootEventPublisher(s.mockPublisher)

	s.sut = command.NewCreateCarHandler(s.idGen, s.now, s.mockCarRepo, s.mockTxSession, s.aggRootEvtPublisher)
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

	recordedEvent := &eventpkg.Event{
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

	recordedEvents := []*eventpkg.Event{recordedEvent}

	expectedCar := domain2.CarFromModel(&model.Car{
		VinNumber: validVinNumber,
		Color:     color,
		AggregateRoot: &modelpkg.AggregateRoot{
			ID:             id,
			CreatedAt:      now,
			UpdatedAt:      now,
			RecordedEvents: recordedEvents,
		},
	})

	s.mockCarRepo.EXPECT().Create(s.ctx, expectedCar).Return(nil)
	s.mockPublisher.EXPECT().Publish(s.ctx, recordedEvents).Return(nil)

	s.mockTxSession.EXPECT().Transaction(s.ctx, gomock.Any()).Do(
		func(ctx context.Context, txSessionFunc func(context.Context) error) {
			s.Require().NoError(txSessionFunc(s.ctx))
		},
	).Return(nil)

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

			invalidVinErr := fmt.Errorf("%v: %s", domain2.ErrInvalidVin, tt.vin)
			err := s.handleSut(uuid.New(), tt.vin)
			s.Require().EqualError(err, invalidVinErr.Error())
		})
	}
}

func (s *createCarUnitSuite) TestReturnRepositoryErr() {
	repoErr := errors.New("repository error")
	s.mockCarRepo.EXPECT().Create(s.ctx, gomock.Any()).Return(repoErr)

	s.mockTxSession.EXPECT().Transaction(s.ctx, gomock.Any()).Do(
		func(ctx context.Context, txSessionFunc func(context.Context) error) {
			s.Require().EqualError(txSessionFunc(s.ctx), repoErr.Error())
		},
	).Return(repoErr)

	err := s.handleSut(uuid.New(), validVinNumber)
	s.Require().EqualError(err, repoErr.Error())
}

func (s *createCarUnitSuite) handleSut(id uuid.UUID, vin string) error {
	return s.sut.Handle(s.ctx, &command.CreateCar{ID: id, VIN: vin, Color: color})
}
