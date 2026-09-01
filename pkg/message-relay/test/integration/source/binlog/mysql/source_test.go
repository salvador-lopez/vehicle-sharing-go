//go:build integration

package mysql

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/stretchr/testify/suite"

	binlogmysql "vehicle-sharing-go/pkg/message-relay/source/binlog/mysql"
)

const tableName = "message_relay_binlog_source_test"

type sourceIntegrationSuite struct {
	suite.Suite
	mysql *mysqlHelper

	canal        *canal.Canal
	handler      *capturingEventHandler
	stateHandler *inMemoryStateHandler
	sut          *binlogmysql.Source

	runErr  chan error
	runDone chan struct{}
}

func TestSourceIntegrationSuite(t *testing.T) {
	suite.Run(t, new(sourceIntegrationSuite))
}

func (s *sourceIntegrationSuite) SetupSuite() {
	helper, err := newMySQLHelper()
	s.Require().NoError(err)

	s.mysql = helper
}

func (s *sourceIntegrationSuite) TearDownSuite() {
	s.Require().NoError(s.mysql.close())
}

func (s *sourceIntegrationSuite) SetupTest() {
	s.Require().NoError(s.mysql.recreateTable(tableName))
	s.startSource()
}

func (s *sourceIntegrationSuite) startSource() {
	cfg := canal.NewDefaultConfig()
	cfg.Addr = fmt.Sprintf("%s:%s", envOrDefault("MYSQL_HOST", "localhost"), envOrDefault("MYSQL_PORT", "3308"))
	cfg.User = envOrDefault("MYSQL_USER", "inventory")
	cfg.Password = envOrDefault("MYSQL_PASSWORD", "inventory")
	cfg.Dump.ExecutionPath = ""
	cfg.IncludeTableRegex = []string{fmt.Sprintf(".*\\.%s", tableName)}

	c, err := canal.NewCanal(cfg)
	s.Require().NoError(err)
	s.canal = c

	masterPos, err := c.GetMasterPos()
	s.Require().NoError(err)

	s.stateHandler = newInMemoryStateHandler(masterPos)
	handler, posCh := newCapturingEventHandler()
	s.handler = handler

	s.sut = binlogmysql.NewSource(c, handler, posCh, s.stateHandler, 50*time.Millisecond)

	s.runErr = make(chan error, 1)
	s.runDone = make(chan struct{})
	go func() {
		s.runErr <- s.sut.Run()
		close(s.runDone)
	}()
}

func (s *sourceIntegrationSuite) TearDownTest() {
	s.sut.Close()
	<-s.runDone

	s.Require().NoError(s.mysql.dropTable(tableName))
}

func (s *sourceIntegrationSuite) TestRunStreamsInsertedRowFromRealMySQLBinlog() {
	s.Require().NoError(s.mysql.insertRow(tableName, 1, "hello-binlog"))

	select {
	case row := <-s.handler.rows:
		s.Equal("insert", row.Action)
		s.Require().Len(row.Rows, 1)
		s.Equal(int32(1), row.Rows[0][0])
		s.Equal("hello-binlog", row.Rows[0][1])
	case <-time.After(10 * time.Second):
		s.Fail("timed out waiting for the inserted row to be streamed from the binlog")
	}

	s.Eventually(func() bool {
		return s.stateHandler.lastPosition().Compare(mysql.Position{}) > 0
	}, 5*time.Second, 50*time.Millisecond, "expected the binlog position to be checkpointed")
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
