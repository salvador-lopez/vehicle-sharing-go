package gorm

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/stretchr/testify/suite"

	gormpkg "vehicle-sharing-go/pkg/infrastructure/database/gorm"
)

type databaseSuite struct {
	suite.Suite
	ctx       context.Context
	cancelFun context.CancelFunc
	conn      *gormpkg.Connection
}

func (s *databaseSuite) SetupSuite() {
	s.createDb()
}

func (s *databaseSuite) createDb() {
	port, err := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	s.Require().NoError(err)

	conn, err := gormpkg.NewConnectionFromConfig(&gormpkg.Config{
		UserName:     os.Getenv("MYSQL_USER"),
		Password:     os.Getenv("MYSQL_PASSWORD"),
		DatabaseName: os.Getenv("MYSQL_DATABASE"),
		Host:         os.Getenv("MYSQL_HOST"),
		Port:         port,
	})
	s.Require().NoError(err)

	s.conn = conn
}

func (s *databaseSuite) SetupTest() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	s.ctx = ctx
	s.cancelFun = cancel
}

func (s *databaseSuite) TearDownTest() {
	s.cancelFun()
}
