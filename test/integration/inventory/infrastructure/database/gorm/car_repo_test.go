//go:build integration_inventory

package gorm

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"vehicle-sharing-go/internal/inventory/vehicle/domain"
	gormvehicle "vehicle-sharing-go/internal/inventory/vehicle/infrastructure/database/gorm"
	"vehicle-sharing-go/internal/inventory/vehicle/infrastructure/database/gorm/model"
)

type carRepoIntegrationSuite struct {
	suite.Suite
	ctx       context.Context
	cancelFun context.CancelFunc
	db        *gorm.DB
	carId     uuid.UUID
	sut       *gormvehicle.CarRepository
}

func (s *carRepoIntegrationSuite) SetupSuite() {
	s.initDb()
	s.carId = uuid.New()
	s.sut = gormvehicle.NewCarRepository(s.db)
}

func (s *carRepoIntegrationSuite) initDb() {
	s.db = createDb(s.Require())
	s.Require().NoError(s.db.AutoMigrate(&model.Car{}))
}

func (s *carRepoIntegrationSuite) SetupTest() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	s.ctx = ctx
	s.cancelFun = cancel
}

func (s *carRepoIntegrationSuite) TearDownTest() {
	s.db.Delete(&model.Car{}, s.carId)
	s.cancelFun()
}

func (s *carRepoIntegrationSuite) TearDownSuite() {
	sqlDb, err := s.db.DB()
	s.Require().NoError(err)

	s.Require().NoError(sqlDb.Close())
}

func TestCarRepoIntegrationSuite(t *testing.T) {
	suite.Run(t, new(carRepoIntegrationSuite))
}

func (s *carRepoIntegrationSuite) TestCreateCar() {
	carDTO := &domain.CarDTO{
		VIN:   "4Y1SL65848Z411439",
		Color: "Spectral Blue",
		BaseDTO: &domain.BaseDTO{
			ID:        s.carId,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	car := domain.HydrateCar(carDTO)
	s.Require().NoError(s.sut.Create(s.ctx, car))

	var gormCarStored *model.Car
	s.db.First(&gormCarStored, s.carId)
	s.Require().NotNil(gormCarStored.CarDTO)

	s.Require().Equal(carDTO.VIN, gormCarStored.VIN)
	s.Require().Equal(carDTO.Color, gormCarStored.Color)

	requireEqualDates(carDTO.CreatedAt, gormCarStored.CreatedAt, s.Require())
	requireEqualDates(carDTO.UpdatedAt, gormCarStored.UpdatedAt, s.Require())
}
