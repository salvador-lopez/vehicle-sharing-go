//go:build unit

package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"vehicle-sharing-go/pkg/domain"
	"vehicle-sharing-go/pkg/domain/datamodel"
)

func TestBaseEntityCreatedAt(t *testing.T) {
	createdAt := time.Now()
	aggregate := domain.FromModel(&datamodel.AggregateRoot{
		CreatedAt: createdAt,
	})
	entity := aggregate.BaseEntity

	require.Equal(t, createdAt, entity.CreatedAt())
}

func TestBaseEntityUpdatedAt(t *testing.T) {
	updatedAt := time.Now()
	aggregate := domain.FromModel(&datamodel.AggregateRoot{
		UpdatedAt: updatedAt,
	})
	entity := aggregate.BaseEntity

	require.Equal(t, updatedAt, entity.UpdatedAt())
}
