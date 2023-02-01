package gorm

import (
	"context"

	"gorm.io/gorm"

	"vehicle-sharing-go/internal/inventory/vehicle/application/projection"
)

type CarProjector struct {
	db *gorm.DB
}

func NewCarProjector(db *gorm.DB) *CarProjector {
	return &CarProjector{db: db}
}

func (c *CarProjector) Project(ctx context.Context, car *projection.Car) error {
	return c.db.WithContext(ctx).Create(car).Error
}
