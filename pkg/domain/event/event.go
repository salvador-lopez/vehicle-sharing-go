package event

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID            uuid.UUID
	AggregateID   uuid.UUID
	AggregateType string
	EventType     string
	Payload       any
	Timestamp     time.Time
}
