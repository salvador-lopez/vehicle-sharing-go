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
	nowFun       func() time.Time
	vinValidator domain.VinValidator
	carRepo      domain.CarRepository
}

func NewCreateCarHandler(nowFun func() time.Time, vinValidator domain.VinValidator, carRepo domain.CarRepository) *CreateCarHandler {
	return &CreateCarHandler{nowFun: nowFun, vinValidator: vinValidator, carRepo: carRepo}
}

func (h *CreateCarHandler) Handle(ctx context.Context, cmd *CreateCar) error {
	vin, _ := domain.NewVIN(cmd.VIN, h.vinValidator)
	car := domain.NewCar(cmd.ID, vin, cmd.Color, h.nowFun)

	_ = h.carRepo.Create(ctx, car)

	return nil
}
