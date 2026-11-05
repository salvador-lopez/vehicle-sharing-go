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
	topics         []Topic
	headerMappings []HeaderMapping
}

// NewEventDispatcher creates an EventDispatcher for the configured topics.
func NewEventDispatcher(
	syncProducer sarama.SyncProducer,
	topics []Topic,
	headerMappings []HeaderMapping,
) *EventDispatcher {
	return &EventDispatcher{
		syncProducer:   syncProducer,
		topics:         topics,
		headerMappings: headerMappings,
	}
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
