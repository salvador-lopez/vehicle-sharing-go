package command

import (
	"context"

	"vehicle-sharing-go/internal/inventory/vehicle/domain"
)

//go:generate mockgen -destination=mock/car_repository_mock.go -package=mock . CarRepository
type CarRepository interface {
	Create(context.Context, *domain.Car) error
}
