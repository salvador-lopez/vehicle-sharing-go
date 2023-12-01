package command

import (
	"context"

	"vehicle-sharing-go/internal/inventory/vehicle/domain"
	"vehicle-sharing-go/pkg/domain/event"
)

//go:generate mockgen -destination=mock/transactional_session_mock.go -package=mock . TransactionalSession
type TransactionalSession interface {
	Transaction(ctx context.Context, f func(context.Context) error) error
}

//go:generate mockgen -destination=mock/car_repository_mock.go -package=mock . CarRepository
type CarRepository interface {
	Create(context.Context, *domain.Car) error
}

//go:generate mockgen -destination=mock/domain_events_recorder_mock.go -package=mock . DomainEventsRecorder
type DomainEventsRecorder interface {
	RecordedEvents() []*event.Event
}
