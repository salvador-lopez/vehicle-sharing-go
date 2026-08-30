package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-redis/redis/v8"
)

//go:generate mockgen -destination=mock/redis_client_mock.go -package=mock . RedisClient
type RedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

// StateHandler implements binlog.StateHandler by persisting the MySQL binlog
// position as JSON in a single Redis key.
type StateHandler struct {
	client  *redis.Client
	keyName string
}

// NewStateHandler constructs a StateHandler that reads and writes the binlog
// position under keyName.
func NewStateHandler(client *redis.Client, keyName string) *StateHandler {
	return &StateHandler{client: client, keyName: keyName}
}

// GetLastPosition returns the last persisted binlog position.
// Returns an empty Position (beginning of binlog) if the key does not exist yet.
func (s *StateHandler) GetLastPosition() (mysql.Position, error) {
	val, err := s.client.Get(context.Background(), s.keyName).Result()
	if err == redis.Nil {
		return mysql.Position{}, nil
	}
	if err != nil {
		return mysql.Position{}, err
	}

	var p binlogPosition
	if err := json.Unmarshal([]byte(val), &p); err != nil {
		return mysql.Position{}, err
	}

	return mysql.Position{Name: p.Name, Pos: p.Pos}, nil
}

// SetLastPosition persists the given binlog position with no TTL.
func (s *StateHandler) SetLastPosition(p mysql.Position) error {
	data, err := json.Marshal(binlogPosition{Name: p.Name, Pos: p.Pos})
	if err != nil {
		return err
	}
	return s.client.Set(context.Background(), s.keyName, data, 0).Err()
}

type binlogPosition struct {
	Name string `json:"name,omitempty"`
	Pos  uint32 `json:"pos,omitempty"`
}
