//go:build integration_inventory

package gorm

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"vehicle-sharing-go/internal/inventory/vehicle/domain"
	gormvehicle "vehicle-sharing-go/internal/inventory/vehicle/infrastructure/database/gorm"
	"vehicle-sharing-go/internal/inventory/vehicle/infrastructure/database/gorm/model"
	domainpkg "vehicle-sharing-go/pkg/domain"
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
	s.sut = gormvehicle.NewCarRepository(s.db)
}

func (s *carRepoIntegrationSuite) initDb() {
	s.Require().NoError(s.db.AutoMigrate(&model.Car{}))
}

func (s *carRepoIntegrationSuite) TearDownTest() {
	s.db.Delete(&model.Car{}, s.carId)
	s.databaseSuite.TearDownTest()
}

func TestCarRepoIntegrationSuite(t *testing.T) {
	suite.Run(t, new(carRepoIntegrationSuite))
}

func (s *carRepoIntegrationSuite) TestCreate() {
	carDTO := &domain.CarDTO{
		VinNumber: "4Y1SL65848Z411439",
		Color:     "Spectral Blue",
		AgRootDTO: &domainpkg.AgRootDTO{
			ID:        s.carId,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	car := carDTO.ToAggRoot()
	s.Require().NoError(s.sut.Create(s.ctx, car))

	var gormCarStored *model.Car
	s.db.First(&gormCarStored, s.carId)
	s.Require().NotNil(gormCarStored.CarDTO)

	s.Require().Equal(carDTO.VinNumber, gormCarStored.VinNumber)
	s.Require().Equal(carDTO.Color, gormCarStored.Color)

	requireEqualDates(carDTO.CreatedAt, gormCarStored.CreatedAt, s.Require())
	requireEqualDates(carDTO.UpdatedAt, gormCarStored.UpdatedAt, s.Require())
}
