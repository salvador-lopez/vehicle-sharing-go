package vehicle

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CreateCarCommand struct {
	ID               uuid.UUID
	VIN              string
	BrandName        string
	BrandModel       string
	Color            string
	EngineType       string
	TransmissionType string
}

type CreateCarHandler struct {
	nowFun  func() time.Time
	carRepo CarRepository
}

func NewCreateCarHandler(nowFun func() time.Time, carRepo CarRepository) *CreateCarHandler {
	return &CreateCarHandler{nowFun: nowFun, carRepo: carRepo}
}

func (h *CreateCarHandler) Handle(ctx context.Context, cmd *CreateCarCommand) error {
	return h.carRepo.Create(ctx, &Car{
		id:               cmd.ID,
		createdAt:        h.nowFun(),
		updatedAt:        h.nowFun(),
		vin:              cmd.VIN,
		brandName:        cmd.BrandName,
		brandModel:       cmd.BrandModel,
		color:            cmd.Color,
		engineType:       cmd.EngineType,
		transmissionType: cmd.TransmissionType,
	})
}
