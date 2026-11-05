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
	aggregate := newTestAggregateRoot(id, createdAt, updatedAt)

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
	expected := newTestAggregateRoot(model.ID, createdAt, updatedAt)
	aggregate := domain.FromDataModel(model)

	require.Equal(t, expected, aggregate)
}

func newTestAggregateRoot(id uuid.UUID, createdAt, updatedAt time.Time) *domain.AggregateRoot {
	timestamps := []time.Time{createdAt, updatedAt}
	next := 0
	return domain.NewAggregateRoot(id, func() time.Time {
		timestamp := timestamps[next]
		next++
		return timestamp
	})
}
