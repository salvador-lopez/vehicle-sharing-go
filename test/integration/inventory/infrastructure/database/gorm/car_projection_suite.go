package gorm

import (
	"github.com/google/uuid"

	"vehicle-sharing-go/internal/inventory/vehicle/application/projection"
)

type carProjectionSuite struct {
	databaseSuite
	carId uuid.UUID
}

func (s *carProjectionSuite) SetupSuite() {
	s.databaseSuite.SetupSuite()
	s.initDb()
	s.carId = uuid.New()
}

func (s *carProjectionSuite) initDb() {
	s.Require().NoError(s.db.AutoMigrate(&projection.Car{}))
}

func (s *carProjectionSuite) TearDownTest() {
	s.db.Delete(&projection.Car{}, s.carId)
	s.databaseSuite.TearDownTest()
}

func (s *carProjectionSuite) requireEqualProjections(expected *projection.Car, actual *projection.Car) {
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
}
