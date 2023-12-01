//go:build integration_inventory

package gorm

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"vehicle-sharing-go/internal/inventory/vehicle/domain"
	domainmodel "vehicle-sharing-go/internal/inventory/vehicle/domain/model"
	gormvehicle "vehicle-sharing-go/internal/inventory/vehicle/infrastructure/database/gorm"
	"vehicle-sharing-go/internal/inventory/vehicle/infrastructure/database/gorm/model"
	domainmodelpkg "vehicle-sharing-go/pkg/domain/model"
)

type carRepoIntegrationSuite struct {
	databaseSuite
	carId uuid.UUID
	sut   *gormvehicle.CarRepository
}

func (s *carRepoIntegrationSuite) SetupSuite() {
	s.databaseSuite.SetupSuite()
	s.initDb()
	s.carId = uuid.New()
	s.sut = gormvehicle.NewCarRepository(s.conn)
}

func (s *carRepoIntegrationSuite) initDb() {
	s.Require().NoError(s.conn.Db().AutoMigrate(&model.Car{}))
}

func (s *carRepoIntegrationSuite) TearDownTest() {
	s.conn.Db().Delete(&model.Car{}, s.carId)
	s.databaseSuite.TearDownTest()
}

func TestCarRepoIntegrationSuite(t *testing.T) {
	suite.Run(t, new(carRepoIntegrationSuite))
}

func (s *carRepoIntegrationSuite) TestCreate() {
	carModel := &domainmodel.Car{
		VinNumber: "4Y1SL65848Z411439",
		Color:     "Spectral Blue",
		AggregateRoot: &domainmodelpkg.AggregateRoot{
			ID:        s.carId,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	car := domain.CarFromModel(carModel)
	s.Require().NoError(s.sut.Create(s.ctx, car))

	var gormCarStored *model.Car
	s.conn.Db().First(&gormCarStored, s.carId)
	s.Require().NotNil(gormCarStored.Car)

	s.Require().Equal(carModel.VinNumber, gormCarStored.VinNumber)
	s.Require().Equal(carModel.Color, gormCarStored.Color)

	requireEqualDates(carModel.CreatedAt, gormCarStored.CreatedAt, s.Require())
	requireEqualDates(carModel.UpdatedAt, gormCarStored.UpdatedAt, s.Require())
}
