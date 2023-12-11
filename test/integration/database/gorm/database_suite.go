package gorm

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/stretchr/testify/suite"

	gormpkg "vehicle-sharing-go/pkg/database/gorm"
)

type DatabaseSuite struct {
	suite.Suite
	ctx       context.Context
	cancelFun context.CancelFunc
	conn      *gormpkg.Connection
}

func (s *DatabaseSuite) Ctx() context.Context {
	return s.ctx
}

func (s *DatabaseSuite) Conn() *gormpkg.Connection {
	return s.conn
}

func (s *DatabaseSuite) SetupSuite() {
	s.createDb()
}

func (s *DatabaseSuite) createDb() {
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

func (s *DatabaseSuite) SetupTest() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	s.ctx = ctx
	s.cancelFun = cancel
}

func (s *DatabaseSuite) TearDownTest() {
	s.cancelFun()
}
