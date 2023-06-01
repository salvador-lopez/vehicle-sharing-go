package command

import (
	"context"
	"time"

	"github.com/google/uuid"

	"vehicle-sharing-go/internal/inventory/vehicle/domain"
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
	evtPublisher EventPublisher
}

func NewCreateCarHandler(
	idGen func() uuid.UUID,
	now func() time.Time,
	carRepo CarRepository,
	evtPublisher EventPublisher,
) *CreateCarHandler {
	return &CreateCarHandler{idGen: idGen, now: now, carRepo: carRepo, evtPublisher: evtPublisher}
}

func (h *CreateCarHandler) Handle(ctx context.Context, cmd *CreateCar) error {
	vin, err := domain.NewVIN(cmd.VIN)
	if err != nil {
		return err
	}
	car := domain.NewCar(cmd.ID, vin, cmd.Color, h.idGen, h.now)

	err = h.carRepo.Create(ctx, car)
	if err != nil {
		return err
	}

	_ = h.evtPublisher.Publish(ctx, "inventory", car)

	return nil
}
