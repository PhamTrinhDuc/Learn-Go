package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	auth "learn-go/a2a_mcp/mcp-server/internal/auth"
	"learn-go/a2a_mcp/mcp-server/internal/database"
	"learn-go/a2a_mcp/mcp-server/internal/middleware"
	"learn-go/a2a_mcp/mcp-server/internal/observability"
	"learn-go/a2a_mcp/mcp-server/internal/redis"
	"learn-go/a2a_mcp/mcp-server/internal/server"
	"learn-go/a2a_mcp/mcp-server/internal/tools"
	"learn-go/a2a_mcp/pkg/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defaultPort = "8080"

	defaultDBHost = "localhost"
	defaultDBPort = 5433

	defaultRedisHost = "localhost"
	defaultRedisPort = 6379

	defaultRateLimit = 100 // requests per minutes
)

// Config holds application configuration
type Config struct {
	Port          string
	Database      database.DBConfig
	Redis         redis.RedisConfig
	Telemtry      observability.Config
	EnableTracing bool
	EnableMetrics bool
}

func setupAuth() (*auth.JWTValidator, *rsa.PrivateKey, error) {
	keysDir := utils.GetEnvString("DEMO_KEYS_DIR", "./tmp/demo-keys")
	privateKeyPath := keysDir + "/private_key.pem"
	publicKeyPath := keysDir + "/public_key.pem"

	var privateKey *rsa.PrivateKey
	var publicKeyPEM []byte
	var err error

	// 1. Try to load existing private key
	if _, err := os.Stat(privateKeyPath); err == nil {
		privData, err := os.ReadFile(privateKeyPath)
		if err == nil {
			block, _ := pem.Decode(privData)
			if block != nil && block.Type == "RSA PRIVATE KEY" {
				privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
				if err == nil {
					log.Printf("Loaded existing private key from %s", privateKeyPath)
				}
			}
		}
	}

	// 2. If no existing key, generate new one
	if privateKey == nil {
		log.Println("Generating new RSA key pair...")
		privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate private key: %w", err)
		}

		// Save keys for future use
		if err := os.MkdirAll(keysDir, 0755); err != nil {
			log.Printf("Warning: Failed to create keys directory: %v", err)
		} else {
			// Save Private Key
			privateKeyPEM := pem.EncodeToMemory(&pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
			})
			os.WriteFile(privateKeyPath, privateKeyPEM, 0600)

			// Save Public Key
			publicKeyBytes, _ := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
			publicKeyPEM = pem.EncodeToMemory(&pem.Block{
				Type:  "PUBLIC KEY",
				Bytes: publicKeyBytes,
			})
			os.WriteFile(publicKeyPath, publicKeyPEM, 0644)
			log.Printf("New keys saved to %s", keysDir)
		}
	}

	// 3. Ensure we have publicKeyPEM
	if len(publicKeyPEM) == 0 {
		publicKeyBytes, _ := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
		publicKeyPEM = pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		})
	}

	// 4. Create authValidator
	validator, err := auth.NewJWTValidator(
		auth.Config{
			PublicKeyPEM: string(publicKeyPEM),
			Issuer:       "mcp-server-demo",
			Audience:     "mcp-server",
		},
	)
	if err != nil {
		return nil, nil, err
	}

	return validator, privateKey, nil
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
			MaxConns: int32(utils.GetEnvInt("DB_MAX_CONNS", 25)),
			MinConns: int32(utils.GetEnvInt("DB_MIN_CONNS", 5)),
		},
		Redis: redis.RedisConfig{
			Host:     utils.GetEnvString("REDIS_HOST", defaultRedisHost),
			Port:     utils.GetEnvInt("REDIS_PORT", defaultRedisPort),
			Username: utils.GetEnvString("REDIS_USERNAME", "jiyuu"),
			Password: utils.GetEnvString("REDIS_PASSWORD", "a2amcpgo"),
			DB:       utils.GetEnvInt("REDIS_DB", 0),
			PoolSize: utils.GetEnvInt("REDIS_POOL_SIZE", 10),
			MinCons:  utils.GetEnvInt("REDIS_MIN_CONNS", 2),
		},
		Telemtry: observability.Config{
			ServiceName:    utils.GetEnvString("OTEL_SERVICE", "mcp-server"),
			ServiceVersion: utils.GetEnvString("OTEL_VERSION", "1.0.0"),
			Environment:    utils.GetEnvString("ENVIRONMENT", "development"),
			OTLPEndpoint:   utils.GetEnvString("OTEL_EXPORTER_OTLP_ENDPOINT", "jaeger:4318"),
			SamplingRate:   utils.GetEnvFloat("OTEL_TRACES_SAMPLER_ARG", 1.0),
		},
		EnableTracing: utils.GetEnvBool("OTEL_ENABLE_TRACING", true),
		EnableMetrics: utils.GetEnvBool("OTEL_ENABLE_METRICS", true),
	}
}

func main() {
	// 1. Init context and config
	ctx := context.Background()
	cfg := loadConfig()

	// 2. Init Services
	// a. Init Database
	db, err := database.NewDB(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("Connected to Database failed: %v", err)
	}
	defer db.Close()
	log.Println("Database setup complete")

	// b. Init Redis
	redis, err := redis.NewRedis(ctx, cfg.Redis)
	if err != nil {
		log.Fatalf("Connected to Redis failed: %v", err)
	}
	defer redis.Close()
	log.Println("Redis setup complete")

	// c. Init Telemetry
	telemetry, err := observability.NewTelemetry(ctx, cfg.Telemtry)
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

	// d. Init JWT Validator
	jwtValidator, privateKey, err := setupAuth()
	if err != nil {
		log.Fatalf("failed to setup authenticate: %v", err)
	}
	log.Println("Authentication setup complete")

	// e. Generate a demo token for testing
	demoToken, _ := auth.GenerateDemoToken("tenant_id-123", "user_id-123", []string{"read", "write"}, privateKey)
	log.Printf("\n--- DEMO TOKEN FOR POSTMAN ---\nFull header: Bearer %s\n------------------------------\n", demoToken)

	// 3. Registry Tools MCP
	toolRegistry := tools.NewRegistry()
	toolRegistry.Register(tools.NewHybridSearchTool(db))
	toolRegistry.Register(tools.NewSearchTool(db))
	log.Printf("Registered %d tools", len(toolRegistry.List()))

	// 4. Create MCP handler with telemetry
	mcpHandler := server.NewMCPHandler(toolRegistry, telemetry)

	authMiddleware := middleware.NewAuthMiddleware(jwtValidator)
	rateLimiter := middleware.NewFixedWindowLimiter(redis.Client, 100, time.Minute)
	// tracingMiddleware := middleware.NewTracingMiddleware(telemetry)

	// 5. Create HTTP server with middleware stack
	mux := http.NewServeMux()

	// a. Healthcheck Endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		data := map[string]any{
			"message": "OK",
			"status":  "healthy",
			"time":    time.Now().UTC(),
		}
		json.NewEncoder(w).Encode(data)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		data := map[string]any{
			"message": "MCP Server - Welcome to the MCP server!",
			"code":    http.StatusOK,
			"tags":    []string{"MCP", "Tools"},
		}

		json.NewEncoder(w).Encode(data)
	})

	// b. Metrics endpoint for Prometheus (no auth required)
	if cfg.EnableMetrics {
		mux.Handle("/metrics", promhttp.Handler())
		log.Printf("Metrics endpoint: http://localhost:%s/metrics", cfg.Port)
	}

	// c. MCP Endpoint
	mux.Handle("/mcp",
		rateLimiter.Handler(
			authMiddleware.Handler(
				mcpHandler,
			),
		),
	)

	httpServer := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second, // thời gian tối đa chờ để đọc body request từ client
		WriteTimeout: 15 * time.Second, // thời gian tối đa chờ để ghi response
		IdleTimeout:  60 * time.Second, // thời gian tối đa kết nối rảnh rỗi
	}

	// Start Server in goroutines
	go func() {
		log.Printf("Starting HTTP server on :%s", cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", cfg.Port, err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down gracefully...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
