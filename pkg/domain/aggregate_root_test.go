//go:build unit

package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"vehicle-sharing-go/pkg/domain"
)

func TestNewAggregateRoot(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	aggregate := domain.NewAggregateRoot(id, func() time.Time {
		return now
	})

	require.NotNil(t, aggregate)
	require.Equal(t, now, aggregate.CreatedAt())
	require.Equal(t, now, aggregate.UpdatedAt())
	require.Empty(t, aggregate.RecordedEvents())
}
