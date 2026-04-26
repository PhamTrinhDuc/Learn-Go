package repository

import (
	"backend/domain"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

// Create inserts a new product
func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) error {

	query := `INSERT INTO products (name, description, price_in, price_out, category, usage_type, low_stock_threshold_retail, low_stock_threshold_internal)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	err := r.db.QueryRow(ctx, query,
		product.Name,
		product.Description,
		product.PriceIn,
		product.PriceOut,
		product.Category,
		product.UsageType,
		product.LowStockThresholdRetail,
		product.LowStockThresholdInternal,
	).Scan(&product.ID)

	if err != nil {
		return err
	}
	return nil
}

// GetByID retrieves a product by ID
func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	query := `SELECT id, name, description, price_in, price_out, category, usage_type, low_stock_threshold_retail, low_stock_threshold_internal
							FROM products WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var product domain.Product
	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.PriceIn,
		&product.PriceOut,
		&product.Category,
		&product.UsageType,
		&product.LowStockThresholdRetail,
		&product.LowStockThresholdInternal,
	)

	if err != nil {
		return nil, err
	}
	return &product, nil
}

// List retrieves products with pagination and filters
func (r *ProductRepository) List(ctx context.Context, search string, category string, minPrice *float64, maxPrice *float64, usageType *domain.UsageType, page, limit int) ([]*domain.Product, int64, error) {
	query := `
	SELECT
		id,
		name,
		description,
		price_in,
		price_out,
		category,
		usage_type,
		low_stock_threshold_retail,
		low_stock_threshold_internal,
		created_at,
		updated_at
	FROM products
	WHERE
		($1 = '' OR name ILIKE '%' || $1 || '%' OR category ILIKE '%' || $1 || '%')
		AND ($2 = '' OR category = $2)
		AND ($3::float8 IS NULL OR price_out >= $3)
		AND ($4::float8 IS NULL OR price_out <= $4)
		AND ($5::text IS NULL OR usage_type = $5)
	ORDER BY id DESC
	LIMIT $6 OFFSET $7
	`
	rows, err := r.db.Query(ctx, query, search, category, minPrice, maxPrice, usageType, limit, (page-1)*limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var product domain.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.PriceIn,
			&product.PriceOut,
			&product.Category,
			&product.UsageType,
			&product.LowStockThresholdRetail,
			&product.LowStockThresholdInternal,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, &product)
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM products WHERE
		($1 = '' OR name ILIKE '%' || $1 || '%' OR category ILIKE '%' || $1 || '%')
		AND ($2 = '' OR category = $2)
		AND ($3::float8 IS NULL OR price_out >= $3)
		AND ($4::float8 IS NULL OR price_out <= $4)
		AND ($5::text IS NULL OR usage_type = $5)`

	var total int64
	err = r.db.QueryRow(ctx, countQuery, search, category, minPrice, maxPrice, usageType).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// Update updates a product
func (r *ProductRepository) Update(ctx context.Context, product *domain.Product) error {
	query := `UPDATE products SET name = $1, description = $2, price_in = $3, price_out = $4, category = $5, usage_type = $6, low_stock_threshold_retail = $7, low_stock_threshold_internal = $8 WHERE id = $9`
	_, err := r.db.Exec(ctx, query,
		product.Name,
		product.Description,
		product.PriceIn,
		product.PriceOut,
		product.Category,
		product.UsageType,
		product.LowStockThresholdRetail,
		product.LowStockThresholdInternal,
		product.ID,
	)
	return err
}

// Delete deletes a product
func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
