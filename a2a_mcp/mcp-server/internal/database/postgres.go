package database

import (
	"context"
	"fmt"
	"strings"
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
	CreatedAt time.Time              `json:"creatd_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	CreatedBy *string                `json:"created_by,omitempty"` // use pointer to handle NULL
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

// InsertDocument inserts a new document
func (db *DB) InsertDocument(ctx context.Context, tenantID string, doc *Document) error {
	tx, err := db.BeginTx(ctx, tenantID)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO documents (tenant_id, title, content, metadata, embedding, created_by) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	var embedding interface{}
	if doc.Embedding != nil {
		embedding = pgvector.NewVector(doc.Embedding)
	}
	err = tx.QueryRow(
		ctx,
		query,
		doc.TenantID,
		doc.Title,
		doc.Content,
		doc.Metadata,
		embedding,
		doc.CreatedBy,
	).Scan(&doc.ID, &doc.CreatedAt, &doc.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert document to DB: %w", err)
	}
	return tx.Commit(ctx)
}

// InsertBatchingDocument inserts batching new documents
func (db *DB) InsertBatchingDocument(ctx context.Context, tenantID string, docs []*Document, numFields int) error {
	if len(docs) == 0 {
		return fmt.Errorf("Find 0 documents. Documents is required")
	}

	tx, err := db.BeginTx(ctx, tenantID)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Chuẩn bị các thành phần cho câu query
	index := 1
	valueStrings := make([]string, 0, len(docs))
	valueArgs := make([]interface{}, 0, numFields)

	for _, doc := range docs {
		placeholders := make([]string, numFields)
		for i := range placeholders {
			placeholders[i] = fmt.Sprintf("$%d", index)
			index++
		}
		valueStrings = append(valueStrings, "("+strings.Join(placeholders, ", ")+")")
		var embedding interface{}
		if doc.Embedding != nil {
			embedding = pgvector.NewVector(doc.Embedding)
		}

		valueArgs = append(valueArgs,
			doc.TenantID,
			doc.Title,
			doc.Content,
			doc.Metadata,
			embedding,
			doc.CreatedBy,
		)
	}
	// 2. Ghép nối câu query hoàn chỉnh
	query := fmt.Sprintf(`
		INSERT INTO documents (tenant_id, title, content, metadata, embedding, created_by)
		VALUES %s
		RETURNING id, created_at, updated_at`,
		strings.Join(valueStrings, ","),
	)

	// 3. Thực thi
	rows, err := tx.Query(ctx, query, valueArgs...)
	if err != nil {
		return fmt.Errorf("Error during perform batching query: %w", err)
	}

	defer rows.Close()

	// 4. Scan kết quả ngược lại cho từng Document
	i := 0
	for rows.Next() {
		if err := rows.Scan(&docs[i].ID, &docs[i].CreatedAt, &docs[i].UpdatedAt); err != nil {
			return err
		}
		i++
	}
	return tx.Commit(ctx)
}

// Search document using metadata (title, content, metadata match with query)
func (db *DB) SearchDocuments(ctx context.Context, tenantID string, query string, limit int) ([]*Document, error) {
	tx, err := db.BeginTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	searchQuery := `
		SELECT id, tenant_id, title, content, metadata, created_at, updated_at, created_by
		FROM documents 
		WHERE 
			to_tsvector('simple', f_unaccent(title || ' ' || content)) @@ websearch_to_tsquery('simple', unaccent($1))
		ORDER BY ts_rank_cd(to_tsvector('simple', f_unaccent(title || ' ' || content)), websearch_to_tsquery('simple', unaccent($1))) DESC
		LIMIT $2
	`
	rows, err := tx.Query(ctx, searchQuery, query, limit)
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
			&doc.CreatedAt,
			&doc.UpdatedAt,
			&doc.CreatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan document: %w", err)
		}
		documents = append(documents, doc)
	}
	return documents, nil
}

func (db *DB) GetDocument(ctx context.Context, tenantID, docID string) (*Document, error) {
	tx, err := db.BeginTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := `
		SELECT id, tenant_id, title, content, metadata, embedding, created_at, updated_at, created_by
		FROM documents
		WHERE id = $1
	`
	doc := &Document{}
	var dbEmbedding *pgvector.Vector
	err = tx.QueryRow(ctx, query, docID).Scan(
		&doc.ID, &doc.TenantID, &doc.Title, &doc.Content, &doc.Metadata,
		&dbEmbedding, &doc.CreatedAt, &doc.UpdatedAt, &doc.CreatedBy,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("document not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	if dbEmbedding != nil {
		doc.Embedding = dbEmbedding.Slice()
	}

	return doc, nil
}

// ListDocuments lists all documents for a tenant
func (db *DB) ListDocuments(ctx context.Context, tenantID string, limit int, offset int) ([]*Document, error) {
	tx, err := db.BeginTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := `
		SELECT id, tenant_id, title, content, metadata, created_at, updated_at, created_by
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
			&doc.CreatedAt,
			&doc.UpdatedAt,
			&doc.CreatedBy,
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
func (db *DB) DeleteDocumentByID(ctx context.Context, tenantID string, docID string) error {
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

// DeleteDocument deletes a document by ID
func (db *DB) DeleteDocumentByTenantID(ctx context.Context, tenantID string) error {
	tx, err := db.BeginTx(ctx, tenantID)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `DELETE FROM documents WHERE tenant_id = $1`
	result, err := tx.Exec(ctx, query, tenantID)
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
