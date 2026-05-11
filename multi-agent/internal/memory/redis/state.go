package redis

import (
	"iter"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/adk/session"
)

// redisState implements session.State with Redis persistence.
// It holds the merged (all tiers) state in memory and routes writes to the
// correct Redis key based on the key prefix.
type redisState struct {
	data    map[string]any
	client  *redis.Client
	key     string
	ttl     time.Duration
	service *RedisService
	appName string
	userID  string
}

func newRedisState(initial map[string]any, client *redis.Client, key string, ttl time.Duration, service *RedisService, appName, userID string) *redisState {
	data := make(map[string]any)
	for k, v := range initial {
		data[k] = v
	}
	return &redisState{
		data:    data,
		client:  client,
		key:     key,
		ttl:     ttl,
		service: service,
		appName: appName,
		userID:  userID,
	}
}

func (s *redisState) Data() map[string]any {
	return s.data
}

func (s *redisState) Get(key string) (any, error) {
	return nil, nil
}
func (s *redisState) Set(key string, value any) error {
	return nil
}
func (s *redisState) All() iter.Seq2[string, any] {
	return nil
}

var _ session.State = (*redisState)(nil)
