package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -destination=mock/car_repository_mock.go -package=mock . CarRepository
type CarRepository interface {
	Create(context.Context, *Car) error
}

var ErrInvalidVin = errors.New("invalid vin provided")

//go:generate mockgen -destination=mock/vin_validator_mock.go -package=mock . VinValidator
type VinValidator interface {
	Validate(number string) error
}

type VIN struct {
	number string
}

func NewVIN(number string, validator VinValidator) (*VIN, error) {
	err := validator.Validate(number)
	if err != nil {
		return nil, err
	}
	return &VIN{number: number}, nil
}

type Car struct {
	id        uuid.UUID
	createdAt time.Time
	updatedAt time.Time
	vin       *VIN
	color     string
}

func NewCar(
	id uuid.UUID,
	vin *VIN,
	color string,
	nowFun func() time.Time,
) *Car {
	now := nowFun()
	return &Car{id, now, now, vin, color}
}
