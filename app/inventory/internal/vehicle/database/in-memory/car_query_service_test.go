//go:build unit

package in_memory_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	in_memory "vehicle-sharing-go/app/inventory/internal/vehicle/database/in-memory"
	"vehicle-sharing-go/app/inventory/internal/vehicle/projection"
)

func TestCarQueryServiceFind(t *testing.T) {
	service := in_memory.NewCarQueryService()
	carID := uuid.MustParse("96194205-a21b-4cb6-b499-74cb1da1a20a")

	car, err := service.Find(context.Background(), carID)

	require.NoError(t, err)
	require.Equal(t, &projection.Car{
		ID:        carID,
		CreatedAt: car.CreatedAt,
		UpdatedAt: car.UpdatedAt,
		VINData: &projection.VINData{
			VIN:     "4Y1SL65848Z411439",
			Country: stringPointer("United States of America"),
			Year:    stringPointer("2017"),
		},
		Color: "Spectral Blue",
	}, car)
}

func TestCarQueryServiceFindMissingCar(t *testing.T) {
	service := in_memory.NewCarQueryService()

	car, err := service.Find(context.Background(), uuid.New())

	require.NoError(t, err)
	require.Nil(t, car)
}

func stringPointer(value string) *string {
	return &value
}
