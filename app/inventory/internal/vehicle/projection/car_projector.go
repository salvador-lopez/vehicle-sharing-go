package projection

import (
	"context"
	"time"

	"github.com/google/uuid"

	"vehicle-sharing-go/app/inventory/internal/vehicle/domain/event"
)

//go:generate mockgen -destination=mock/car_repository_mock.go -package=mock . CarRepository
type CarRepository interface {
	Create(ctx context.Context, car *Car) error
}

//go:generate mockgen -destination=mock/vin_decoder_mock.go -package=mock . VINDecoder
type VINDecoder interface {
	Decode(ctx context.Context, vinNumber string) (*VINData, error)
}

type Car struct {
	ID        uuid.UUID `gorm:"<-:create;type:varchar(36)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	*VINData
	Color string `gorm:"type:varchar(255)"`
}

type VINData struct {
	VIN           string  `gorm:"type:varchar(255);unique"`
	Country       *string `gorm:"type:varchar(255)"`
	Manufacturer  *string `gorm:"type:varchar(255)"`
	Brand         *string `gorm:"type:varchar(255)"`
	EngineSize    *string `gorm:"type:varchar(255)"`
	FuelType      *string `gorm:"type:varchar(255)"`
	Model         *string `gorm:"type:varchar(255)"`
	Year          *string `gorm:"type:varchar(255)"`
	AssemblyPlant *string `gorm:"type:varchar(255)"`
	SN            *string `gorm:"type:varchar(255)"`
}

type CarProjector struct {
	vinDecoder VINDecoder
	carRepo    CarRepository
}

func NewCarProjector(vinDecoder VINDecoder, carRepo CarRepository) *CarProjector {
	return &CarProjector{vinDecoder: vinDecoder, carRepo: carRepo}
}

func (cp *CarProjector) ProjectCarCreated(ctx context.Context, carID uuid.UUID, payload *event.CarCreatedPayload) error {
	vinData, err := cp.vinDecoder.Decode(ctx, payload.VinNumber)
	if err != nil {
		return err
	}

	projection := &Car{
		ID:        carID,
		CreatedAt: payload.CreatedAt,
		UpdatedAt: payload.UpdatedAt,
		VINData:   vinData,
		Color:     payload.Color,
	}

	return cp.carRepo.Create(ctx, projection)
}
