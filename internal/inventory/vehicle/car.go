package vehicle

import (
	"context"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -destination=mock/car_repository_mock.go -package=mock . CarRepository
type CarRepository interface {
	Create(context.Context, *Car) error
}

type Car struct {
	id               uuid.UUID
	createdAt        time.Time
	updatedAt        time.Time
	vin              string
	brandName        string
	brandModel       string
	color            string
	engineType       string
	transmissionType string
}

func HydrateCar(
	id uuid.UUID,
	createdAt,
	updatedAt time.Time,
	vin,
	brandName,
	brandModel,
	color,
	engineType,
	transmissionType string,
) *Car {
	return &Car{
		id,
		createdAt,
		updatedAt,
		vin,
		brandName,
		brandModel,
		color,
		engineType,
		transmissionType,
	}
}
