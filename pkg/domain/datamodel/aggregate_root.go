package datamodel

import (
	"time"

	"github.com/google/uuid"
)

type AggregateRoot struct {
	ID        uuid.UUID `gorm:"<-:create;type:varchar(36)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
