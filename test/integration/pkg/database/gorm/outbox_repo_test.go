//go:build integration

package gorm

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/suite"

	"vehicle-sharing-go/internal/inventory/vehicle/domain/event"
	gormpkg "vehicle-sharing-go/pkg/database/gorm"
	modelpkg "vehicle-sharing-go/pkg/database/gorm/model"
	eventpkg "vehicle-sharing-go/pkg/domain/event"
	"vehicle-sharing-go/test/integration/database/gorm"
)

type outboxRepoIntegrationSuite struct {
	gorm.DatabaseSuite
	evtID      uuid.UUID
	appID      string
	kafkaTopic string
	sut        *gormpkg.OutboxRepository
}

func (s *outboxRepoIntegrationSuite) SetupSuite() {
	s.DatabaseSuite.SetupSuite()
	s.initDb()
	s.evtID = uuid.New()
	s.sut = gormpkg.NewOutboxRepository(s.Conn())
}

func (s *outboxRepoIntegrationSuite) initDb() {
	s.Require().NoError(s.Conn().Db().AutoMigrate(&modelpkg.OutboxRecord{}))
}

func (s *outboxRepoIntegrationSuite) TearDownTest() {
	s.Conn().Db().Delete(&modelpkg.OutboxRecord{}, s.evtID)
	s.DatabaseSuite.TearDownTest()
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
	s.Require().NoError(s.sut.Publish(s.Ctx(), events))

	var outboxRecordStored *modelpkg.OutboxRecord
	s.Conn().Db().First(&outboxRecordStored, s.evtID)
	s.Require().NotNil(outboxRecordStored)

	s.Require().Equal(carCreatedEvent.AggregateID, outboxRecordStored.AggregateID)
	s.Require().Equal(carCreatedEvent.AggregateType, outboxRecordStored.AggregateType)
	s.Require().Equal(carCreatedEvent.EventType, outboxRecordStored.EventType)

	var evtPayloadFound *event.CarCreatedPayload
	s.Require().NoError(Decode(outboxRecordStored.Payload, &evtPayloadFound))

	s.Require().Equal(evtPayload.Color, evtPayloadFound.Color)
	s.Require().Equal(evtPayload.VinNumber, evtPayloadFound.VinNumber)
	gorm.RequireEqualDates(evtPayload.CreatedAt, evtPayloadFound.CreatedAt, s.Require())
	gorm.RequireEqualDates(evtPayload.UpdatedAt, evtPayloadFound.UpdatedAt, s.Require())

	gorm.RequireEqualDates(carCreatedEvent.Timestamp, outboxRecordStored.CreatedAt, s.Require())
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