package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"learn-go/workspaces/auth"
	"learn-go/workspaces/protocol"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/redis/go-redis/v9"
)

type FixedWindowLimiter struct {
	redis        *redis.Client
	defaultLimit int
	window       time.Duration
}

type TokenBucketLimiter struct {
	redis    *redis.Client
	capacity int
	rate     float64 // token per senconds
}

var FixedWindowScript = `
local curr = redis.call("INCR", KEYS[1])
if curr == 1 then 
	redis.call("EXPIRE", KEYS[1], ARGV[1])
end
return curr
`

var TokenBucketScript = `
-- KEYS[1]: Key định danh cho Tenant (ví dụ: "ratelimit:tenant_A")
-- ARGV[1]: Capacity (Sức chứa tối đa của xô)
-- ARGV[2]: Refill Rate (Tốc độ nạp thẻ: số thẻ/giây)
-- ARGV[3]: Now (Thời điểm hiện tại - Unix timestamp)

-- BƯỚC 1: lấy dữ liệu cũ từ xô trong Redis 
local bucket = redis.call('hmget', KEYS[1], 'tokens', 'last_updated')
local last_tokens = tonumber(bucket[1]) or tonumber(ARGV[1]) --Nếu xô mới tinh, cho đầy thẻ luôn
local last_updated = tonumber(bucket[2]) or tonumber(ARGV[3]) 

-- BƯỚC 2: "TÍNH NHẨM" - Quá trình Lazy Refill
local duration = math.max(0, ARGV[3] - last_updated) --kể từ lần cuối request tới nay được bao lâu rồi? 
local added_tokens = duration * ARGV[2] --số tokens được nạp từ lần cuối tới nay 

-- BƯỚC 3: CẬP NHẬT XÔ (Không cho vượt quá Capacity)
local tokens_cu = math.min(added_tokens+last_tokens, ARGV[1]) --số lượng nạp > sức chứa thì cũng giới hạn lại

-- BƯỚC 4: QUYẾT ĐỊNH
if tokens_cu >= 1 then 
	tokens_cu = tokens_cu - 1 --nếu còn >=1 token thì trừ đi cho lần request này 
	-- cập nhật trạng thái hiện tại và token còn lại 
	redis.call('hmset', KEYS[1], 'tokens', tokens_cu, 'last_updated', ARGV[3])
	redis.call('expire', KEYS[1], 600) -- set hết hạn sau 10p
	return 1
else 
	return 0 -- không còn đủ token -> limit 
end
`

// Khởi tạo
func NewFixedWindowLimiter(redisClient *redis.Client, defaultLimit int, window time.Duration) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		redis:        redisClient,
		defaultLimit: defaultLimit,
		window:       window,
	}
}

func NewTokenBucketLimiter(redisClient *redis.Client, capacity int, rate float64) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		redis:    redisClient,
		capacity: capacity,
		rate:     rate,
	}
}

func (rl *TokenBucketLimiter) TokenBucketCheckLimit(ctx context.Context, tenantID string) (bool, error) {
	key := fmt.Sprintf("ratelimit:%s", tenantID)
	now := time.Now().Unix()

	result, err := rl.redis.Eval(
		ctx,
		TokenBucketScript,
		[]string{key},
		rl.capacity,
		rl.rate,
		now,
	).Int()

	if err != nil {
		return false, fmt.Errorf("Lua script Token Bucket failed: %w", err)
	}
	return result == 1, nil
}

// Rate limit với cửa số 1 phút giới hạn 100 request (100 request/1 phút)
func (rl *FixedWindowLimiter) WindowCheckLimit(ctx context.Context, tenantID string) (bool, error) {
	key := fmt.Sprintf("ratelimit:%s%d", tenantID, time.Now().Unix()/60) // cửa sổ 1 phút
	// Execute script
	// Keys: key
	// Args: rl.window (s)
	result, err := rl.redis.Eval(
		ctx,
		FixedWindowScript,
		[]string{key},
		int(rl.window.Seconds())).Result()
	if err != nil {
		return false, fmt.Errorf("lua script Fixed Window failed: %w", err)
	}
	count, ok := result.(int64)
	if !ok {
		return false, fmt.Errorf("unexpected return type from redis: %T", result)
	}
	return int(count) <= rl.defaultLimit, nil
}

func (rl *FixedWindowLimiter) SendError(w http.ResponseWriter, id interface{}, code int, message string) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusTooManyRequests)

	response := protocol.NewErrorResponse(id, code, message, nil)
	json.NewEncoder(w).Encode(response)
}

func (rl *TokenBucketLimiter) SendError(w http.ResponseWriter, id interface{}, code int, message string) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusTooManyRequests)

	response := protocol.NewErrorResponse(id, code, message, nil)
	json.NewEncoder(w).Encode(response)
}

func (rl *FixedWindowLimiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Extract tenant ID
		// 2. Check limit (if err => don't block request. if !allowed => block request)
		ctx := r.Context()

		TenantID, err := auth.ExtractTenantID(ctx)
		if err != nil {
			// If no tenant ID, skip rate limiting ((for unauthenticated requests))
			next.ServeHTTP(w, r)
			return
		}

		allowed, err := rl.WindowCheckLimit(ctx, TenantID)
		if err != nil {
			// Log error but don't block request
			fmt.Printf("Rate limit check error: %v\n", err)
			next.ServeHTTP(w, r)
			return
		}

		if !allowed {
			rl.SendError(w, nil, protocol.RateLimitExceeded, "Rate limit exeeded for tenant")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rl *TokenBucketLimiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Extract tenant ID
		// 2. Check limit (if err => don't block request. if !allowed => block request)
		ctx := r.Context()

		TenantID, err := auth.ExtractTenantID(ctx)
		if err != nil {
			// If no tenant ID, skip rate limiting ((for unauthenticated requests))
			next.ServeHTTP(w, r)
			return
		}

		allowed, err := rl.TokenBucketCheckLimit(ctx, TenantID)
		if err != nil {
			// Log error but don't block request
			fmt.Printf("Rate limit check error: %v\n", err)
			next.ServeHTTP(w, r)
			return
		}

		if !allowed {
			rl.SendError(w, nil, protocol.RateLimitExceeded, "Rate limit exeeded for tenant")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func FixedWindowMockRequest(limiter *FixedWindowLimiter, tenantID string) int {
	ctx := context.WithValue(context.Background(), auth.ContextKeyTenantID, tenantID)

	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// init request and store response
	req := httptest.NewRequest("POST", "/mcp", nil).WithContext(ctx)
	rr := httptest.NewRecorder()

	handler := limiter.Handler(testHandler)
	handler.ServeHTTP(rr, req)

	return rr.Code
}

func TokenBucketMockRequest(limiter *TokenBucketLimiter, tenantID string) int {
	ctx := context.WithValue(context.Background(), auth.ContextKeyTenantID, tenantID)

	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// init request and store response
	req := httptest.NewRequest("POST", "/mcp", nil).WithContext(ctx)
	rr := httptest.NewRecorder()

	handler := limiter.Handler(testHandler)
	handler.ServeHTTP(rr, req)

	return rr.Code
}
