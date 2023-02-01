//go:build integration_inventory

package gorm

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"vehicle-sharing-go/internal/inventory/vehicle/application/projection"
	gormvehicle "vehicle-sharing-go/internal/inventory/vehicle/infrastructure/database/gorm"
)

type carServiceIntegrationSuite struct {
	carProjectionSuite
	sut *gormvehicle.CarService
}

func (s *carServiceIntegrationSuite) SetupSuite() {
	s.carProjectionSuite.SetupSuite()
	s.sut = gormvehicle.NewCarService(s.db)
}

func TestCarServiceIntegrationSuite(t *testing.T) {
	suite.Run(t, new(carServiceIntegrationSuite))
}

func (s *carServiceIntegrationSuite) TestFind() {
	carProjectionExpected := &projection.Car{
		ID:            s.carId,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		VIN:           "SCBFR7ZA5CC072256",
		Country:       "UNITED KINGDOM (UK)",
		Manufacturer:  "BENTLEY MOTORS LIMITED",
		Brand:         "BENTLEY",
		EngineSize:    "6L",
		FuelType:      "Flexible Fuel Vehicle (FFV)",
		Model:         "Continental",
		Year:          "2012",
		AssemblyPlant: "-",
		SN:            "411439",
		Color:         "Blue Spectral",
	}
	s.Require().NoError(s.db.WithContext(s.ctx).Create(carProjectionExpected).Error)

	carProjection, err := s.sut.Find(s.ctx, s.carId)
	s.Require().NoError(err)
	s.requireEqualProjections(carProjectionExpected, carProjection)
}
