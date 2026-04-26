package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	auth_internal "mcp-server/internal/auth"
	"mcp-server/internal/database"
	"mcp-server/internal/middleware"
	"mcp-server/internal/observability"
	"mcp-server/internal/redis"
	"mcp-server/internal/server"
	"mcp-server/internal/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defaultPort      = "8080"
	defaultDBHost    = "localhost"
	defaultDBPort    = 5433
	defaultRedisHost = "localhost"
	defaultRedisPort = 6379
)

type Config struct {
	Port          string
	Database      database.DBConfig
	Redis         redis.RedisConfig
	Telemetry     observability.Config
	EnableTracing bool
	EnableMetrics bool
}

func setupAuth() (*auth_internal.JWTValidator, *rsa.PrivateKey, error) {
	keysDir := utils.GetEnvString("DEMO_KEYS_DIR", "./tmp/demo-keys")
	privateKeyPath := keysDir + "/private_key.pem"
	privData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, nil, err
	}
	block, _ := pem.Decode(privData)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, nil, err
	}

	publicKeyPath := keysDir + "/public_key.pem"
	pubData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, nil, err
	}

	validator, err := auth_internal.NewJWTValidator(
		auth_internal.Config{
			PublicKeyPEM: string(pubData),
			Issuer:       "mcp-server-demo",
			Audience:     "mcp-server",
		},
	)
	return validator, privateKey, err
}

func loadConfig() Config {
	return Config{
		Port: utils.GetEnvString("PORT", defaultPort),
		Database: database.DBConfig{
			Host:     utils.GetEnvString("DB_HOST", defaultDBHost),
			Port:     utils.GetEnvInt("DB_PORT", defaultDBPort),
			User:     utils.GetEnvString("DB_USER", "mcp_user"),
			Password: utils.GetEnvString("DB_PASSWORD", "mcp_password"),
			DBName:   utils.GetEnvString("DB_NAME", "mcp_db"),
			SSLMode:  utils.GetEnvString("DB_SSLMODE", "disable"),
		},
		Redis: redis.RedisConfig{
			Host:     utils.GetEnvString("REDIS_HOST", defaultRedisHost),
			Port:     utils.GetEnvInt("REDIS_PORT", defaultRedisPort),
			Password: utils.GetEnvString("REDIS_PASSWORD", "a2amcpgo"),
		},
		Telemetry: observability.Config{
			ServiceName: "mcp-server",
		},
		EnableTracing: true,
		EnableMetrics: true,
	}
}

func main() {
	ctx := context.Background()
	cfg := loadConfig()

	db, err := database.NewDB(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("DB fail: %v", err)
	}
	defer db.Close()

	redis, err := redis.NewRedis(ctx, cfg.Redis)
	if err != nil {
		log.Fatalf("Redis fail: %v", err)
	}
	defer redis.Close()

	jwtValidator, _, err := setupAuth()
	if err != nil {
		log.Printf("Warning: Auth setup failed: %v", err)
	}

	r := gin.Default()

	// 1. Health & Metrics
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	if cfg.EnableMetrics {
		r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}

	// 2. MCP Endpoint with Gin
	authMid := middleware.NewAuthMiddleware(jwtValidator)
	mcpHandler := server.NewSSEHandler(db)

	mcpGroup := r.Group("/mcp")
	if jwtValidator != nil {
		mcpGroup.Use(authMid.Handler())
	}
	{
		// Official SDK SSE transport needs wildcards for session handling
		mcpGroup.Any("/*path", gin.WrapH(mcpHandler))
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("MCP Server starting on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}
