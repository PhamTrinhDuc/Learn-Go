package middleware

import (
	"context"
	auth "learn-go/workspaces/auth"
	my_redis "learn-go/workspaces/redis"
	utils "learn-go/workspaces/utils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func SetUpRedis(t *testing.T) *my_redis.RedisClient {
	client, err := my_redis.NewRedis(context.Background(), my_redis.RedisConfig{
		Host:     utils.GetEnvOrDefault("REDIS_HOST", "localhost"),
		Port:     6379,
		Username: utils.GetEnvOrDefault("REDIS_USERNAME", "jiyuu"),
		Password: utils.GetEnvOrDefault("REDIS_PASSWORD", "a2amcpgo"),
	})
	require.NoError(t, err)
	assert.NotNil(t, client)
	return client
}

func TestNewFixedWindowLimiter(t *testing.T) {
	client := SetUpRedis(t)
	limiter := NewFixedWindowLimiter(client.Client, 100, time.Minute)

	assert.NotNil(t, limiter)
	assert.Equal(t, 100, limiter.defaultLimit)
	assert.Equal(t, time.Minute, limiter.window)
}

func TestFixedWindowLimiter_Valid(t *testing.T) {
	client := SetUpRedis(t)
	limiter := NewFixedWindowLimiter(client.Client, 100, time.Minute)
	ctx := context.WithValue(context.Background(), auth.ContextKeyTenantID, "a2a-mcp-123")

	// Create test handler
	handleCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleCalled = true
		tenantID, err := auth.ExtractTenantID(r.Context())
		assert.NoError(t, err)
		assert.Equal(t, tenantID, "a2a-mcp-123")

		w.WriteHeader(http.StatusOK)
	})

	// init request and store response
	req := httptest.NewRequest("POST", "/mcp", nil).WithContext(ctx)
	rr := httptest.NewRecorder()

	handler := limiter.Handler(testHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, handleCalled, true)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestFixedWindowLimiter_NoTenantID(t *testing.T) {
	client := SetUpRedis(t)
	limiter := NewFixedWindowLimiter(client.Client, 100, time.Minute)

	// Create test handler
	handleCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleCalled = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("POST", "/mcp", nil)
	rr := httptest.NewRecorder()

	handler := limiter.Handler(testHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, handleCalled, true)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestFixedWindow_CheckLimit_OneTenantID(t *testing.T) {
	tests := []struct {
		name          string
		requestCount  int64
		limit         int
		expectAllowed bool
	}{
		{
			name:          "first request",
			requestCount:  1,
			limit:         100,
			expectAllowed: true,
		},
		{
			name:          "within limit",
			requestCount:  50,
			limit:         100,
			expectAllowed: true,
		},
		{
			name:          "at limit",
			requestCount:  100,
			limit:         100,
			expectAllowed: true,
		},
		{
			name:          "exceeded limit",
			requestCount:  101,
			limit:         100,
			expectAllowed: false,
		},
		{
			name:          "far exceeded",
			requestCount:  500,
			limit:         100,
			expectAllowed: false,
		},
	}

	tenantID := "a2a-mcp-go"
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := SetUpRedis(t)
			// Reset database for a fresh test environment
			client.FlushDB(context.Background())
			defer client.Close()

			limiter := NewFixedWindowLimiter(client.Client, tc.limit, time.Minute)

			var lastCode int
			for i := 0; i < int(tc.requestCount); i++ {
				code := FixedWindowMockRequest(limiter, tenantID) // send request and check limit int method Handler
				lastCode = code
			}
			if tc.expectAllowed {
				assert.Equal(t, http.StatusOK, lastCode)
			} else {
				assert.Equal(t, http.StatusTooManyRequests, lastCode)
			}
		})
	}
}

func TestFixedWindow_CheckLimit_ManyTenantID(t *testing.T) {
	tests := []struct {
		name          string
		tenantID      string
		requestCount  int64
		limit         int
		expectAllowed bool
	}{
		{
			name:          "user 1 request first",
			requestCount:  50,
			tenantID:      "mcp-a2a-1",
			limit:         100,
			expectAllowed: true,
		},
		{
			name:          "user 2 request second",
			requestCount:  50,
			tenantID:      "mcp-a2a-1",
			limit:         100,
			expectAllowed: true,
		},
		{
			name:          "user 2 request first",
			requestCount:  100,
			tenantID:      "mcp-a2a-2",
			limit:         100,
			expectAllowed: true,
		},
		{
			name:          "user 2 request second",
			requestCount:  2,
			tenantID:      "mcp-a2a-2",
			limit:         100,
			expectAllowed: false,
		},
		{
			name:          "user 3 request",
			requestCount:  99,
			tenantID:      "mcp-a2a-3",
			limit:         100,
			expectAllowed: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := SetUpRedis(t)
			// Reset database for a fresh test environment
			// client.FlushDB(context.Background())
			// defer client.Close()

			limiter := NewFixedWindowLimiter(client.Client, tc.limit, time.Minute)

			var lastCode int
			for i := 0; i < int(tc.requestCount); i++ {
				code := FixedWindowMockRequest(limiter, tc.tenantID) // send request and check limit int method Handler
				lastCode = code
			}
			if tc.expectAllowed {
				assert.Equal(t, http.StatusOK, lastCode)
			} else {
				assert.Equal(t, http.StatusTooManyRequests, lastCode)
			}
		})
	}
}

func TestNewTokenBucketLimiter(t *testing.T) {
	client := SetUpRedis(t)
	limiter := NewTokenBucketLimiter(client.Client, 20, 0.2) // xô chứa tối đa 20 tokens, 5s nạp 1 token

	assert.NotNil(t, limiter)
	assert.Equal(t, 20, limiter.capacity)
	assert.Equal(t, 0.2, limiter.rate)
}

func TestTokenBucketLimiter_Valid(t *testing.T) {
	client := SetUpRedis(t)
	limiter := NewTokenBucketLimiter(client.Client, 20, 0.2)
	ctx := context.WithValue(context.Background(), auth.ContextKeyTenantID, "a2a-mcp-123")

	// Create test handler
	handleCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleCalled = true
		tenantID, err := auth.ExtractTenantID(r.Context())
		assert.NoError(t, err)
		assert.Equal(t, tenantID, "a2a-mcp-123")

		w.WriteHeader(http.StatusOK)
	})

	// init request and store response
	req := httptest.NewRequest("POST", "/mcp", nil).WithContext(ctx)
	rr := httptest.NewRecorder()

	handler := limiter.Handler(testHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, handleCalled, true)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestTokenBucketLimiter_NoTenantID(t *testing.T) {
	client := SetUpRedis(t)
	limiter := NewTokenBucketLimiter(client.Client, 20, 0.2)

	// Create test handler
	handleCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleCalled = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("POST", "/mcp", nil)
	rr := httptest.NewRecorder()

	handler := limiter.Handler(testHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, handleCalled, true)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestTokenBucket_CheckLimit_OneTenantID(t *testing.T) {
	tests := []struct {
		name          string
		requestCount  int64
		capacity      int
		rate          float64
		expectAllowed bool
	}{
		// request đầu tiên => xô rỗng => đổ 20 tokens vào xô => request được chấp nhận và xô còn 19 tokens
		{
			name:          "first request",
			requestCount:  1,
			capacity:      20,
			rate:          0.2,
			expectAllowed: true,
		},
		// request đầu => đổ đầu 20 tokens vào xô và request hết 20 tokens => 20 request đầu tiên được chấp nhận
		{
			name:          "at limit",
			requestCount:  20,
			capacity:      20,
			rate:          0.2,
			expectAllowed: true,
		},
		// request đầu => đổ đầu 20 tokens vào xô và request hết 20 tokens => 20 request đầu tiên được chấp nhận, request thứ 21 => xô rỗng => request bị từ chối
		{
			name:          "exceeded limit",
			requestCount:  21,
			capacity:      20,
			rate:          0.2,
			expectAllowed: false,
		},
		// request đầu => đổ đầu 20 tokens vào xô và request hết 20 tokens => 20 request đầu tiên được chấp nhận, request thứ 21 => xô rỗng => request bị từ chối, sau 5s => xô được nạp lại 1 token => request thứ 22 được chấp nhận và xô còn 0 token
		{
			name:          "exceeded limit then allowed after refill",
			requestCount:  22,
			capacity:      20,
			rate:          0.2,
			expectAllowed: true,
		},
	}

	tenantID := "a2a-mcp-go"
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := SetUpRedis(t)
			// Reset database for a fresh test environment
			client.FlushDB(context.Background())
			defer client.Close()

			limiter := NewTokenBucketLimiter(client.Client, tc.capacity, tc.rate)

			var lastCode int
			for i := 0; i < int(tc.requestCount); i++ {
				code := TokenBucketMockRequest(limiter, tenantID) // send request and check limit int method Handler
				if i == 20 && tc.name == "exceeded limit then allowed after refill" {
					// Simulate waiting for 5 seconds to allow token refill
					time.Sleep(5 * time.Second)
				}
				lastCode = code
			}
			if tc.expectAllowed {
				assert.Equal(t, http.StatusOK, lastCode)
			} else {
				assert.Equal(t, http.StatusTooManyRequests, lastCode)
			}
		})
	}
}

func TestTokenBucket_CheckLimit_ManyTenantID(t *testing.T) {
	tests := []struct {
		name          string
		tenantID      string
		requestCount  int64
		capacity      int
		rate          float64
		expectAllowed bool
	}{
		{
			name:          "user 1 request first",
			requestCount:  1,
			tenantID:      "mcp-a2a-1",
			capacity:      20,
			rate:          0.2,
			expectAllowed: true,
		},
		{
			name:          "user 2 request second",
			requestCount:  19,
			tenantID:      "mcp-a2a-1",
			capacity:      20,
			rate:          0.2,
			expectAllowed: true,
		},
		{
			name:          "user 2 request first",
			requestCount:  20,
			tenantID:      "mcp-a2a-2",
			capacity:      20,
			rate:          0.2,
			expectAllowed: true,
		},
		{
			name:          "user 2 request second",
			requestCount:  1,
			tenantID:      "mcp-a2a-2",
			capacity:      20,
			rate:          0.2,
			expectAllowed: false,
		},
		{
			name:          "user 3 request",
			requestCount:  20,
			tenantID:      "mcp-a2a-3",
			capacity:      20,
			rate:          0.2,
			expectAllowed: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := SetUpRedis(t)
			// Reset database for a fresh test environment
			// client.FlushDB(context.Background())
			// defer client.Close()

			limiter := NewTokenBucketLimiter(client.Client, tc.capacity, tc.rate)

			var lastCode int
			for i := 0; i < int(tc.requestCount); i++ {
				code := TokenBucketMockRequest(limiter, tc.tenantID) // send request and check limit int method Handler
				lastCode = code
			}
			if tc.expectAllowed {
				assert.Equal(t, http.StatusOK, lastCode)
			} else {
				assert.Equal(t, http.StatusTooManyRequests, lastCode)
			}
		})
	}
}
