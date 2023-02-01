//go:build integration_inventory

package gorm

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"vehicle-sharing-go/internal/inventory/vehicle/application/query/projection"
	gormvehicle "vehicle-sharing-go/internal/inventory/vehicle/infrastructure/database/gorm"
)

type carServiceIntegrationSuite struct {
	suite.Suite
	ctx       context.Context
	cancelFun context.CancelFunc
	db        *gorm.DB
	carId     uuid.UUID
	sut       *gormvehicle.CarService
}

func (s *carServiceIntegrationSuite) SetupSuite() {
	s.initDb()
	s.carId = uuid.New()
	s.sut = gormvehicle.NewCarService(s.db)
}

func (s *carServiceIntegrationSuite) initDb() {
	s.db = createDb(s.Require())
	s.Require().NoError(s.db.AutoMigrate(&projection.Car{}))
}

func (s *carServiceIntegrationSuite) SetupTest() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	s.ctx = ctx
	s.cancelFun = cancel
}

func (s *carServiceIntegrationSuite) TearDownTest() {
	s.db.Delete(&projection.Car{}, s.carId)
	s.cancelFun()
}

func (s *carServiceIntegrationSuite) TearDownSuite() {
	sqlDb, err := s.db.DB()
	s.Require().NoError(err)

	s.Require().NoError(sqlDb.Close())
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
	}
	s.Require().NoError(s.db.WithContext(s.ctx).Create(carProjectionExpected).Error)

	carProjection, err := s.sut.Find(s.ctx, s.carId)
	s.Require().NoError(err)
	s.requireEqualProjections(carProjectionExpected, carProjection)

}

func (s *carServiceIntegrationSuite) requireEqualProjections(expected *projection.Car, actual *projection.Car) {
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
