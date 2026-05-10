package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/adk/session"
)

// RedisService: implement session.Service adk using Redis backend
type RedisService struct {
	client       *redis.Client
	ttl          time.Duration
	appStateTTL  time.Duration
	userStateTTL time.Duration
}

type redisSession struct {
	id             string // ID Session
	appName        string
	userID         string
	state          *redisState
	events         *redisEvents
	lastUpdateTime time.Time
}

func (s *redisSession) ID() string                { return s.id }
func (s *redisSession) AppName() string           { return s.appName }
func (s *redisSession) UserID() string            { return s.userID }
func (s *redisSession) State() session.State      { return s.state }
func (s *redisSession) Events() session.Events    { return s.events }
func (s *redisSession) LastUpdateTime() time.Time { return s.lastUpdateTime }

type redisState struct {
	data    map[string]any
	client  *redis.Client
	key     string
	ttl     time.Duration
	service *RedisService
	appName string
	userID  string
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

type redisEvents struct {
}

func (s *redisEvents) All() iter.Seq[*session.Event] {
	return nil
}

func (s *redisEvents) Len() int {
	return 0
}

func (s *redisEvents) At(i int) *session.Event {
	return nil
}

var _ session.Session = (*redisSession)(nil)
var _ session.State = (*redisState)(nil)
var _ session.Events = (*redisEvents)(nil)

// storableSession is the JSON-serializable representation of a session.
// State only contains session-scoped keys (no app: or user: prefixed keys).
type storableSession struct {
	ID             string         `json:"id"`
	AppName        string         `json:"app_name"`
	UserID         string         `json:"user_id"`
	State          map[string]any `json:"state"`
	LastUpdateTime time.Time      `json:"last_update_time"`
}

// RedisSessionServiceConfig holds configuration for RedisSessionService.
type RedisConfig struct {
	Host     string
	Port     int
	Username string
	// Password for Redis authentication (optional)
	Password string
	// DB is the Redis database number
	DB int
	// TTL is the session expiration time (default: 24 hours)
	TTL time.Duration
	// Timeout is the time connect to Redis
	Timeout time.Duration
	// AppStateTTL is the expiration time for app-scoped state.
	// Defaults to 0 (no expiration), matching the canonical ADK behaviour
	// where app state outlives individual sessions.
	AppStateTTL time.Duration
	// UserStateTTL is the expiration time for user-scoped state.
	// Defaults to 0 (no expiration), matching the canonical ADK behaviour
	// where user state outlives individual sessions.
	UserStateTTL time.Duration
}

func NewRedisService(cfg *RedisConfig) (*RedisService, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if cfg.Timeout == 0 {
		cfg.Timeout = 30 // 30s timeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	if cfg.TTL == 0 {
		cfg.TTL = 24 * time.Hour
	}

	return &RedisService{
		client:       client,
		ttl:          cfg.TTL,
		appStateTTL:  cfg.AppStateTTL,
		userStateTTL: cfg.UserStateTTL,
	}, nil
}

// Close close connect to Redis
func (s *RedisService) Close() error {
	return s.client.Close()
}

func (s *RedisService) FlushDB(ctx context.Context) error {
	return s.client.FlushDB(ctx).Err()
}

// appStateKey return specific app key
func (s *RedisService) appStateKey(appName string) string {
	return fmt.Sprintf("appstate:%s", appName)
}

// userStateKey return specific user key
func (s *RedisService) userStateKey(appName string, userID string) string {
	return fmt.Sprintf("users:%s:%s", appName, userID)
}

// sesionStateKey return key specific session key
func (s *RedisService) sessionKey(appName string, userID string, sessionKey string) string {
	return fmt.Sprintf("session:%s:%s:%s", appName, userID, sessionKey)
}

// sessionKeys return list session keys
func (s *RedisService) sessionKeys(appName string, userID string) string {
	return fmt.Sprintf("sessions:%s:%s", appName, userID)
}

// Create creates a new session. It returns an error if a session with the
// same ID already exists, matching the canonical ADK behaviour.
func (s *RedisService) Create(ctx context.Context, req *session.CreateRequest) (*session.CreateResponse, error) {
	// Hàm tạo 1 session mới, nếu giống CURD thông thường thì chỉ return redisSession là xong.
	// Nhưng tại sao lại phải có những chỗ updateAppState và updateAppUser?
	// Xét TH sau với hàm updateAppUser: khách hàng ấn vào 1 sản phẩm giày nike trên web và chọn chat ngay:
	// Lúc này việc update có nhiệm vụ cập nhật các state như: {"user:last_viewed_product": "Giày Nike", "current_page": "Promotion"}
	// Các thông tin này sẽ có ích cho AI khi chat

	// Coding start:
	// 1. check exists session key
	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	key := s.sessionKey(req.AppName, req.UserID, sessionID) // key of specific session
	if exists, _ := s.client.Exists(ctx, key).Result(); exists > 0 {
		return nil, fmt.Errorf("session %s already exists", sessionID)
	}
	// 2. Update appState and UserState
	appDelta, userDelta, sessionDelta := extractStateDeltas(req.State)
	appState, err := s.updateAppState(ctx, req.AppName, appDelta)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	userState, err := s.updateUserState(ctx, req.AppName, req.UserID, userDelta)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	mergedState := s.mergeState(appState, userState, sessionDelta)

	redisSession := &redisSession{
		id:             sessionID,
		appName:        req.AppName,
		userID:         req.UserID,
		state:          &redisState{},
		events:         &redisEvents{},
		lastUpdateTime: time.Now(),
	}

	storableSess := storableSession{
		ID:             sessionID,
		AppName:        req.AppName,
		UserID:         req.UserID,
		State:          mergedState,
		LastUpdateTime: redisSession.lastUpdateTime,
	}

	data, err := json.Marshal(storableSess)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal session data: %w", err)
	}
	// store session
	if err := s.client.Set(ctx, key, data, s.ttl).Err(); err != nil {
		return nil, fmt.Errorf("failed to store session to redis: %w", err)
	}
	// store key session to Set
	idxSessKey := s.sessionKeys(req.AppName, req.UserID)
	if err := s.client.SAdd(ctx, idxSessKey).Err(); err != nil {
		return nil, fmt.Errorf("failed to update session index: %w", err)
	}
	// set expire
	s.client.Expire(ctx, idxSessKey, s.ttl)
	return &session.CreateResponse{Session: redisSession}, nil
}

func (s *RedisService) mergeState(appState map[string]any, userState map[string]any, sessionState map[string]any) map[string]any {
	data := make(map[string]any)
	for k, v := range sessionState {
		data[k] = v
	}
	for k, v := range userState {
		data[session.KeyPrefixUser+k] = v
	}
	for k, v := range appState {
		data[session.KeyPrefixApp+k] = v
	}
	return data
}

// updateAppState update data State of User
func (s *RedisService) updateUserState(
	ctx context.Context,
	appName string,
	userID string,
	userDelta map[string]any) (map[string]any, error) {

	if userDelta == nil {
		return map[string]any{}, nil
	}
	key := s.userStateKey(appName, userID)
	fields := marshalHashFields(userDelta)
	if _, err := s.client.HSet(ctx, key, fields).Result(); err != nil {
		return nil, fmt.Errorf("failed to update user State: %w", err)
	}
	if s.userStateTTL > 0 {
		s.client.Expire(ctx, key, s.userStateTTL)
	} else {
		s.client.Persist(ctx, key)
	}
	return s.loadUserState(ctx, appName, userID), nil
}

// updateAppState update data State of App
func (s *RedisService) updateAppState(ctx context.Context, appName string, appDelta map[string]any) (map[string]any, error) {
	// if not available data => return old data store in key
	if appDelta == nil {
		return s.loadAppState(ctx, appName), nil
	}

	key := s.appStateKey(appName)
	fields := marshalHashFields(appDelta)
	if _, err := s.client.HSet(ctx, key, fields).Result(); err != nil {
		return nil, fmt.Errorf("failed to update app state: %w", err)
	}
	if s.appStateTTL > 0 {
		s.client.Expire(ctx, key, s.appStateTTL)
	} else {
		s.client.Persist(ctx, key)
	}
	return s.loadAppState(ctx, appName), nil
}

// loadAppState return Hash data of App inherit loadHashState
func (s *RedisService) loadAppState(ctx context.Context, appName string) map[string]any {
	return s.loadHashState(ctx, appName)
}

// loadAppState return Hash data of User inherit loadHashState
func (s *RedisService) loadUserState(ctx context.Context, appName string, userID string) map[string]any {
	return s.loadHashState(ctx, s.userStateKey(appName, userID))
}

// loadHashState load Hash data from Redis
func (s *RedisService) loadHashState(ctx context.Context, key string) map[string]any {
	result, err := s.client.HGetAll(ctx, key).Result()
	if err != nil || len(result) == 0 {
		log.Print("failed to get hash data: %w", err)
	}
	return unmarshalHashFields(result)
}
