package vehicle

import (
	"context"
	"time"

	"github.com/google/uuid"

	"vehicle-sharing-go/internal/inventory/vehicle/domain"
)

type CreateCarCommand struct {
	ID  uuid.UUID
	VIN string
}

type CreateCarHandler struct {
	nowFun  func() time.Time
	carRepo domain.CarRepository
}

func NewCreateCarHandler(nowFun func() time.Time, carRepo domain.CarRepository) *CreateCarHandler {
	return &CreateCarHandler{nowFun: nowFun, carRepo: carRepo}
}

func (h *CreateCarHandler) Handle(ctx context.Context, cmd *CreateCarCommand) error {
	return h.carRepo.Create(
		ctx,
		domain.NewCar(cmd.ID, cmd.VIN, h.nowFun),
	)
}
