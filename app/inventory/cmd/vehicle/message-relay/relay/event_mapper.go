package relay

import (
	"errors"
	"fmt"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/schema"
)

// EventMapper converts canal RowsEvents into []OutboxEvent values.
type EventMapper struct {
	aggregateIDColumnName   string
	aggregateTypeColumnName string
	payloadColumnName       string
}

var errNotInsert = errors.New("row-event is not an insert")

// Map returns an OutboxEvent for every row in the given RowsEvent.
// Only INSERT events are accepted; any other action returns errNotInsert.
func (e *EventMapper) Map(event *canal.RowsEvent) ([]OutboxEvent, error) {
	if event.Action != canal.InsertAction {
		return nil, errNotInsert
	}

	if err := assertRowSizesValid(event); err != nil {
		return nil, err
	}

	oes := make([]OutboxEvent, 0, len(event.Rows))
	for _, row := range event.Rows {
		cols := buildColumns(event.Table.Columns, row)

		aggregateID, aggregateType, payload, err := e.mainColumnValues(cols)
		if err != nil {
			return nil, err
		}

		oes = append(oes, OutboxEvent{
			AggregateID:                aggregateID,
			AggregateType:              aggregateType,
			Payload:                    payload,
			Columns:                    cols,
			EventTimestampFromDatabase: event.Header.Timestamp,
		})
	}

	return oes, nil
}

func assertRowSizesValid(event *canal.RowsEvent) error {
	for _, row := range event.Rows {
		if len(event.Table.Columns) != len(row) {
			return errors.New("unexpected row length: table column count does not match row value count")
		}
	}
	return nil
}

func buildColumns(tableColumns []schema.TableColumn, rowValues []interface{}) []Column {
	cols := make([]Column, 0, len(tableColumns))
	for i, tc := range tableColumns {
		if rowValues[i] == nil {
			cols = append(cols, Column{Name: []byte(tc.Name)})
			continue
		}

		cv, ok := rowValues[i].([]byte)
		if !ok {
			cv = []byte(fmt.Sprintf("%v", rowValues[i]))
		}

		cols = append(cols, Column{Name: []byte(tc.Name), Value: cv})
	}
	return cols
}

func (e *EventMapper) mainColumnValues(cols []Column) (aggregateID, aggregateType, payload []byte, err error) {
	for _, c := range cols {
		if c.Name == nil {
			continue
		}
		switch string(c.Name) {
		case e.aggregateIDColumnName:
			aggregateID = c.Value
		case e.aggregateTypeColumnName:
			aggregateType = c.Value
		case e.payloadColumnName:
			payload = c.Value
		}
	}

	if aggregateID == nil {
		return nil, nil, nil, fmt.Errorf("column %q not found in row", e.aggregateIDColumnName)
	}
	if aggregateType == nil {
		return nil, nil, nil, fmt.Errorf("column %q not found in row", e.aggregateTypeColumnName)
	}
	if payload == nil {
		return nil, nil, nil, fmt.Errorf("column %q not found in row", e.payloadColumnName)
	}

	return aggregateID, aggregateType, payload, nil
}
