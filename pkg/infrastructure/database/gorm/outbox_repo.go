package gorm

import (
	"context"

	"vehicle-sharing-go/pkg/domain/event"
	"vehicle-sharing-go/pkg/infrastructure/database/gorm/model"
)

type OutboxRepository struct {
	conn *Connection
}

func (o *OutboxRepository) Publish(ctx context.Context, events []*event.Event) error {
	var gormEvents []*model.Event
	for _, evt := range events {
		gormEvents = append(gormEvents, &model.Event{Event: evt})
	}

	return o.conn.Db().WithContext(ctx).Create(gormEvents).Error
}

func NewOutboxRepository(conn *Connection) *OutboxRepository {
	return &OutboxRepository{conn: conn}
}
