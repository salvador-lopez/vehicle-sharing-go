package command

import (
	"context"

	"vehicle-sharing-go/pkg/domain/event"
)

//go:generate mockgen -destination=mock/event_publisher_mock.go -package=mock . EventPublisher
type EventPublisher interface {
	Publish(ctx context.Context, topic string, aggRoot event.Recorder) error
}
