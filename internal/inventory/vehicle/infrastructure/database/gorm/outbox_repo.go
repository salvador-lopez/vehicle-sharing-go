package gorm

import (
	"context"

	"gorm.io/gorm"

	"vehicle-sharing-go/internal/inventory/vehicle/application/command"
)

type OutboxRepository struct {
	db *gorm.DB
}

func NewOutboxRepository(db *gorm.DB) *OutboxRepository {
	return &OutboxRepository{db: db}
}

func (o *OutboxRepository) Append(ctx context.Context, eventsRecorder command.DomainEventsRecorder) error {
	return o.db.WithContext(ctx).Create(eventsRecorder.RecordedEvents()).Error
}
