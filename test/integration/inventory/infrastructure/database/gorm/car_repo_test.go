//go:build integration

package gorm_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"vehicle-sharing-go/internal/inventory/vehicle/domain"
	gormvehicle "vehicle-sharing-go/internal/inventory/vehicle/infrastructure/database/gorm"
)

type carRepoIntegrationSuite struct {
	suite.Suite
	db  *gorm.DB
	sut *gormvehicle.CarRepository
}

func (s *carRepoIntegrationSuite) SetupSuite() {
	dsn := fmt.Sprintf(`%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC`,
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	s.Require().NoError(err)

	s.db = db
	s.sut = gormvehicle.NewCarRepository(s.db)
}

func (s *carRepoIntegrationSuite) TearDownTest() {
	s.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&domain.CarDTO{})
}

func (s *carRepoIntegrationSuite) TearDownSuite() {
	sqlDb, err := s.db.DB()
	s.Require().NoError(err)

	err = sqlDb.Close()
	s.Require().NoError(err)
}

func TestCarRepoIntegrationSuite(t *testing.T) {
	suite.Run(t, new(carRepoIntegrationSuite))
}
