package kafka

import (
	"fmt"
	"regexp"

	"github.com/Shopify/sarama"

	"vehicle-sharing-go/app/inventory/cmd/vehicle/message-relay/relay"
)

// Topic describes a Kafka topic and the aggregate type pattern whose events
// should be routed to it.
type Topic struct {
	Name          string
	TopicDetail   *sarama.TopicDetail
	AggregateType *regexp.Regexp
}

// HeaderMapping maps an outbox column name to a Kafka message header name.
type HeaderMapping struct {
	ColumnName string
	HeaderName string
}

// EventDispatcher implements relay.EventDispatcher by routing OutboxEvents to
// Kafka topics and publishing them via a sarama.SyncProducer.
type EventDispatcher struct {
	syncProducer   sarama.SyncProducer
	admin          sarama.ClusterAdmin
	topics         []Topic
	headerMappings []HeaderMapping
}

// NewEventDispatcher creates an EventDispatcher, ensuring all configured topics
// exist in Kafka before returning.
func NewEventDispatcher(
	syncProducer sarama.SyncProducer,
	admin sarama.ClusterAdmin,
	topics []Topic,
	headerMappings []HeaderMapping,
) (*EventDispatcher, error) {
	if err := createMissingTopics(topics, admin); err != nil {
		return nil, err
	}

	return &EventDispatcher{
		syncProducer:   syncProducer,
		admin:          admin,
		topics:         topics,
		headerMappings: headerMappings,
	}, nil
}

// Dispatch routes the event to every topic whose AggregateType regexp matches
// the event's AggregateType, then publishes it.
func (d *EventDispatcher) Dispatch(event relay.OutboxEvent) error {
	for _, topic := range d.topics {
		if !topic.AggregateType.MatchString(string(event.AggregateType)) {
			continue
		}

		headers, err := d.mapHeaders(event.Columns)
		if err != nil {
			return err
		}

		_, _, err = d.syncProducer.SendMessage(&sarama.ProducerMessage{
			Key:     sarama.ByteEncoder(event.AggregateID),
			Topic:   topic.Name,
			Value:   sarama.ByteEncoder(event.Payload),
			Headers: headers,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func createMissingTopics(topics []Topic, admin sarama.ClusterAdmin) error {
	names := make([]string, 0, len(topics))
	for _, t := range topics {
		names = append(names, t.Name)
	}

	metadata, err := admin.DescribeTopics(names)
	if err != nil {
		return err
	}

	for _, m := range metadata {
		if m.Err != sarama.ErrUnknownTopicOrPartition {
			continue
		}
		for _, t := range topics {
			if t.Name != m.Name {
				continue
			}
			if err := admin.CreateTopic(t.Name, t.TopicDetail, false); err != nil {
				return err
			}
		}
	}

	return nil
}

func (d *EventDispatcher) mapHeaders(cols []relay.Column) ([]sarama.RecordHeader, error) {
	headers := make([]sarama.RecordHeader, 0, len(d.headerMappings))

outerLoop:
	for _, mapping := range d.headerMappings {
		for _, col := range cols {
			if mapping.ColumnName == string(col.Name) {
				headers = append(headers, sarama.RecordHeader{
					Key:   []byte(mapping.HeaderName),
					Value: col.Value,
				})
				continue outerLoop
			}
		}
		return nil, fmt.Errorf(
			"column %q not found for header mapping to %q",
			mapping.ColumnName, mapping.HeaderName,
		)
	}

	return headers, nil
}
