package middleware

import (
	"mcp-server/internal/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddle struct {
	validator *auth.JWTValidator
}

func NewAuthMiddleware(validator *auth.JWTValidator) *AuthMiddle {
	return &AuthMiddle{
		validator: validator,
	}
}

// Handler is the Gin-compatible auth middleware
func (am *AuthMiddle) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := am.validator.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			c.Abort()
			return
		}

		// Update context with auth claims
		ctx := auth.WithAuth(c.Request.Context(), claims)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
