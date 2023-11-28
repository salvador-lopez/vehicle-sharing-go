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
	idGen       func() uuid.UUID
	now         func() time.Time
	repoFactory RepositoryFactory
	txSession   TransactionalSession
}

func NewCreateCarHandler(idGen func() uuid.UUID, now func() time.Time, repoFactory RepositoryFactory, txSession TransactionalSession) *CreateCarHandler {
	return &CreateCarHandler{idGen: idGen, now: now, repoFactory: repoFactory, txSession: txSession}
}

func (h *CreateCarHandler) Handle(ctx context.Context, cmd *CreateCar) error {
	vin, err := domain.NewVIN(cmd.VIN)
	if err != nil {
		return err
	}
	car := domain.NewCar(cmd.ID, vin, cmd.Color, h.idGen, h.now)

	err = h.repoFactory.CarRepository().Create(ctx, car)
	if err != nil {
		return err
	}

	// err = h.repoFactory.OutboxRepository().Append(ctx, car)
	// if err != nil {
	// 	return err
	// }

	return nil
}
