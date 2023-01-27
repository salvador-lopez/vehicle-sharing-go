package domain

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
	id        uuid.UUID
	createdAt time.Time
	updatedAt time.Time
	vin       string
}

func NewCar(
	id uuid.UUID,
	vin string,
	nowFun func() time.Time,
) *Car {
	now := nowFun()
	return &Car{id, now, now, vin}
}

func HydrateCar(
	id uuid.UUID,
	createdAt,
	updatedAt time.Time,
	vin string,
) *Car {
	return &Car{id, createdAt, updatedAt, vin}
}
