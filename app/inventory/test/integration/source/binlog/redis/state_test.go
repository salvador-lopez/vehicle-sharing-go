//go:build integration

package redis

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-mysql-org/go-mysql/mysql"
	redisclient "github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/suite"

	binlogredis "vehicle-sharing-go/app/inventory/cmd/vehicle/message-relay/source/binlog/redis"
)

type stateHandlerIntegrationSuite struct {
	suite.Suite
	client  *redisclient.Client
	ctx     context.Context
	keyName string
	sut     *binlogredis.StateHandler
}

func TestStateHandlerIntegrationSuite(t *testing.T) {
	suite.Run(t, new(stateHandlerIntegrationSuite))
}

func (s *stateHandlerIntegrationSuite) SetupSuite() {
	s.client = redisclient.NewClient(&redisclient.Options{
		Addr: fmt.Sprintf("%s:%s", envOrDefault("REDIS_HOST", "localhost"), envOrDefault("REDIS_PORT", "6379")),
		DB:   0,
	})
}

func (s *stateHandlerIntegrationSuite) SetupTest() {
	s.ctx = context.Background()
	s.keyName = "inventory:integration:test:binlog-position"
	s.Require().NoError(s.client.Del(s.ctx, s.keyName).Err())
	s.sut = binlogredis.NewStateHandler(s.client, s.keyName)
}

func (s *stateHandlerIntegrationSuite) TearDownTest() {
	s.Require().NoError(s.client.Del(s.ctx, s.keyName).Err())
}

func (s *stateHandlerIntegrationSuite) TearDownSuite() {
	s.Require().NoError(s.client.Close())
}

func (s *stateHandlerIntegrationSuite) TestGetLastPositionReturnsZeroValueWhenStateIsMissing() {
	pos, err := s.sut.GetLastPosition()

	s.Require().NoError(err)
	s.Require().Equal(mysql.Position{}, pos)
}

func (s *stateHandlerIntegrationSuite) TestSetLastPositionAndGetLastPositionRoundtrip() {
	expected := mysql.Position{
		Name: "mysql-bin.000007",
		Pos:  1280,
	}

	s.Require().NoError(s.sut.SetLastPosition(expected))

	actual, err := s.sut.GetLastPosition()
	s.Require().NoError(err)
	s.Require().Equal(expected, actual)
}

func (s *stateHandlerIntegrationSuite) TestGetLastPositionReturnsErrorWhenStoredPayloadIsInvalidJSON() {
	s.Require().NoError(s.client.Set(s.ctx, s.keyName, "not-json", time.Minute).Err())

	_, err := s.sut.GetLastPosition()

	s.Require().Error(err)
}

func (s *stateHandlerIntegrationSuite) TestGetLastPositionReturnsErrorWhenRedisClientIsUnavailable() {
	sut := binlogredis.NewStateHandler(s.closedClient(), s.keyName)

	_, err := sut.GetLastPosition()

	s.Require().Error(err)
}

func (s *stateHandlerIntegrationSuite) TestSetLastPositionReturnsErrorWhenRedisClientIsUnavailable() {
	sut := binlogredis.NewStateHandler(s.closedClient(), s.keyName)

	err := sut.SetLastPosition(mysql.Position{Name: "mysql-bin.000007", Pos: 1280})

	s.Require().Error(err)
}

// closedClient returns a Redis client that is already closed, so any command
// deterministically fails with a connection error without relying on
// network flakiness.
func (s *stateHandlerIntegrationSuite) closedClient() *redisclient.Client {
	client := redisclient.NewClient(s.client.Options())
	s.Require().NoError(client.Close())
	return client
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
