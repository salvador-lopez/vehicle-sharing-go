//go:build unit

package redis_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-mysql-org/go-mysql/mysql"
	redisclient "github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	binlogredis "vehicle-sharing-go/pkg/message-relay/source/binlog/redis"
	"vehicle-sharing-go/pkg/message-relay/source/binlog/redis/mock"
)

const (
	keyName      = "binlog-position"
	positionJSON = `{"name":"mysql-bin.000007","pos":1280}`
)

func newTestPosition() mysql.Position {
	return mysql.Position{Name: "mysql-bin.000007", Pos: 1280}
}

type stateSuite struct {
	suite.Suite
	ctx        context.Context
	mockCtrl   *gomock.Controller
	mockClient *mock.MockClient
	sut        *binlogredis.StateHandler
}

func TestStateSuite(t *testing.T) {
	suite.Run(t, new(stateSuite))
}

func (s *stateSuite) SetupTest() {
	s.ctx = context.Background()
	s.mockCtrl = gomock.NewController(s.T())
	s.mockClient = mock.NewMockClient(s.mockCtrl)
	s.sut = binlogredis.NewStateHandler(s.mockClient, keyName)
}

func (s *stateSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func (s *stateSuite) TestGetLastPositionReturnsEmptyPositionWhenKeyIsMissing() {
	s.mockClient.EXPECT().
		Get(s.ctx, keyName).
		Return(redisclient.NewStringResult("", redisclient.Nil))

	position, err := s.sut.GetLastPosition(s.ctx)

	s.Require().NoError(err)
	s.Require().Equal(mysql.Position{}, position)
}

func (s *stateSuite) TestGetLastPositionReturnsEmptyPositionForWrappedMissingKeyError() {
	s.mockClient.EXPECT().
		Get(s.ctx, keyName).
		Return(redisclient.NewStringResult("", fmt.Errorf("get key: %w", redisclient.Nil)))

	position, err := s.sut.GetLastPosition(s.ctx)

	s.Require().NoError(err)
	s.Require().Equal(mysql.Position{}, position)
}

func (s *stateSuite) TestGetLastPositionReturnsRedisError() {
	redisErr := errors.New("redis unavailable")
	s.mockClient.EXPECT().
		Get(s.ctx, keyName).
		Return(redisclient.NewStringResult("", redisErr))

	_, err := s.sut.GetLastPosition(s.ctx)

	s.Require().ErrorIs(err, redisErr)
}

func (s *stateSuite) TestGetLastPositionReturnsJSONError() {
	s.mockClient.EXPECT().
		Get(s.ctx, keyName).
		Return(redisclient.NewStringResult("not-json", nil))

	_, err := s.sut.GetLastPosition(s.ctx)

	s.Require().Error(err)
}

func (s *stateSuite) TestGetLastPositionReturnsStoredPosition() {
	s.mockClient.EXPECT().
		Get(s.ctx, keyName).
		Return(redisclient.NewStringResult(positionJSON, nil))

	position, err := s.sut.GetLastPosition(s.ctx)

	s.Require().NoError(err)
	s.Require().Equal(newTestPosition(), position)
}

func (s *stateSuite) TestSetLastPositionStoresJSONWithoutExpiration() {
	position := newTestPosition()
	s.mockClient.EXPECT().
		Set(s.ctx, keyName, []byte(positionJSON), time.Duration(0)).
		Return(redisclient.NewStatusResult("OK", nil))

	s.Require().NoError(s.sut.SetLastPosition(s.ctx, position))
}

func (s *stateSuite) TestSetLastPositionReturnsRedisError() {
	redisErr := errors.New("redis unavailable")
	s.mockClient.EXPECT().
		Set(s.ctx, keyName, gomock.Any(), time.Duration(0)).
		Return(redisclient.NewStatusResult("", redisErr))

	s.Require().ErrorIs(s.sut.SetLastPosition(s.ctx, newTestPosition()), redisErr)
}
