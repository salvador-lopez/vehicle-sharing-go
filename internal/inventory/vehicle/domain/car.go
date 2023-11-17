package domain

import (
	"time"

	"github.com/google/uuid"

	"vehicle-sharing-go/internal/inventory/vehicle/domain/event"
	"vehicle-sharing-go/internal/inventory/vehicle/domain/model"
	"vehicle-sharing-go/pkg/domain"
)

type Car struct {
	*domain.AggregateRoot
	vin   *VIN
	color string
}

func NewCar(
	id uuid.UUID,
	vin *VIN,
	color string,
	evtIdGen func() uuid.UUID,
	now func() time.Time,
) *Car {
	car := &Car{domain.NewAggregateRoot(id, now), vin, color}
	car.recordCreatedEvent(evtIdGen(), now())

	return car
}

func (c *Car) recordCreatedEvent(evtID uuid.UUID, timestamp time.Time) {
	c.RecordEvent(evtID, "CarCreatedEvent", "Car", &event.CarCreatedPayload{
		VinNumber: c.vin.number,
		Color:     c.color,
		CreatedAt: c.CreatedAt(),
		UpdatedAt: c.UpdatedAt(),
	}, timestamp)
}

func CarFromModel(model *model.Car) *Car {
	return &Car{
		domain.AggregateRootFromModel(model.AggregateRoot),
		&VIN{model.VinNumber},
		model.Color,
	}
}

func (c *Car) ToModel() *model.Car {
	return &model.Car{
		VinNumber:     c.vin.number,
		Color:         c.color,
		AggregateRoot: c.AggregateRoot.ToDataModel(),
	}
}
