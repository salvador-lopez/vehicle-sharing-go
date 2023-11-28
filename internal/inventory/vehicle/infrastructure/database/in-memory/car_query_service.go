package in_memory

import (
	"context"
	"time"

	"github.com/google/uuid"

	"vehicle-sharing-go/internal/inventory/vehicle/application/projection"
)

type CarQueryService struct {
	projections map[string]*projection.Car
}

func NewCarQueryService() *CarQueryService {
	carID, _ := uuid.Parse("96194205-a21b-4cb6-b499-74cb1da1a20a")
	country := "United States of America"
	year := "2017"
	return &CarQueryService{projections: map[string]*projection.Car{carID.String(): {
		ID:        carID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		VINData: &projection.VINData{
			VIN:     "4Y1SL65848Z411439",
			Country: &country,
			Year:    &year,
		},
		Color: "Spectral Blue",
	}}}
}

func (c CarQueryService) Find(ctx context.Context, id uuid.UUID) (*projection.Car, error) {
	return c.projections[id.String()], nil
}
