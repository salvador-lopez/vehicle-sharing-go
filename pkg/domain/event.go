package domain

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	id            uuid.UUID
	aggregateID   uuid.UUID
	aggregateType string
	eventType     string
	payload       any
	timestamp     time.Time
}

func NewEvent(
	id uuid.UUID,
	aggregateID uuid.UUID,
	aggregateType string,
	eventType string,
	payload any,
	timestamp time.Time,
) *Event {
	return &Event{
		id:            id,
		aggregateID:   aggregateID,
		aggregateType: aggregateType,
		eventType:     eventType,
		payload:       payload,
		timestamp:     timestamp,
	}
}

type PayloadDTO interface {
	ToPayload() any
}

type EventDTO struct {
	ID            uuid.UUID
	AggregateID   uuid.UUID
	AggregateType string
	EventType     string
	PayloadDTO    PayloadDTO
	Timestamp     time.Time
}

func (dto EventDTO) ToEvent() *Event {
	return &Event{
		id:            dto.ID,
		aggregateID:   dto.AggregateID,
		aggregateType: dto.AggregateType,
		eventType:     dto.EventType,
		payload:       dto.PayloadDTO.ToPayload(),
		timestamp:     dto.Timestamp,
	}
}
