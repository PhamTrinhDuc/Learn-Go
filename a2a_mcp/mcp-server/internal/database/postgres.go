package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
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

type Document struct {
	ID        string                 `json:"id"`
	TenantID  string                 `json:"tenant_id"`
	Title     string                 `json:"title"`
	Content   string                 `json:"content"`
	Metadata  map[string]interface{} `json:"metadata"`
	Embedding []float32              `json:"embedding,omitempty"`
	CreateAt  time.Time              `json:"create_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	CreateBy  *string                `json:"create_by,omitempty"` // use pointer to handle NULL
}

type DocumentRAG struct {
	Document Document
	Score    float64
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

// Search document using metadata (title, content, metadata match with query)
func (db *DB) SearchDocuments(
	ctx context.Context,
	tenantID string,
	query string,
	limit int) ([]*Document, error) {
	tx, err := db.BeginTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	searchQuery := `
		SELECT id, tenant_id, title, content, metadata, create_at, updated_at, created_by
		FROM documents 
		WHERE 
			title ILIKE $1 OR 
			content ILIKE $1 OR 
			metadata::text ILIKE $1 
		ORDER BY created_at DESC 
		LIMIT $2
	`
	searchPattern := "%" + query + "%"
	rows, err := tx.Query(ctx, searchQuery, searchPattern, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search documents: %w", err)
	}

	defer rows.Close()

	var documents []*Document
	for rows.Next() {
		doc := &Document{}
		err := rows.Scan(
			&doc.ID,
			&doc.TenantID,
			&doc.Title,
			&doc.Content,
			&doc.Metadata,
			&doc.CreateAt,
			&doc.UpdatedAt,
			&doc.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan document: %w", err)
		}
		documents = append(documents, doc)
	}
	return documents, nil
}

// VectorSearch performs similarity search using pgvector
func (db *DB) VectorSearch(
	ctx context.Context,
	tenantID string,
	query string,
	embedding []float32,
	limit int,
) ([]*DocumentRAG, error) {
	tx, err := db.BeginTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query = `
		SELECT id, tenant_id, title, content, metadata, create_at, updated_at, created_by, 
		1 - (embedding <=>$1) AS similarity_score
		FROM documents 
		WHERE embedding IS NOT NULL 
		ORDER BY embedding <=> $1
		LIMIT $2
	`

	vec := pgvector.NewVector(embedding)
	rows, err := tx.Query(ctx, query, vec, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to perform vector search: %w", err)
	}

	defer rows.Close()

	var documents []*DocumentRAG

	for rows.Next() {
		doc := &Document{}
		var score float64
		var dbEmbedding pgvector.Vector

		err := rows.Scan(
			&doc.ID,
			&doc.TenantID,
			&doc.Title,
			&doc.Content,
			&doc.Metadata,
			&dbEmbedding,
			&doc.CreateAt,
			&doc.UpdatedAt,
			&doc.UpdatedAt,
			&score,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to Scan documents: %w", err)
		}
		doc.Embedding = dbEmbedding.Slice()
		documents = append(documents, &DocumentRAG{
			Document: *doc,
			Score:    score,
		})
	}
	return documents, nil
}

// ListDocuments lists all documents for a tenant
func (db *DB) ListDocuments(ctx context.Context, tenantID string, limit int, offset int) ([]*Document, error) {
	tx, err := db.BeginTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := `
		SELECT id, tenant_id, title, content, metadata, create_at, updated_at, created_by
		FROM documents 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}

	defer rows.Close()

	documents := []*Document{}
	for rows.Next() {
		doc := &Document{}
		err := rows.Scan(
			&doc.ID,
			&doc.TenantID,
			&doc.Title,
			&doc.Content,
			&doc.Metadata,
			&doc.CreateAt,
			&doc.UpdatedAt,
			&doc.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to Scan documents: %w", err)
		}
		documents = append(documents, doc)
	}
	return documents, nil
}

// UpdateDocument updates an existing document
func (db *DB) UpdateDocument(ctx context.Context, tenantID string, doc *Document) error {
	tx, err := db.BeginTx(ctx, tenantID)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		UPDATE documents
		SET title = $1, content = $2, metadata = $3, embedding = $4
		WHERE id = $5
		RETURNING updated_at
	`
	err = tx.QueryRow(ctx, query, doc.Title, doc.Content, doc.Metadata, doc.Embedding, doc.ID).Scan(&doc.UpdatedAt)
	if err == pgx.ErrNoRows {
		return fmt.Errorf("document not found")
	}
	if err != nil {
		return nil
	}
	return tx.Commit(ctx)
}

// DeleteDocument deletes a document by ID
func (db *DB) DeleteDocument(ctx context.Context, tenantID string, docID string) error {
	tx, err := db.BeginTx(ctx, tenantID)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `DELETE FROM documents WHERE id = $1`
	result, err := tx.Exec(ctx, query, docID)
	if err != nil {
		return nil
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("document not found")
	}
	return tx.Commit(ctx)

}

func (db *DB) GetTenantSettings(ctx context.Context, tenantID string) (map[string]interface{}, error) {
	query := `SELECT settings FROM tenants WHERE id = $1 AND is_active = true`

	var settings map[string]interface{}
	err := db.pool.QueryRow(ctx, query, tenantID).Scan(&settings)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("tenant not found or inactive")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant settings: %w", err)
	}

	return settings, nil
}

func main() {
	cfg := DBConfig{
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
