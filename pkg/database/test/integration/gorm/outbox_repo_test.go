//go:build integration

package gorm

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/suite"

	"vehicle-sharing-go/pkg/database/gorm"
	"vehicle-sharing-go/pkg/database/gorm/model"
	"vehicle-sharing-go/pkg/domain/event"
)

type outboxRepoIntegrationSuite struct {
	DatabaseSuite
	evtID      uuid.UUID
	appID      string
	kafkaTopic string
	sut        *gorm.OutboxRepository
}

func (s *outboxRepoIntegrationSuite) SetupSuite() {
	s.DatabaseSuite.SetupSuite()
	s.initDb()
	s.evtID = uuid.New()
	s.sut = gorm.NewOutboxRepository(s.Conn())
}

func (s *outboxRepoIntegrationSuite) initDb() {
	s.Require().NoError(s.Conn().Db().AutoMigrate(&model.OutboxRecord{}))
}

func (s *outboxRepoIntegrationSuite) TearDownTest() {
	s.Conn().Db().Delete(&model.OutboxRecord{}, s.evtID)
	s.DatabaseSuite.TearDownTest()
}

func TestOutboxRepoIntegrationSuite(t *testing.T) {
	suite.Run(t, new(outboxRepoIntegrationSuite))
}

func (s *outboxRepoIntegrationSuite) TestPublish() {
	now := time.Now()

	var events []*event.Event
	evtPayload := &CarCreatedPayload{
		VinNumber: "4Y1SL65848Z411439",
		Color:     "Spectral Blue",
		CreatedAt: now,
		UpdatedAt: now,
	}
	carCreatedEvent := &event.Event{
		ID:            s.evtID,
		AggregateID:   uuid.New(),
		AggregateType: "Car",
		EventType:     "CarCreatedEvent",
		Payload:       evtPayload,
		Timestamp:     now,
	}
	events = append(events, carCreatedEvent)
	s.Require().NoError(s.sut.Publish(s.Ctx(), events))

	var outboxRecordStored *model.OutboxRecord
	s.Conn().Db().First(&outboxRecordStored, s.evtID)
	s.Require().NotNil(outboxRecordStored)

	s.Require().Equal(carCreatedEvent.AggregateID, outboxRecordStored.AggregateID)
	s.Require().Equal(carCreatedEvent.AggregateType, outboxRecordStored.AggregateType)
	s.Require().Equal(carCreatedEvent.EventType, outboxRecordStored.EventType)

	var evtPayloadFound *CarCreatedPayload
	s.Require().NoError(decode(outboxRecordStored.Payload, &evtPayloadFound))

	s.Require().Equal(evtPayload.Color, evtPayloadFound.Color)
	s.Require().Equal(evtPayload.VinNumber, evtPayloadFound.VinNumber)
	RequireEqualDates(evtPayload.CreatedAt, evtPayloadFound.CreatedAt, s.Require())
	RequireEqualDates(evtPayload.UpdatedAt, evtPayloadFound.UpdatedAt, s.Require())

	RequireEqualDates(carCreatedEvent.Timestamp, outboxRecordStored.CreatedAt, s.Require())
}

func toTimeHookFunc() mapstructure.DecodeHookFunc {
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

func decode(input any, result interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: nil,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			toTimeHookFunc()),
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

type CarCreatedPayload struct {
	VinNumber string `gorm:"type:varchar(255);unique"`
	Color     string `gorm:"type:varchar(255)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
