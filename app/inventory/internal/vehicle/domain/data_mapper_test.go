//go:build unit

package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"vehicle-sharing-go/app/inventory/internal/vehicle/domain/datamodel"
	pkgdomain "vehicle-sharing-go/pkg/domain"
	modelpkg "vehicle-sharing-go/pkg/domain/datamodel"
)

func TestFromDataModel(t *testing.T) {
	id := uuid.New()
	createdAt := time.Now()
	updatedAt := createdAt.Add(time.Minute)
	model := &datamodel.Car{
		VinNumber: "4Y1SL65848Z411439",
		Color:     "Spectral Blue",
		AggregateRoot: &modelpkg.AggregateRoot{
			ID:        id,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
	}
	timestamps := []time.Time{createdAt, updatedAt}
	nextTimestamp := 0
	expected := &Car{
		AggregateRoot: pkgdomain.NewAggregateRoot(id, func() time.Time {
			timestamp := timestamps[nextTimestamp]
			nextTimestamp++
			return timestamp
		}),
		vin:   &VIN{number: "4Y1SL65848Z411439"},
		color: "Spectral Blue",
	}

	car := FromDataModel(model)

	require.Equal(t, expected, car)
}

func TestToDataModel(t *testing.T) {
	id := uuid.New()
	createdAt := time.Now()
	updatedAt := createdAt.Add(time.Minute)
	timestamps := []time.Time{createdAt, updatedAt}
	nextTimestamp := 0
	car := &Car{
		AggregateRoot: pkgdomain.NewAggregateRoot(id, func() time.Time {
			timestamp := timestamps[nextTimestamp]
			nextTimestamp++
			return timestamp
		}),
		vin:   &VIN{number: "4Y1SL65848Z411439"},
		color: "Spectral Blue",
	}
	car.AggregateRoot.RecordEvent(uuid.New(), "CarCreatedEvent", "Car", nil, createdAt)

	result := ToDataModel(car)

	require.Equal(t, &datamodel.Car{
		VinNumber: "4Y1SL65848Z411439",
		Color:     "Spectral Blue",
		AggregateRoot: &modelpkg.AggregateRoot{
			ID:        id,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
	}, result)
}
