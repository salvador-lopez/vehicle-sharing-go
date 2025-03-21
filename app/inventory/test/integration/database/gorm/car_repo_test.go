//go:build integration

package gorm

import (
	"testing"
	"time"
	pkgdomain "vehicle-sharing-go/pkg/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	gormvehicle "vehicle-sharing-go/app/inventory/internal/vehicle/database/gorm"
	"vehicle-sharing-go/app/inventory/internal/vehicle/database/gorm/model"
	"vehicle-sharing-go/app/inventory/internal/vehicle/domain"
	domainmodel "vehicle-sharing-go/app/inventory/internal/vehicle/domain/model"

	"vehicle-sharing-go/pkg/database/test/integration/gorm"
	domainmodelpkg "vehicle-sharing-go/pkg/domain/model"
)

type carRepoIntegrationSuite struct {
	gorm.DatabaseSuite
	carId uuid.UUID
	sut   *gormvehicle.CarRepository
}

func (s *carRepoIntegrationSuite) SetupSuite() {
	s.DatabaseSuite.SetupSuite()
	s.initDb()
	s.carId = uuid.New()
	s.sut = gormvehicle.NewCarRepository(s.Conn())
}

func (s *carRepoIntegrationSuite) initDb() {
	s.Require().NoError(s.Conn().Db().AutoMigrate(&model.Car{}))
}

func (s *carRepoIntegrationSuite) TearDownTest() {
	s.Conn().Db().Delete(&model.Car{}, s.carId)
	s.DatabaseSuite.TearDownTest()
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
	s.Require().NoError(s.sut.Create(s.Ctx(), car))

	var gormCarStored *model.Car
	s.Conn().Db().First(&gormCarStored, s.carId)
	s.Require().NotNil(gormCarStored.Car)

	s.Require().Equal(carModel.VinNumber, gormCarStored.VinNumber)
	s.Require().Equal(carModel.Color, gormCarStored.Color)

	gorm.RequireEqualDates(carModel.CreatedAt, gormCarStored.CreatedAt, s.Require())
	gorm.RequireEqualDates(carModel.UpdatedAt, gormCarStored.UpdatedAt, s.Require())
}

func (s *carRepoIntegrationSuite) TestCreateSameCarTwiceReturnErrCarAlreadyExist() {
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
	s.Require().NoError(s.sut.Create(s.Ctx(), car))
	err := s.sut.Create(s.Ctx(), car)
	s.Require().ErrorIs(err, pkgdomain.ErrConflict)
	s.Require().EqualError(err, "domain conflict: car already exist")
}
