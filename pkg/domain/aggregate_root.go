package domain

import (
	"time"

	"github.com/google/uuid"

	"vehicle-sharing-go/pkg/domain/event"
)

type AggregateRoot struct {
	*BaseEntity
	recordedEvents []*event.Event
}

func NewAggregateRoot(id uuid.UUID, nowFun func() time.Time) *AggregateRoot {
	return &AggregateRoot{BaseEntity: &BaseEntity{id: id, createdAt: nowFun(), updatedAt: nowFun()}}
}

func (a *AggregateRoot) RecordedEvents() []*event.Event {
	return a.recordedEvents
}

func (a *AggregateRoot) RecordEvent(eventID uuid.UUID, eventType, aggregateType string, payload any, timestamp time.Time) {
	a.recordedEvents = append(
		a.recordedEvents,
		&event.Event{
			ID:            eventID,
			AggregateID:   a.id,
			AggregateType: aggregateType,
			EventType:     eventType,
			Payload:       payload,
			Timestamp:     timestamp,
		},
	)
}
