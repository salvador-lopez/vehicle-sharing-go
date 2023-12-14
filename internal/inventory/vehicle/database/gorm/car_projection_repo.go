package gorm

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"vehicle-sharing-go/internal/inventory/vehicle/database/gorm/model"
	"vehicle-sharing-go/internal/inventory/vehicle/projection"
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

func (c *CarProjectionRepository) Find(ctx context.Context, id uuid.UUID) (*projection.Car, error) {
	var carModel *model.CarProjection
	result := c.db.WithContext(ctx).Find(&carModel, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return carModel.Car, nil
}
