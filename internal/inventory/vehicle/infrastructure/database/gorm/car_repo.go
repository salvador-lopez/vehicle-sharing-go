package gorm

import (
	"context"

	"gorm.io/gorm"

	"vehicle-sharing-go/internal/inventory/vehicle/domain"
)

type CarRepository struct {
	db *gorm.DB
}

func NewCarRepository(db *gorm.DB) *CarRepository {
	return &CarRepository{db: db}
}

func (c CarRepository) Create(ctx context.Context, car *domain.Car) error {
	return nil
}
