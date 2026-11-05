//go:build unit

package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"vehicle-sharing-go/pkg/domain"
	"vehicle-sharing-go/pkg/domain/datamodel"
)

func TestToDataModel(t *testing.T) {
	id := uuid.New()
	createdAt := time.Now()
	updatedAt := createdAt.Add(time.Second)
	timestamps := []time.Time{createdAt, updatedAt}
	nextTimestamp := 0
	aggregate := domain.NewAggregateRoot(id, func() time.Time {
		timestamp := timestamps[nextTimestamp]
		nextTimestamp++
		return timestamp
	})

	result := domain.ToDataModel(aggregate)

	require.Equal(t, &datamodel.AggregateRoot{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, result)
}

func TestFromDataModel(t *testing.T) {
	createdAt := time.Now()
	updatedAt := createdAt.Add(time.Second)
	model := &datamodel.AggregateRoot{
		ID:        uuid.New(),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	timestamps := []time.Time{createdAt, updatedAt}
	nextTimestamp := 0
	expected := domain.NewAggregateRoot(model.ID, func() time.Time {
		timestamp := timestamps[nextTimestamp]
		nextTimestamp++
		return timestamp
	})
	aggregate := domain.FromDataModel(model)

	require.Equal(t, expected, aggregate)
}
