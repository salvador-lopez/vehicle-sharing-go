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
type eventRecorder interface {
	RecordedEvents() []*Event
}

func PublishRecordedEvents(ctx context.Context, topic string, aggRoot eventRecorder, publisher EventPublisher) error {
	recordedEvents := aggRoot.RecordedEvents()
	err := publisher.Publish(ctx, topic, recordedEvents)

	if err != nil {
		recordedEventsBytes, _ := json.Marshal(recordedEvents)
		return fmt.Errorf("%w, cause: %v, events: %v", publishEventsErr, errors.New("publish error"), string(recordedEventsBytes))
	}

	return err
}
