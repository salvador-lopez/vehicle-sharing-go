package domain

import (
	"vehicle-sharing-go/app/inventory/internal/vehicle/domain/datamodel"
	"vehicle-sharing-go/pkg/domain"
)

func FromDataModel(m *datamodel.Car) *Car {
	return &Car{
		domain.FromModel(m.AggregateRoot),
		&VIN{m.VinNumber},
		m.Color,
	}
}

func ToDataModel(c *Car) *datamodel.Car {
	return &datamodel.Car{
		VinNumber:     c.vin.number,
		Color:         c.color,
		AggregateRoot: domain.ToDataModel(c.AggregateRoot),
	}
}
