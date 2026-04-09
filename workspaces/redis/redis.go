package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int //Redis có nhiều database (thường 0-15)
	PoolSize int // tương đương MaxConns trong PG - tối đa X connect 1 lúc
	MinCons  int // tương đương MinConns trong PG - tối thiếu X connect 1 lúc
}

type RedisClient struct {
	Client *redis.Client
}

func (cfg *RedisConfig) Validate() error {
	if cfg.Host == "" {
		return fmt.Errorf("redis host required")
	}
	if cfg.Port <= 0 {
		return fmt.Errorf("port cannot be less than or equal to 0")
	}
	if cfg.PoolSize <= 0 {
		cfg.PoolSize = 10
	}
	if cfg.MinCons < 0 {
		cfg.MinCons = 2
	}
	return nil
}

func NewRedis(ctx context.Context, cfg RedisConfig) (*RedisClient, error) {

}
