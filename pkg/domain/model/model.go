package model

import (
	"time"

	"github.com/google/uuid"

	"vehicle-sharing-go/pkg/domain/event"
)

type AggregateRoot struct {
	ID             uuid.UUID `gorm:"<-:create;type:varchar(36)"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	RecordedEvents []*event.Event `gorm:"-"`
}
