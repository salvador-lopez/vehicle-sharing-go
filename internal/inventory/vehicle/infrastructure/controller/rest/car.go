package rest

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"vehicle-sharing-go/internal/inventory/vehicle/application/command"
	"vehicle-sharing-go/internal/inventory/vehicle/application/projection"
	"vehicle-sharing-go/internal/inventory/vehicle/infrastructure/controller/gen/car"
)

//go:generate mockgen -destination=mock/find_car_query_service_mock.go -package=mock . FindCarQueryService
type FindCarQueryService interface {
	Find(ctx context.Context, id uuid.UUID) (*projection.Car, error)
}

//go:generate mockgen -destination=mock/create_car_command_handler_mock.go -package=mock . CreateCarCommandHandler
type CreateCarCommandHandler interface {
	Handle(ctx context.Context, cmd *command.CreateCar) error
}

type CarController struct {
	commandHandler CreateCarCommandHandler
	queryService   FindCarQueryService
}

func NewCarController(ch CreateCarCommandHandler, qs FindCarQueryService) *CarController {
	return &CarController{commandHandler: ch, queryService: qs}
}

func (v CarController) Get(ctx context.Context, payload *car.GetPayload) (res *car.CarResource, err error) {
	carID, _ := uuid.Parse(payload.ID)
	carProjection, err := v.queryService.Find(ctx, carID)
	if err != nil {
		err = car.MakeInternal(err)
		return
	}
	if carProjection == nil {
		err = car.MakeNotFound(errors.New("car not found"))
		return
	}

	res = &car.CarResource{
		ID:        carProjection.ID.String(),
		CreatedAt: carProjection.CreatedAt.String(),
		UpdatedAt: carProjection.UpdatedAt.String(),
		Color:     carProjection.Color,
		VinData: &car.VinData{
			Vin:           car.Vin(carProjection.VIN),
			Country:       carProjection.Country,
			Manufacturer:  carProjection.Manufacturer,
			Brand:         carProjection.Brand,
			EngineSize:    carProjection.EngineSize,
			FuelType:      carProjection.FuelType,
			Model:         carProjection.Model,
			Year:          carProjection.Year,
			AssemblyPlant: carProjection.AssemblyPlant,
			SN:            carProjection.SN,
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
