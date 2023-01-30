package domain

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -destination=mock/car_repository_mock.go -package=mock . CarRepository
type CarRepository interface {
	Create(context.Context, *Car) error
}

type VIN struct {
	number string
}

func NewVIN(number string) (*VIN, error) {
	err := guardVIN(number)
	if err != nil {
		return nil, err
	}
	return &VIN{number: number}, nil
}

var ErrInvalidVin = errors.New("invalid vin provided")

func guardVIN(number string) error {
	matches, _ := regexp.Match("^[A-HJ-NPR-Z\\d]{8}[\\dX][A-HJ-NPR-Z\\d]{8}$", []byte(number))
	if !matches {
		return fmt.Errorf("%v: %s", ErrInvalidVin, number)
	}

	return nil
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
