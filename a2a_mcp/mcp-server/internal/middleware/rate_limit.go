package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type FixedWindowLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

func NewFixedWindowLimiter(client *redis.Client, limit int, window time.Duration) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		client: client,
		limit:  limit,
		window: window,
	}
}

func (l *FixedWindowLimiter) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Rate limiting logic would go here
		c.Next()
	}
}


