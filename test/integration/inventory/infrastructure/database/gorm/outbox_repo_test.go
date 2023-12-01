//go:build integration_inventory

package gorm

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/suite"

	"vehicle-sharing-go/internal/inventory/vehicle/domain/event"
	eventpkg "vehicle-sharing-go/pkg/domain/event"
	"vehicle-sharing-go/pkg/infrastructure/database/gorm"
	modelpkg "vehicle-sharing-go/pkg/infrastructure/database/gorm/model"
)

type outboxRepoIntegrationSuite struct {
	databaseSuite
	evtID uuid.UUID
	sut   *gorm.OutboxRepository
}

func (s *outboxRepoIntegrationSuite) SetupSuite() {
	s.databaseSuite.SetupSuite()
	s.initDb()
	s.evtID = uuid.New()
	s.sut = gorm.NewOutboxRepository(s.conn)
}

func (s *outboxRepoIntegrationSuite) initDb() {
	s.Require().NoError(s.conn.Db().AutoMigrate(&modelpkg.Event{}))
}

func (s *outboxRepoIntegrationSuite) TearDownTest() {
	s.conn.Db().Delete(&eventpkg.Event{}, s.evtID)
	s.databaseSuite.TearDownTest()
}

func TestOutboxRepoIntegrationSuite(t *testing.T) {
	suite.Run(t, new(outboxRepoIntegrationSuite))
}

func (s *outboxRepoIntegrationSuite) TestPublish() {
	now := time.Now()

	var events []*eventpkg.Event
	evtPayload := &event.CarCreatedPayload{
		VinNumber: "4Y1SL65848Z411439",
		Color:     "Spectral Blue",
		CreatedAt: now,
		UpdatedAt: now,
	}
	carCreatedEvent := &eventpkg.Event{
		ID:            s.evtID,
		AggregateID:   uuid.New(),
		AggregateType: "Car",
		EventType:     "CarCreatedEvent",
		Payload:       evtPayload,
		Timestamp:     now,
	}
	events = append(events, carCreatedEvent)
	s.Require().NoError(s.sut.Publish(s.ctx, events))

	var gormEvtStored *modelpkg.Event
	s.conn.Db().First(&gormEvtStored, s.evtID)
	s.Require().NotNil(gormEvtStored)

	s.Require().Equal(carCreatedEvent.AggregateID, gormEvtStored.AggregateID)
	s.Require().Equal(carCreatedEvent.AggregateType, gormEvtStored.AggregateType)

	var evtPayloadFound *event.CarCreatedPayload
	s.Require().NoError(Decode(gormEvtStored.Payload, &evtPayloadFound))

	s.Require().Equal(evtPayload.Color, evtPayloadFound.Color)
	s.Require().Equal(evtPayload.VinNumber, evtPayloadFound.VinNumber)
	requireEqualDates(evtPayload.CreatedAt, evtPayloadFound.CreatedAt, s.Require())
	requireEqualDates(evtPayload.UpdatedAt, evtPayloadFound.UpdatedAt, s.Require())

	requireEqualDates(carCreatedEvent.Timestamp, gormEvtStored.Timestamp, s.Require())
}

func ToTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.String:
			return time.Parse(time.RFC3339, data.(string))
		case reflect.Float64:
			return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
		case reflect.Int64:
			return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
		default:
			return data, nil
		}
	}
}

func Decode(input any, result interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: nil,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			ToTimeHookFunc()),
		Result: result,
	})
	if err != nil {
		return err
	}

	if err := decoder.Decode(input); err != nil {
		return err
	}
	return err
}
