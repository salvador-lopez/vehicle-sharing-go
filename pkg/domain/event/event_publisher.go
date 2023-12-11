package event

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

var publishEventsErr = errors.New("failed to publish events recorded in the aggregate root")

//go:generate mockgen -destination=mock/publisher_mock.go -package=mock . Publisher
type Publisher interface {
	Publish(ctx context.Context, events []*Event) error
}

//go:generate mockgen -destination=mock/recorder_mock.go -package=mock . Recorder
type Recorder interface {
	RecordedEvents() []*Event
}

type AgRootEventPublisher struct {
	publisher Publisher
}

func NewAgRootEventPublisher(publisher Publisher) *AgRootEventPublisher {
	return &AgRootEventPublisher{publisher: publisher}
}

func (ep *AgRootEventPublisher) Publish(ctx context.Context, evtRecorder Recorder) error {
	recordedEvents := evtRecorder.RecordedEvents()
	err := ep.publisher.Publish(ctx, recordedEvents)

	if err != nil {
		recordedEventsBytes, _ := json.Marshal(recordedEvents)
		return fmt.Errorf("%w, cause: %v, events: %v", publishEventsErr, errors.New("publish error"), string(recordedEventsBytes))
	}

	return err
}
