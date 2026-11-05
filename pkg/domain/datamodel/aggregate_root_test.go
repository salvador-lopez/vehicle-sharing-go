//go:build unit

package datamodel_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"vehicle-sharing-go/pkg/domain/datamodel"
)

func TestAggregateRootHasNoRecordedEvents(t *testing.T) {
	_, exists := reflect.TypeOf(datamodel.AggregateRoot{}).FieldByName("RecordedEvents")

	require.False(t, exists)
}
