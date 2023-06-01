package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

var publishEventsErr = errors.New("failed to publish events recorded in the aggregate root")

//go:generate mockgen -destination=mock/event_publisher_mock.go -package=mock . EventPublisher
type EventPublisher interface {
	Publish(ctx context.Context, topic string, events []*Event) error
}

//go:generate mockgen -destination=mock/event_recorder_mock.go -package=mock . EventRecorder
type EventRecorder interface {
	RecordedEvents() []*Event
}

type AgRootEventPublisher struct {
	publisher EventPublisher
}

func NewAgRootEventPublisher(publisher EventPublisher) *AgRootEventPublisher {
	return &AgRootEventPublisher{publisher: publisher}
}

func (ep *AgRootEventPublisher) Publish(ctx context.Context, topic string, evtRecorder EventRecorder) error {
	recordedEvents := evtRecorder.RecordedEvents()
	err := ep.publisher.Publish(ctx, topic, recordedEvents)

	if err != nil {
		recordedEventsBytes, _ := json.Marshal(recordedEvents)
		return fmt.Errorf("%w, cause: %v, events: %v", publishEventsErr, errors.New("publish error"), string(recordedEventsBytes))
	}

	return err
}
