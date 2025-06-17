package messaging

import (
	"context"
	"log"
	"vehicle-sharing-go/pkg/domain/event"
)

type OutboxReader interface {
	Poll(ctx context.Context, limit int) ([]*event.Event, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, events []*event.Event) error
}

type PollingPublisher struct {
	outboxReader OutboxReader
	publisher    EventPublisher
	batchSize    int
	logger       *log.Logger
}

func NewPollingPublisher(
	outboxReader OutboxReader,
	publisher EventPublisher,
	logger *log.Logger,
) *PollingPublisher {
	return &PollingPublisher{
		outboxReader: outboxReader,
		publisher:    publisher,
		logger:       logger,
	}
}

func (p *PollingPublisher) RunBatch(ctx context.Context, batchSize int) error {
	events, err := p.outboxReader.Poll(ctx, batchSize)
	if err != nil {
		return err
	}

	return p.publisher.Publish(ctx, events)
}
