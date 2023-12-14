package gorm

import (
	"context"

	"vehicle-sharing-go/internal/inventory/vehicle/command"
	"vehicle-sharing-go/internal/inventory/vehicle/database/gorm/model"
	"vehicle-sharing-go/internal/inventory/vehicle/domain"
	"vehicle-sharing-go/pkg/database/gorm"
)

type CarRepository struct {
	conn *gorm.Connection
}

func NewCarRepository(conn *gorm.Connection) *CarRepository {
	return &CarRepository{conn: conn}
}

func (c *CarRepository) Create(ctx context.Context, car *domain.Car) error {
	carModel := &model.Car{Car: car.ToModel()}

	err := c.conn.Db().WithContext(ctx).Create(carModel).Error

	if err != nil {
		if c.conn.IsDuplicateEntryErr(err) {
			return command.ErrCarAlreadyExists
		}

		return err
	}

	return nil
}
