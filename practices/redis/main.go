package main

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host        string
	Port        int
	Username    string
	Password    string
	DB          int //Redis có nhiều database (thường 0-15)
	PoolSize    int // tương đương MaxConns trong PG - tối đa X connect 1 lúc
	MinCons     int // tương đương MinConns trong PG - tối thiếu X connect 1 lúc
	FixedWindow interface{}
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

func GetRedisConfig() RedisConfig {
	return RedisConfig{
		Host:     "localhost",
		Port:     6379,
		Username: "jiyuu",
		Password: "a2amcpgo",
	}
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

type Session struct {
	ID   string            `json:"id"`
	User string            `json:"user"`
	Data map[string]string `json:"data"`
}

func GetSet(ctx context.Context, client *redis.Client) {
	// --- LUỒNG 1: TẠO SESSION (Dùng kiểu STRING trong Redis) ---
	sessionID := "sess-001"
	sessKey := "session:" + sessionID
	mySess := Session{
		ID:   sessionID,
		User: "Jiyuu",
		Data: map[string]string{"username": "Jiyuu", "mail": "Jiyuu@gmail.com"},
	}
	// Chuyển Object thành JSON (Marshalling)
	sessData, _ := json.Marshal(mySess)

	// Lưu vào Redis với thời gian sống (TTL) là 10 phút
	client.Set(ctx, sessKey, sessData, 1*time.Minute)
	// Lấy ra xem dữ liệu đã được ghi vào redis chưa
	insertedData := client.Get(ctx, sessKey)
	fmt.Println(insertedData.String())
}

type Message struct {
	Role    string `json:"role"`
	Message string `json:"message"`
}

func List(ctx context.Context, client *redis.Client) {
	sessionID := "sess-001"
	eventKey := "event:" + sessionID

	data1, _ := json.Marshal(Message{
		Role:    "user",
		Message: "Xin chào",
	})

	data2, _ := json.Marshal(Message{
		Role:    "assistant",
		Message: "Tôi có thể giúp gì cho bạn",
	})

	client.Del(ctx, eventKey) // xóa eventkey đang lưu dữ liệu
	client.RPush(ctx, eventKey, data1)
	client.Expire(ctx, eventKey, 5*time.Second) // xét TTL
	client.RPush(ctx, eventKey, data2)
	client.Expire(ctx, eventKey, 5*time.Second)

	dataInserted, err := client.LRange(ctx, eventKey, 0, -1).Result()
	if err != nil {
		fmt.Printf("failed to get data from event key: %s", eventKey)
	}

	var messages []Message
	for _, item := range dataInserted {
		var msg Message
		err := json.Unmarshal([]byte(item), &msg)
		if err != nil {
			fmt.Printf("failed to parse json data for event key: %s", eventKey)
		}
		messages = append(messages, msg)
	}
	fmt.Println(reflect.TypeOf(messages))
	fmt.Println(messages)
}

type InfoUser struct {
	Name  string  `json:"name"`
	Level int     `json"level"`
	Power float64 `json:"power"`
}

func Hash(ctx context.Context, client *redis.Client) {
	userState := "userstate:app01:user001"
	// Lưu ý: HSet không hiểu struct và json, chỉ: map hoặc cặp field-value
	client.HSet(ctx, userState, map[string]interface{}{
		"name":  "Nam",
		"level": 12,
		"power": 9999,
	})
	insertedData, _ := client.HGetAll(ctx, userState).Result()
	fmt.Println("Dữ liệu được insert vào Redis:")
	fmt.Println(insertedData)

	client.HSet(ctx, userState,
		"name", "Huy",
	)
	insertedData, _ = client.HGetAll(ctx, userState).Result()
	fmt.Println("Dữ liệu sau khi được cập nhật:")
	fmt.Println(insertedData) // => dữ liệu được cập nhật

	client.Del(ctx, userState)
}

func Set(ctx context.Context, client *redis.Client) {
	indexKey := "sessions:app01:user001"

	client.SAdd(ctx, indexKey, "session-001")
	client.SAdd(ctx, indexKey, "session-002")
	client.SAdd(ctx, indexKey, "session-001")

	allSessions, _ := client.SMembers(ctx, indexKey).Result()
	fmt.Println("Danh sách Session ID của User này (Set):")
	for _, id := range allSessions {
		fmt.Printf("   - ID: %s\n", id)
	}

	client.Del(ctx, indexKey)
}

func main() {
	ctx := context.Background()

	rdb, err := NewRedis(ctx, GetRedisConfig())
	if err != nil {
		fmt.Println("failed to connect to redis")
	}
	client := rdb.Client
	// GetSet(ctx,client)
	// List(ctx, client)
	Hash(ctx, client)
	// Set(ctx, client)
}
