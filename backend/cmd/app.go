package main

import (
	"backend/api/middleware"
	"backend/api/route"
	"backend/bootstrap"
	"backend/internal/auth"
	"backend/internal/observability"
	"backend/internal/utils"
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	defaultPort      = "8080"
	defaultDBHost    = "localhost"
	defaultDBPort    = 5433
	defaultRedisHost = "localhost"
	defaultRedisPort = 6379
)

type Config struct {
	Port     string
	Database bootstrap.DBConfig
	// Redis         redis.RedisConfig
	Telemetry     observability.Config
	EnableTracing bool
	EnableMetrics bool
	BackendURL    string
}

func loadConfig() Config {
	return Config{
		Port: utils.GetEnvString("PORT", defaultPort),
		Database: bootstrap.DBConfig{
			Host:     utils.GetEnvString("DB_HOST", defaultDBHost),
			Port:     utils.GetEnvInt("DB_PORT", defaultDBPort),
			User:     utils.GetEnvString("DB_USER", "mcp_user"),
			Password: utils.GetEnvString("DB_PASSWORD", "mcp_password"),
			DBName:   utils.GetEnvString("DB_NAME", "salon_chain"),
			SSLMode:  utils.GetEnvString("DB_SSLMODE", "disable"),
			MaxConns: int32(utils.GetEnvInt("DB_MAX_CONNS", 25)),
			MinConns: int32(utils.GetEnvInt("DB_MIN_CONNS", 5)),
		},

		// Redis: redis.RedisConfig{
		// 	Host:     utils.GetEnvString("REDIS_HOST", defaultRedisHost),
		// 	Port:     utils.GetEnvInt("REDIS_PORT", defaultRedisPort),
		// 	Username: utils.GetEnvString("REDIS_USERNAME", "jiyuu"),
		// 	Password: utils.GetEnvString("REDIS_PASSWORD", "a2amcpgo"),
		// 	DB:       utils.GetEnvInt("REDIS_DB", 0),
		// 	PoolSize: utils.GetEnvInt("REDIS_POOL_SIZE", 10),
		// 	MinCons:  utils.GetEnvInt("REDIS_MIN_CONNS", 2),
		// },

		Telemetry: observability.Config{
			ServiceName:    utils.GetEnvString("OTEL_SERVICE", "backend-server"),
			ServiceVersion: utils.GetEnvString("OTEL_VERSION", "1.0.0"),
			Environment:    utils.GetEnvString("ENVIRONMENT", "development"),
			OTLPEndpoint:   utils.GetEnvString("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4318"),
			SamplingRate:   utils.GetEnvFloat("OTEL_TRACES_SAMPLER_ARG", 1.0),
			EnableTracing:  utils.GetEnvBool("OTEL_ENABLE_TRACING", true),
			EnableMetrics:  utils.GetEnvBool("OTEL_ENABLE_METRICS", true),
		},
	}
}

func main() {
	// Khởi tạo Gin engine
	r := gin.Default()

	// Tạo context lắng nghe tín hiệu từ hệ điều hành
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := loadConfig()
	// 1. init database
	db, err := bootstrap.NewDB(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("Không thể kết nối DB: %v", err)
	}

	// 2. Init telemetry
	telemetry, err := observability.NewTelemetry(ctx, cfg.Telemetry)
	if err != nil {
		log.Fatalf("Failed to initialize telemetry: %v", err)
	}

	// tạo context mới để shutdown thay vì shutdown trực tiếp để tránh cancel context chính của app
	defer func() {
		shudownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := telemetry.Shutdown(shudownCtx); err != nil {
			log.Fatalf("Error shutting down telemetry: %v", err)
		}
	}()
	log.Println("OpenTelemetry initialized successfully")

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware() // Khởi tạo middleware của bạn
	tracingMiddleware := middleware.NewTracingMiddleware(telemetry)

	r.Use(tracingMiddleware.Handler())

	// 4. Init RSA keys for testing
	keyDir := filepath.Join("tmp", "demo-keys")
	if err := auth.EnsureKeysExist(keyDir); err != nil {
		log.Printf("Warning: Failed to ensure keys: %v", err)
	} else {
		log.Printf("RSA keys ensured in %s", keyDir)
		// Generate a test token for debugging
		testToken, err := auth.GenerateTokenWithPrivateKey("admin-123", "admin@example.com", "admin")
		if err == nil {
			log.Printf("\n--- TEST TOKEN (COPY EVERYTHING BETWEEN LINES) ---\n%s\n--- END TOKEN ---\n\n", testToken)
		} else {
			log.Printf("Warning: failed to generate token: %v", err)
		}
	}

	//  Group các API lại (ví dụ prefix /api/v1)
	container := bootstrap.NewContainer(db.GetPool(), telemetry)

	v1 := r.Group("/api/v1")
	route.SetupBranchRoutes(v1, container.BranchCtl, authMiddleware)
	route.SetupUserRoutes(v1, container.UserCtl, authMiddleware)
	route.SetupStylistRoutes(v1, container.StylistCtl, authMiddleware)
	route.SetupStylistScheduleRoutes(v1, container.StylistScheduleCtl, authMiddleware)

	// 6. Chạy Server
	r.Run(":8080") // Server sẽ chạy tại http://localhost:8080
}
