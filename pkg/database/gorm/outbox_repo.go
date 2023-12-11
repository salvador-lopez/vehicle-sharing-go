package gorm

import (
	"context"

	"vehicle-sharing-go/pkg/database/gorm/model"
	"vehicle-sharing-go/pkg/domain/event"
)

type OutboxRepository struct {
	conn *Connection
}

func NewOutboxRepository(conn *Connection) *OutboxRepository {
	return &OutboxRepository{conn: conn}
}

func (o *OutboxRepository) Publish(ctx context.Context, events []*event.Event) error {
	var outboxRecords []*model.OutboxRecord
	for _, evt := range events {
		outboxRecords = append(outboxRecords, &model.OutboxRecord{
			ID:            evt.ID,
			CreatedAt:     evt.Timestamp,
			EventType:     evt.EventType,
			AggregateType: evt.AggregateType,
			AggregateID:   evt.AggregateID,
			Payload:       evt.Payload,
		})
	}

	return o.conn.Db().WithContext(ctx).Create(outboxRecords).Error
}
