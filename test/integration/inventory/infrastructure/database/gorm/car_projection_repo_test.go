//go:build integration_inventory

package gorm

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"vehicle-sharing-go/internal/inventory/vehicle/application/projection"
	gormvehicle "vehicle-sharing-go/internal/inventory/vehicle/infrastructure/database/gorm"
)

type carProjectionRepoIntegrationSuite struct {
	carProjectionSuite
	sut *gormvehicle.CarProjectionRepository
}

func (s *carProjectionRepoIntegrationSuite) SetupSuite() {
	s.carProjectionSuite.SetupSuite()
	s.sut = gormvehicle.NewCarProjectionRepository(s.db)
}

func TestCarProjectorIntegrationSuite(t *testing.T) {
	suite.Run(t, new(carProjectionRepoIntegrationSuite))
}

func (s *carProjectionRepoIntegrationSuite) TestProject() {
	carProjectionExpected := &projection.Car{
		ID:        s.carId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		VINData: &projection.VINData{
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
		},
		Color: "Spectral Blue",
	}
	s.Require().NoError(s.sut.Create(s.ctx, carProjectionExpected))

	var carProjection *projection.Car
	s.Require().NoError(s.db.WithContext(s.ctx).Find(&carProjection, s.carId).Error)

	s.requireEqualProjections(carProjectionExpected, carProjection)
}
