package messaging

import (
	"context"
	"log"
	"time"
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
	pollInterval time.Duration
	batchSize    int
	logger       *log.Logger
}

func NewPollingPublisher(
	outboxReader OutboxReader,
	publisher EventPublisher,
	pollInterval time.Duration,
	batchSize int,
	logger *log.Logger,
) *PollingPublisher {
	return &PollingPublisher{
		outboxReader: outboxReader,
		publisher:    publisher,
		pollInterval: pollInterval,
		batchSize:    batchSize,
		logger:       logger,
	}
}

func (p *PollingPublisher) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			p.logger.Println("Polling publisher shutting down")
			return nil
		default:
		}

		events, err := p.outboxReader.Poll(ctx, p.batchSize)
		if err != nil {
			p.logger.Printf("Error polling outbox: %v", err)
			time.Sleep(p.pollInterval) // back off on error
			continue
		}

		if len(events) == 0 {
			time.Sleep(p.pollInterval)
			continue
		}

		err = p.publisher.Publish(ctx, events)
		if err != nil {
			p.logger.Printf("Error publishing events: %v", err)
			time.Sleep(p.pollInterval) // back off on error
			continue
		}
	}
}
