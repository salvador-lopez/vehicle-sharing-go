package relay

import (
	"errors"
	"regexp"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/sirupsen/logrus"
)

const (
	defaultAggregateIDColumnName   = "aggregate_id"
	defaultAggregateTypeColumnName = "aggregate_type"
	defaultPayloadColumnName       = "payload"
)

// OutboxEvent is a decoded row from the outbox table.
type OutboxEvent struct {
	AggregateID                []byte
	AggregateType              []byte
	Payload                    []byte
	Columns                    []Column
	EventTimestampFromDatabase uint32
}

// Column is a single column name/value pair from a binlog row event.
type Column struct {
	Name  []byte
	Value []byte
}

// EventDispatcher routes and publishes an OutboxEvent to its destination.
type EventDispatcher interface {
	Dispatch(event OutboxEvent) error
}

// AggregateTypeTopicPair maps a regexp to a topic name (kept for API compatibility).
type AggregateTypeTopicPair struct {
	AggregateTypeRegexp *regexp.Regexp
	Topic               string
}

// NewEventHandler constructs an EventHandler that filters canal row events and
// dispatches decoded OutboxEvents to the given EventDispatcher.
// Column name arguments default to "aggregate_id", "aggregate_type", and
// "payload" when empty.
// The returned channel carries the latest binlog position as reported by
// OnPosSynced; pass it to the BinlogSource so it can checkpoint progress.
func NewEventHandler(
	eventDispatcher EventDispatcher,
	aggregateIDColumnName string,
	aggregateTypeColumnName string,
	payloadColumnName string,
) (*EventHandler, <-chan mysql.Position, error) {
	if aggregateIDColumnName == "" {
		aggregateIDColumnName = defaultAggregateIDColumnName
	}
	if aggregateTypeColumnName == "" {
		aggregateTypeColumnName = defaultAggregateTypeColumnName
	}
	if payloadColumnName == "" {
		payloadColumnName = defaultPayloadColumnName
	}

	posCh := make(chan mysql.Position)

	return &EventHandler{
		eventMapper: &EventMapper{
			aggregateIDColumnName:   aggregateIDColumnName,
			aggregateTypeColumnName: aggregateTypeColumnName,
			payloadColumnName:       payloadColumnName,
		},
		eventDispatcher: eventDispatcher,
		positionChan:    posCh,
	}, posCh, nil
}

// EventHandler implements canal.EventHandler and bridges binlog row events into
// the relay pipeline.
type EventHandler struct {
	canal.DummyEventHandler

	eventMapper     *EventMapper
	eventDispatcher EventDispatcher
	positionChan    chan mysql.Position
}

// OnRow is called by canal for each row-level DML event. Only INSERT actions
// are forwarded; UPDATE/DELETE events are silently skipped.
func (h *EventHandler) OnRow(e *canal.RowsEvent) error {
	logrus.Debug("reading row-event")

	oes, err := h.eventMapper.Map(e)
	if err != nil && errors.Is(err, errNotInsert) {
		logrus.Info("skipping row-event that is not an insert")
		return nil
	}
	if err != nil {
		return err
	}

	for _, oe := range oes {
		if err = h.eventDispatcher.Dispatch(oe); err != nil {
			return err
		}
		logrus.WithField("event", oe).Debug("event dispatched")
	}

	return nil
}

// OnPosSynced forwards the current binlog position to the runner so it can be
// checkpointed to the state store on the next tick.
func (h *EventHandler) OnPosSynced(p mysql.Position, _ mysql.GTIDSet, _ bool) error {
	h.positionChan <- p
	return nil
}

func (h *EventHandler) String() string { return "EventHandler" }
