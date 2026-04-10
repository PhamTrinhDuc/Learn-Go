package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DB       int //Redis có nhiều database (thường 0-15)
	PoolSize int // tương đương MaxConns trong PG - tối đa X connect 1 lúc
	MinCons  int // tương đương MinConns trong PG - tối thiếu X connect 1 lúc
}

type RedisClient struct {
	Client *redis.Client
}

func (cfg *RedisConfig) Validate() error {
	if cfg.Host == "" || cfg.Username == "" || cfg.Password == "" {
		return fmt.Errorf("redis host, username required")
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
	// 0. Validate
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}
	// 1. Khởi tạo options cho go redis
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	opts := &redis.Options{
		Addr:     addr,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.DB,

		// 2. Cấu hình Connection Pool
		PoolSize:        cfg.PoolSize,
		MinIdleConns:    cfg.MinCons,
		PoolTimeout:     30 * time.Second, // thời gian chờ để lấy 1 kết nối
		ConnMaxIdleTime: 1 * time.Minute,  // thời gian 1 kết nối rảnh rỗi
		ConnMaxLifetime: 1 * time.Minute,  // thời gian 1 kết nối tồn tại
	}
	// 3. Khởi tạo client
	client := redis.NewClient(opts)
	// 4. Test connection (Khác với DB truyền thống, Redis thường Ping ngay để biết client sống hay không)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}
	return &RedisClient{Client: client}, nil
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
}

func (r *RedisClient) FlushDB(ctx context.Context) error {
	return r.Client.FlushDB(ctx).Err()
}
