package query

import (
	"context"

	"github.com/google/uuid"

	"vehicle-sharing-go/internal/inventory/vehicle/application/projection"
)

//go:generate mockgen -destination=mock/car_service_mock.go -package=mock . CarService
type CarService interface {
	Find(ctx context.Context, id uuid.UUID) (*projection.Car, error)
}
