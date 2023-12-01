package event

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID            uuid.UUID `gorm:"<-:create;type:varchar(36)"`
	AggregateID   uuid.UUID `gorm:"type:varchar(36);unique"`
	AggregateType string    `gorm:"type:varchar(255)"`
	EventType     string    `gorm:"type:varchar(255)"`
	Payload       any       `gorm:"serializer:json"`
	Timestamp     time.Time
}
