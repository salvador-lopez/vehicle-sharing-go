package command

import (
	"context"
	"time"

	"github.com/google/uuid"
	"vehicle-sharing-go/app/inventory/internal/vehicle/domain"
	"vehicle-sharing-go/pkg/domain/event"
)

type CreateCar struct {
	ID    uuid.UUID
	VIN   string
	Color string
}

type CreateCarHandler struct {
	idGen        func() uuid.UUID
	now          func() time.Time
	carRepo      CarRepository
	txSession    TransactionalSession
	evtPublisher *event.AgRootEventPublisher
}

func NewCreateCarHandler(
	idGen func() uuid.UUID, now func() time.Time,
	cr CarRepository,
	txSession TransactionalSession,
	ep *event.AgRootEventPublisher,
) *CreateCarHandler {
	return &CreateCarHandler{idGen: idGen, now: now, carRepo: cr, txSession: txSession, evtPublisher: ep}
}

func (h *CreateCarHandler) Handle(ctx context.Context, cmd *CreateCar) error {
	vin, err := domain.NewVIN(cmd.VIN)
	if err != nil {
		return err
	}
	car := domain.NewCar(cmd.ID, vin, cmd.Color, h.idGen, h.now)

	return h.txSession.Transaction(ctx, func(ctx context.Context) error {
		err = h.carRepo.Create(ctx, car)
		if err != nil {
			return err
		}

		return h.evtPublisher.Publish(ctx, car)
	})
}
