package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int32 // tối đa X connect 1 lúc
	MinConns int32 // tối thiếu X connect 1 lúc
}

type DB struct {
	pool *pgxpool.Pool
}

func (cfg *DBConfig) Validate() error {
	// 1. Kiểm tra các trường BẮT BUỘC phải có
	if cfg.Host == "" || cfg.User == "" || cfg.DBName == "" {
		return fmt.Errorf("host, user, and dbname are required")
	}

	// 2. Gán giá trị MẶC ĐỊNH cho các trường tùy chọn nếu chúng trống
	if cfg.Port <= 0 {
		cfg.Port = 5432
	}
	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}
	if cfg.MaxConns <= 0 {
		cfg.MaxConns = 10
	}
	if cfg.MinConns < 0 {
		cfg.MinConns = 2
	}
	return nil
}

func NewDB(ctx context.Context, cfg DBConfig) (*DB, error) {
	// 0.Validate config
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate fields in config: %w", err)
	}

	// 1. Concate connection string
	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s pool_max_conns=%d pool_min_conns=%d",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode, cfg.MaxConns, cfg.MinConns,
	)

	// 2. parse config
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse connection string: %w", err)
	}

	// 3. Config connection pool
	poolConfig.MaxConns = cfg.MaxConns             // override connString
	poolConfig.MinConns = cfg.MinConns             // override connString
	poolConfig.MaxConnLifetime = time.Hour         // thời gian sống tối đa của 1 kết nối
	poolConfig.MaxConnIdleTime = 30 * time.Minute  // thời gian tối đa của 1 kết nối không hoạt động
	poolConfig.HealthCheckPeriod = 1 * time.Minute // thời gian kiểm tra kết nối

	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		return nil
	}
	// 4. Create pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to create connection pool: %w", err)
	}
	// 5. Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("Failed ping to database %w", err)
	}
	// 6. Return DB{pool}
	return &DB{pool: pool}, nil
}

func (db *DB) Close() {
	db.pool.Close()
}

// SetTenantContext sets the tenant ID for row-level security
func (db *DB) setTententContext(ctx context.Context, tx pgx.Tx, tenantID string) error {
	// Note: SET commands don't support parameter binding ($1), so we use fmt.Sprintf
	// The tenantID is validated to be a UUID by the JWT validator, so this is safe
	query := fmt.Sprintf("SET LOCAL app.current_tenant_id = '%s'", tenantID)
	_, err := tx.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to set tenant content: %w", err)
	}
	return nil
}

// BeginTx starts a new transaction with tenant context
func (db *DB) BeginTx(ctx context.Context, tenantID string) (pgx.Tx, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	if err := db.setTententContext(ctx, tx, tenantID); err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("failed to set tenant context: %w", err)
	}
	return tx, nil
}

func main() {
	cfg := Config{
		Host:     "localhost",
		Port:     5433,
		User:     "postgres",
		Password: "postgres",
		DBName:   "vchatbot",
		SSLMode:  "disable",
		MaxConns: 10,
		MinConns: 2,
	}

	ctx := context.Background()

	_, err := NewDB(ctx, cfg)
	if err != nil {
		fmt.Printf("failed to initialize database: %v\n", err)
		return
	}
	fmt.Println("Successfully connected to database!")
	// row := db.pool.QueryRow(ctx, "SELECT * FROM companies")
}
