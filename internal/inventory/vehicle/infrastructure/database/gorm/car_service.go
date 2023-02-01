package gorm

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"vehicle-sharing-go/internal/inventory/vehicle/application/query/projection"
)

type CarService struct {
	db *gorm.DB
}

func NewCarService(db *gorm.DB) *CarService {
	return &CarService{db: db}
}

func (c *CarService) Find(ctx context.Context, id uuid.UUID) (*projection.Car, error) {
	var carProjection *projection.Car
	result := c.db.WithContext(ctx).Find(&carProjection, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return carProjection, nil
}
