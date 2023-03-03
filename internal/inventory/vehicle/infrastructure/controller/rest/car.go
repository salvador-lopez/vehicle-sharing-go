package rest

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"vehicle-sharing-go/internal/inventory/vehicle/application/command"
	"vehicle-sharing-go/internal/inventory/vehicle/application/query"
	"vehicle-sharing-go/internal/inventory/vehicle/infrastructure/controller/gen/car"
)

type CarController struct {
	commandHandler *command.CreateCarHandler
	queryService   query.CarService
}

func NewCarController(commandHandler *command.CreateCarHandler, queryService query.CarService) *CarController {
	return &CarController{commandHandler: commandHandler, queryService: queryService}
}

func (v CarController) Get(ctx context.Context, payload *car.GetPayload) (res *car.CarResource, err error) {
	carID, _ := uuid.Parse(payload.ID)
	carDTO, err := v.queryService.Find(ctx, carID)
	if err != nil {
		err = car.MakeInternal(err)
		return
	}
	if carDTO == nil {
		err = car.MakeNotFound(errors.New("car not found"))
		return
	}

	res = &car.CarResource{
		ID:        carDTO.ID.String(),
		CreatedAt: carDTO.CreatedAt.String(),
		UpdatedAt: carDTO.UpdatedAt.String(),
		Color:     carDTO.Color,
		VinData: &car.VinData{
			Vin:           car.Vin(carDTO.VIN),
			Country:       carDTO.Country,
			Manufacturer:  carDTO.Manufacturer,
			Brand:         carDTO.Brand,
			EngineSize:    carDTO.EngineSize,
			FuelType:      carDTO.FuelType,
			Model:         carDTO.Model,
			Year:          carDTO.Year,
			AssemblyPlant: carDTO.AssemblyPlant,
			SN:            carDTO.SN,
		},
	}

	return
}

func (v CarController) Create(ctx context.Context, payload *car.CreatePayload) (err error) {
	carID, _ := uuid.Parse(payload.ID) // We can omit the error handling because uuid format constraint is defined in the api goa design file
	err = v.commandHandler.Handle(ctx, &command.CreateCar{
		ID:    carID,
		VIN:   string(payload.Vin),
		Color: payload.Color,
	})
	if err != nil {
		err = car.MakeInternal(err)
	}

	return
}
