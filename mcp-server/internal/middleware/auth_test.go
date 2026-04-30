package middleware

import (
	"mcp-server/internal/auth"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware_Handler(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	gin.SetMode(gin.TestMode)

	t.Run("Valid Token", func(t *testing.T) {
		r := gin.New()
		middleware := NewAuthMiddleware()
		r.Use(middleware.Handler())

		r.GET("/test", func(c *gin.Context) {
			userID, _ := c.Get(string(auth.ContextKeyUserID))
			assert.Equal(t, "user123", userID)
			c.String(http.StatusOK, "success")
		})

		token, err := auth.GenerateTokenWithPrivateKey("user123", "test@test.com", "admin")
		require.NoError(t, err)

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "success", w.Body.String())
	})

	t.Run("Missing Token", func(t *testing.T) {
		r := gin.New()
		middleware := NewAuthMiddleware()
		r.Use(middleware.Handler())

		r.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "success")
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		r := gin.New()
		middleware := NewAuthMiddleware()
		r.Use(middleware.Handler())

		r.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "success")
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestAuthMiddleware_OptionalHandler(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	gin.SetMode(gin.TestMode)

	t.Run("Valid Token - Passes Claims", func(t *testing.T) {
		r := gin.New()
		middleware := NewAuthMiddleware()
		r.Use(middleware.OptionalHandler())

		r.GET("/test", func(c *gin.Context) {
			userID, exists := c.Get(string(auth.ContextKeyUserID))
			assert.True(t, exists)
			assert.Equal(t, "user123", userID)
			c.String(http.StatusOK, "success")
		})

		token, err := auth.GenerateTokenWithPrivateKey("user123", "test@test.com", "admin")
		require.NoError(t, err)

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Missing Token - Proceeds Without Claims", func(t *testing.T) {
		r := gin.New()
		middleware := NewAuthMiddleware()
		r.Use(middleware.OptionalHandler())

		r.GET("/test", func(c *gin.Context) {
			_, exists := c.Get(string(auth.ContextKeyUserID))
			assert.False(t, exists)
			c.String(http.StatusOK, "success")
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAuthMiddleware_RequireRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Has Required Role", func(t *testing.T) {
		r := gin.New()
		middleware := NewAuthMiddleware()

		// Mock role in context
		r.Use(func(c *gin.Context) {
			c.Set(string(auth.ContextKeyRole), "admin")
			c.Next()
		})
		r.Use(middleware.RequireRole("admin"))

		r.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "success")
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Insufficient Role", func(t *testing.T) {
		r := gin.New()
		middleware := NewAuthMiddleware()

		// Mock role in context
		r.Use(func(c *gin.Context) {
			c.Set(string(auth.ContextKeyRole), "user")
			c.Next()
		})
		r.Use(middleware.RequireRole("admin"))

		r.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "success")
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Missing Role", func(t *testing.T) {
		r := gin.New()
		middleware := NewAuthMiddleware()

		r.Use(middleware.RequireRole("admin"))

		r.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "success")
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
