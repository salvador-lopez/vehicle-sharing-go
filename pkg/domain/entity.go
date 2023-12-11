package domain

import (
	"time"

	"github.com/google/uuid"
)

type BaseEntity struct {
	id        uuid.UUID
	createdAt time.Time
	updatedAt time.Time
}

func (e *BaseEntity) CreatedAt() time.Time {
	return e.createdAt
}

func (e *BaseEntity) UpdatedAt() time.Time {
	return e.updatedAt
}
