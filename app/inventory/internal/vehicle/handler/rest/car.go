package rest

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"vehicle-sharing-go/app/inventory/internal/vehicle/command"
	"vehicle-sharing-go/app/inventory/internal/vehicle/handler/rest/gen/car"
	"vehicle-sharing-go/app/inventory/internal/vehicle/projection"
)

//go:generate mockgen -destination=mock/find_car_query_service_mock.go -package=mock . FindCarQueryService
type FindCarQueryService interface {
	Find(ctx context.Context, id uuid.UUID) (*projection.Car, error)
}

//go:generate mockgen -destination=mock/create_car_command_handler_mock.go -package=mock . CreateCarCommandHandler
type CreateCarCommandHandler interface {
	Handle(ctx context.Context, cmd *command.CreateCar) error
}

type CarHandler struct {
	commandHandler CreateCarCommandHandler
	queryService   FindCarQueryService
}

func NewCarHandler(ch CreateCarCommandHandler, qs FindCarQueryService) *CarHandler {
	return &CarHandler{commandHandler: ch, queryService: qs}
}

func (v CarHandler) Get(ctx context.Context, payload *car.GetPayload) (res *car.CarResource, err error) {
	carID, _ := uuid.Parse(payload.ID)
	carProjection, err := v.queryService.Find(ctx, carID)
	if err != nil {
		err = car.MakeInternal(ErrInternal)
		return
	}
	if carProjection == nil {
		err = car.MakeNotFound(ErrNotFound)
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

func (v CarHandler) Create(ctx context.Context, payload *car.CreatePayload) (err error) {
	carID, _ := uuid.Parse(payload.ID) // We can omit the error handling because uuid format constraint is defined in the api goa design file
	err = v.commandHandler.Handle(ctx, &command.CreateCar{
		ID:    carID,
		VIN:   string(payload.Vin),
		Color: payload.Color,
	})
	if err != nil {
		if errors.Is(err, command.ErrCarAlreadyExists) {
			err = car.MakeConflict(err)
			return
		}

		err = car.MakeInternal(ErrInternal)
	}

	return
}
