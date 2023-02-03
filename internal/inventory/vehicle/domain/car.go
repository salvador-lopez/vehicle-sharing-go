package domain

import (
	"context"
	"time"

	"github.com/google/uuid"

	"vehicle-sharing-go/pkg/domain"
)

//go:generate mockgen -destination=mock/car_repository_mock.go -package=mock . CarRepository
type CarRepository interface {
	Create(context.Context, *Car) error
}

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
	c.RecordEvent(evtID, "CarCreatedEvent", "Car", &CarCreatedEventPayload{
		c.vin,
		c.color,
		c.CreatedAt(),
		c.UpdatedAt(),
	}, timestamp)
}

type CarDTO struct {
	VinNumber string `gorm:"type:varchar(255);unique"`
	Color     string `gorm:"type:varchar(255)"`
	*domain.AgRootDTO
}

func (c *CarDTO) ToAggRoot() *Car {
	return &Car{
		c.AgRootDTO.ToAggRoot(),
		&VIN{c.VinNumber},
		c.Color,
	}
}

func (c *Car) ToDTO() *CarDTO {
	return &CarDTO{
		c.vin.number,
		c.color,
		c.AggregateRoot.ToDTO(),
	}
}
