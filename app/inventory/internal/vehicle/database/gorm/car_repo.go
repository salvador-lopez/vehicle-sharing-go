package gorm

import (
	"context"

	"vehicle-sharing-go/app/inventory/internal/vehicle/command"
	"vehicle-sharing-go/app/inventory/internal/vehicle/database/gorm/model"
	"vehicle-sharing-go/app/inventory/internal/vehicle/domain"

	gormpkg "vehicle-sharing-go/pkg/database/gorm"
)

type CarRepository struct {
	conn *gormpkg.Connection
}

func NewCarRepository(conn *gormpkg.Connection) *CarRepository {
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
