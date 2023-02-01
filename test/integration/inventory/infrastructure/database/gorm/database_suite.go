package gorm

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type databaseSuite struct {
	suite.Suite
	ctx       context.Context
	cancelFun context.CancelFunc
	db        *gorm.DB
}

func (s *databaseSuite) SetupSuite() {
	s.createDb()
}

func (s *databaseSuite) createDb() {
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
}

func (s *databaseSuite) SetupTest() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	s.ctx = ctx
	s.cancelFun = cancel
}

func (s *databaseSuite) TearDownTest() {
	s.cancelFun()
}

func (s *databaseSuite) TearDownSuite() {
	sqlDb, err := s.db.DB()
	s.Require().NoError(err)

	s.Require().NoError(sqlDb.Close())
}
