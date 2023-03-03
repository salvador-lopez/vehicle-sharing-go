package rest

import (
	"context"

	"github.com/google/uuid"

	"vehicle-sharing-go/internal/inventory/vehicle/application/command"
	"vehicle-sharing-go/internal/inventory/vehicle/infrastructure/controller/gen/car"
)

type CarController struct {
	commandHandler *command.CreateCarHandler
}

func NewCarController(commandHandler *command.CreateCarHandler) *CarController {
	return &CarController{commandHandler: commandHandler}
}

func (v CarController) Create(ctx context.Context, payload *car.CreatePayload) (err error) {
	carID, _ := uuid.Parse(payload.ID) // We can omit the error handling because uuid format constraint is defined in the api goa design file
	err = v.commandHandler.Handle(ctx, &command.CreateCar{
		ID:    carID,
		VIN:   payload.Vin,
		Color: payload.Color,
	})
	if err != nil {
		return car.MakeInternal(err)
	}

	return err
}
