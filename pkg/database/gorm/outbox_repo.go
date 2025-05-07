package gorm

import (
	"context"
	"time"

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
	var records []*model.OutboxRecord
	for _, evt := range events {
		records = append(records, &model.OutboxRecord{
			ID:            evt.ID,
			CreatedAt:     evt.Timestamp,
			EventType:     evt.EventType,
			AggregateType: evt.AggregateType,
			AggregateID:   evt.AggregateID,
			Payload:       evt.Payload,
		})
	}

	return o.conn.Db().WithContext(ctx).Create(records).Error
}

func (o *OutboxRepository) PollAfter(ctx context.Context, after time.Time, limit int) ([]*event.Event, error) {
	var records []*model.OutboxRecord

	err := o.conn.Db().
		WithContext(ctx).
		Where("created_at > ?", after).
		Order("created_at ASC").
		Limit(limit).
		Find(&records).Error

	if err != nil {
		return nil, err
	}

	var events []*event.Event
	for _, r := range records {
		events = append(events, &event.Event{
			ID:            r.ID,
			Timestamp:     r.CreatedAt,
			EventType:     r.EventType,
			AggregateType: r.AggregateType,
			AggregateID:   r.AggregateID,
			Payload:       r.Payload,
		})
	}

	return events, nil
}
