//go:build integration_inventory

package gorm

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"vehicle-sharing-go/internal/inventory/vehicle/application/projection"
	gormvehicle "vehicle-sharing-go/internal/inventory/vehicle/infrastructure/database/gorm"
)

type carProjectionRepoIntegrationSuite struct {
	databaseSuite
	carId uuid.UUID
	sut   *gormvehicle.CarProjectionRepository
}

func (s *carProjectionRepoIntegrationSuite) SetupSuite() {
	s.databaseSuite.SetupSuite()
	s.initDb()
	s.carId = uuid.New()
	s.sut = gormvehicle.NewCarProjectionRepository(s.db)
}

func (s *carProjectionRepoIntegrationSuite) initDb() {
	s.Require().NoError(s.db.AutoMigrate(&projection.Car{}))
}

func TestCarProjectorIntegrationSuite(t *testing.T) {
	suite.Run(t, new(carProjectionRepoIntegrationSuite))
}

func (s *carProjectionRepoIntegrationSuite) TestProject() {
	carProjectionExpected := &projection.Car{
		ID:        s.carId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		VINData: s.buildVinDataProjection(
			"SCBFR7ZA5CC072256",
			"UNITED KINGDOM (UK)",
			"BENTLEY MOTORS LIMITED",
			"BENTLEY",
			"6L",
			"Flexible Fuel Vehicle (FFV)",
			"Continental",
			"2012",
			"-",
			"411439",
		),
		Color: "Spectral Blue",
	}
	s.Require().NoError(s.sut.Create(s.ctx, carProjectionExpected))

	var carProjection *projection.Car
	s.Require().NoError(s.db.WithContext(s.ctx).Find(&carProjection, s.carId).Error)

	s.requireEqualProjections(carProjectionExpected, carProjection)
}

func (s *carProjectionRepoIntegrationSuite) TestFind() {
	carProjectionExpected := &projection.Car{
		ID:        s.carId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		VINData: s.buildVinDataProjection(
			"AJBFR7ZA5JF072267",
			"SPAIN (ES)",
			"SEAT SA",
			"SEAT",
			"1L",
			"Diesel",
			"ARONA",
			"2013",
			"-",
			"312438",
		),
		Color: "Sapphire Graphite",
	}

	s.Require().NoError(s.db.WithContext(s.ctx).Create(carProjectionExpected).Error)

	var carProjection *projection.Car
	s.Require().NoError(s.db.WithContext(s.ctx).Find(&carProjection, s.carId).Error)

	s.requireEqualProjections(carProjectionExpected, carProjection)
}

func (s *carProjectionRepoIntegrationSuite) TearDownTest() {
	s.db.Delete(&projection.Car{}, s.carId)
	s.databaseSuite.TearDownTest()
}

func (s *carProjectionRepoIntegrationSuite) buildVinDataProjection(
	vinNumber,
	country,
	manufacturer,
	brand,
	engineSize,
	fuelType,
	model,
	year,
	assemblyPlant,
	sn string,
) *projection.VINData {
	return &projection.VINData{
		VIN:           vinNumber,
		Country:       &country,
		Manufacturer:  &manufacturer,
		Brand:         &brand,
		EngineSize:    &engineSize,
		FuelType:      &fuelType,
		Model:         &model,
		Year:          &year,
		AssemblyPlant: &assemblyPlant,
		SN:            &sn,
	}
}

func (s *carProjectionRepoIntegrationSuite) requireEqualProjections(expected *projection.Car, actual *projection.Car) {
	s.Require().Equal(expected.ID, actual.ID)
	requireEqualDates(expected.CreatedAt, actual.CreatedAt, s.Require())
	requireEqualDates(expected.UpdatedAt, actual.UpdatedAt, s.Require())
	s.Require().Equal(expected.VIN, actual.VIN)
	s.Require().Equal(expected.Country, actual.Country)
	s.Require().Equal(expected.Manufacturer, actual.Manufacturer)
	s.Require().Equal(expected.Brand, actual.Brand)
	s.Require().Equal(expected.EngineSize, actual.EngineSize)
	s.Require().Equal(expected.FuelType, actual.FuelType)
	s.Require().Equal(expected.Model, actual.Model)
	s.Require().Equal(expected.Year, actual.Year)
	s.Require().Equal(expected.AssemblyPlant, actual.AssemblyPlant)
	s.Require().Equal(expected.SN, actual.SN)
	s.Require().Equal(expected.Color, actual.Color)
}
