package middleware

import (
	"auth"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	protocol "protocal"
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
	redis.set("EXPIRE", KEYS[1], ARGV[1])
end
return curr
`

func (rl *FixedWindowLimiter) WindowCheckLimit(ctx context.Context, tenantID string) (bool, error) {
	key := fmt.Sprint("ratelimit:%s%d", tenantID, time.Now().Unix()/(60))

	// Execute script
	// Keys: key
	// Args: rl.window (s)
	result, err := rl.redis.Eval(
		ctx,
		FixedWindowScript,
		[]string{key},
		int(rl.window.Seconds())).Result()
	if err != nil {
		return false, fmt.Errorf("lua script failed: %w", err)
	}
	count, ok := result.(int)
	if !ok {
		return false, fmt.Errorf("unexpected return type from redis: %T", result)
	}
	return count <= rl.defaultLimit, nil
}

func (rl *FixedWindowLimiter) SendError(w http.ResponseWriter, id interface{}, code int, message string) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

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
