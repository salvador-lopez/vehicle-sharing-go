package model

import (
	"vehicle-sharing-go/app/inventory/internal/vehicle/domain/model"
)

type Car struct {
	*model.Car
}

// TableName This is not needed but added in order to exemplify why we are composing the gorm model.Car with the domain model.Car
func (c *Car) TableName() string {
	return "cars"
}
