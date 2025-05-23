package command

import (
	"context"
	"errors"
	"vehicle-sharing-go/app/inventory/internal/vehicle/domain"
)

//go:generate mockgen -destination=mock/transactional_session_mock.go -package=mock . TransactionalSession
type TransactionalSession interface {
	Transaction(ctx context.Context, f func(context.Context) error) error
}

var ErrCarAlreadyExists = errors.New("car already exist")

//go:generate mockgen -destination=mock/car_repository_mock.go -package=mock . CarRepository
type CarRepository interface {
	Create(context.Context, *domain.Car) error
}
