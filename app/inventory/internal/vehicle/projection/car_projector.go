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
	ID        uuid.UUID `gorm:"<-:create;type:varchar(36)" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	*VINData
	Color string `gorm:"type:varchar(255)"`
}

type VINData struct {
	VIN           string  `gorm:"type:varchar(255);unique" json:"vin"`
	Country       *string `gorm:"type:varchar(255)" json:"country"`
	Manufacturer  *string `gorm:"type:varchar(255)" json:"manufacturer"`
	Brand         *string `gorm:"type:varchar(255)" json:"brand"`
	EngineSize    *string `gorm:"type:varchar(255)" json:"engineSize"`
	FuelType      *string `gorm:"type:varchar(255)" json:"fuelType"`
	Model         *string `gorm:"type:varchar(255)" json:"model"`
	Year          *string `gorm:"type:varchar(255)" json:"year"`
	AssemblyPlant *string `gorm:"type:varchar(255)" json:"assemblyPlant"`
	SN            *string `gorm:"type:varchar(255)" json:"sn"`
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
