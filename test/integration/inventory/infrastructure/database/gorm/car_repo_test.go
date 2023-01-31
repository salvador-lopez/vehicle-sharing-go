//go:build integration_inventory

package gorm_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
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
	dsn := fmt.Sprintf(`%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC`,
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	s.Require().NoError(err)

	s.Require().NoError(db.AutoMigrate(&model.Car{}))

	s.db = db
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

	tFormat := time.RFC3339
	s.Require().Equal(carDTO.CreatedAt.UTC().Format(tFormat), gormCarStored.CreatedAt.Format(tFormat))
	s.Require().Equal(carDTO.UpdatedAt.UTC().Format(tFormat), gormCarStored.UpdatedAt.Format(tFormat))
}
