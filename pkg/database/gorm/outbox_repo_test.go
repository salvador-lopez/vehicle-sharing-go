//go:build unit

package gorm_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	gormdb "gorm.io/gorm"

	databasegorm "vehicle-sharing-go/pkg/database/gorm"
	"vehicle-sharing-go/pkg/database/gorm/model"
	"vehicle-sharing-go/pkg/domain/event"
)

type outboxRecordsCapture struct {
	records []*model.OutboxRecord
}

func newOutboxTestConnection(t *testing.T) (*databasegorm.Connection, *outboxRecordsCapture) {
	t.Helper()

	connection, err := databasegorm.NewConnectionFromConfig(
		&databasegorm.Config{SkipInitializeWithVersion: true},
		databasegorm.DSNConfig{
			UserName:     "gorm",
			Password:     "gorm",
			DatabaseName: "gorm",
			Host:         "localhost",
			Port:         9910,
		},
		&gormdb.Config{
			DryRun:                 true,
			DisableAutomaticPing:   true,
			SkipDefaultTransaction: true,
		},
	)
	require.NoError(t, err)

	capture := &outboxRecordsCapture{}
	require.NoError(t, connection.Db().Callback().Create().Before("gorm:create").Register("capture_outbox_records", func(tx *gormdb.DB) {
		capture.records, _ = tx.Statement.Dest.([]*model.OutboxRecord)
	}))

	return connection, capture
}

func newOutboxTestEvent(eventType, vin string) *event.Event {
	return &event.Event{
		ID:            uuid.New(),
		AggregateID:   uuid.New(),
		AggregateType: "Vehicle",
		EventType:     eventType,
		Payload:       map[string]string{"vin": vin},
		Timestamp:     time.Now(),
	}
}

func TestOutboxRepositoryPublishMapsEventsToRecords(t *testing.T) {
	firstEvent := newOutboxTestEvent("VehicleCreated", "vehicle-123")
	secondEvent := newOutboxTestEvent("VehicleUpdated", "vehicle-456")

	tests := []struct {
		name   string
		events []*event.Event
	}{
		{name: "single event", events: []*event.Event{firstEvent}},
		{name: "multiple events", events: []*event.Event{firstEvent, secondEvent}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connection, capture := newOutboxTestConnection(t)

			err := databasegorm.NewOutboxRepository(connection).Publish(context.Background(), tt.events)
			require.NoError(t, err)
			require.Len(t, capture.records, len(tt.events))
			for i, evt := range tt.events {
				require.Equal(t, &model.OutboxRecord{
					ID:            evt.ID,
					CreatedAt:     evt.Timestamp,
					AggregateType: evt.AggregateType,
					AggregateID:   evt.AggregateID,
					EventType:     evt.EventType,
					Payload:       evt.Payload,
				}, capture.records[i])
			}
		})
	}
}

func TestOutboxRepositoryPublishReturnsCreateError(t *testing.T) {
	connection, _ := newOutboxTestConnection(t)
	createErr := errors.New("failed to create outbox records")
	require.NoError(t, connection.Db().Callback().Create().Before("gorm:create").Register("fail_outbox_create", func(tx *gormdb.DB) {
		_ = tx.AddError(createErr)
	}))

	evt := newOutboxTestEvent("VehicleCreated", "vehicle-123")

	err := databasegorm.NewOutboxRepository(connection).Publish(context.Background(), []*event.Event{evt})

	require.ErrorIs(t, err, createErr)
}
