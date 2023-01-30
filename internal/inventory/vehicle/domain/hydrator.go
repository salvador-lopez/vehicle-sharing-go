package domain

import (
	"time"

	"github.com/google/uuid"
)

type BaseDTO struct {
	ID        uuid.UUID `gorm:"<-:create;type:varchar(36)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CarDTO struct {
	VIN   string `gorm:"type:varchar(255);unique"`
	Color string `gorm:"type:varchar(255)"`
	*BaseDTO
}

func HydrateCar(carDTO *CarDTO) *Car {
	return &Car{
		carDTO.ID,
		carDTO.CreatedAt,
		carDTO.UpdatedAt,
		&VIN{carDTO.VIN},
		carDTO.Color,
	}
}

func DeHydrateCar(
	car *Car,
) *CarDTO {
	return &CarDTO{
		car.vin.number,
		car.color,
		&BaseDTO{ID: car.id, CreatedAt: car.createdAt, UpdatedAt: car.updatedAt},
	}
}
