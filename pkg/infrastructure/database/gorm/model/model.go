package model

import (
	"time"

	"github.com/google/uuid"
)

type OutboxRecord struct {
	ID                uuid.UUID `gorm:"<-:create;type:varchar(36)"`
	CreatedAt         time.Time
	AggregateType     string    `gorm:"type:varchar(100)"`
	AggregateID       uuid.UUID `gorm:"type:varchar(36);unique"`
	Payload           any       `gorm:"serializer:json"`
	EventType         string    `gorm:"type:varchar(255)"`
	KafkaHeaderKeys   []string  `gorm:"serializer:json"`
	KafkaHeaderValues []string  `gorm:"serializer:json"`
}

func (o *OutboxRecord) TableName() string {
	return "outbox"
}
