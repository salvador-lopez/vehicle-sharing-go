package gorm

import (
	"context"

	"gorm.io/gorm"

	"vehicle-sharing-go/internal/inventory/vehicle/application/projection"
)

type CarProjectionRepository struct {
	db *gorm.DB
}

func NewCarProjectionRepository(db *gorm.DB) *CarProjectionRepository {
	return &CarProjectionRepository{db: db}
}

func (c *CarProjectionRepository) Create(ctx context.Context, car *projection.Car) error {
	return c.db.WithContext(ctx).Create(car).Error
}
