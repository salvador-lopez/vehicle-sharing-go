package kafka

import (
	"context"
	"encoding/json"
	"vehicle-sharing-go/pkg/domain/event"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type EventPublisher struct {
	producer *kafka.Producer
	topic    string
}

func NewEventPublisher(producer *kafka.Producer, topic string) *EventPublisher {
	return &EventPublisher{producer: producer, topic: topic}
}

func (p EventPublisher) Publish(ctx context.Context, events []*event.Event) error {
	for _, evt := range events {
		payloadBytes, err := json.Marshal(evt.Payload)
		if err != nil {
			return err
		}

		headers := []kafka.Header{
			{Key: "event_type", Value: []byte(evt.EventType)},
			{Key: "aggregate_type", Value: []byte(evt.AggregateType)},
		}

		kafkaMsg := &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
			Key:            []byte(evt.AggregateID.String()),
			Value:          payloadBytes,
			Headers:        headers,
		}

		err = p.producer.Produce(kafkaMsg, nil)
		if err != nil {
			return err
		}
	}

	// Flush to ensure delivery
	p.producer.Flush(1000)

	return nil
}
