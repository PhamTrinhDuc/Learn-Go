package middleware

import (
	"backend/internal/auth"
	"backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	validator *auth.JWTValidator
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(validator *auth.JWTValidator) *AuthMiddleware {
	return &AuthMiddleware{
		validator: validator,
	}
}

// Handler validates the JWT token and adds claims to the context
func (m *AuthMiddleware) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, utils.HTTPResponse{Error: "Authorization header required"})
			ctx.Abort()
			return
		}

		// Validation token
		claims, err := m.validator.ValidateToken(authHeader)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, utils.HTTPResponse{Error: "Invalid token: " + err.Error()})
			ctx.Abort()
			return
		}

		// Add auth to context
		ctx.Set(string(auth.ContextKeyUserID), claims.UserID)
		ctx.Set(string(auth.ContextKeyScopes), claims.Scopes)

		ctx.Next()
	}
}

// OptionalHandler validates the token if present, but doesn't block if missing
func (m *AuthMiddleware) OptionalHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader != "" {
			claims, err := m.validator.ValidateToken(authHeader)
			if err == nil {
				ctx.Set(string(auth.ContextKeyUserID), claims.UserID)
				ctx.Set(string(auth.ContextKeyScopes), claims.Scopes)
			}
		}
		ctx.Next()
	}
}
