package domain

import (
	"backend/internal/utils"
	"context"
	"time"

	"github.com/google/uuid"
)

type UsageType string

const (
	UsageBoth     UsageType = "both"
	UsageRetaul   UsageType = "retail"
	UsageInternal UsageType = "internal"
)

type Product struct {
	ID                        uuid.UUID `json:"id" db:"id"`
	Name                      string    `json:"name" db:"name"`
	Category                  string    `json:"category" db:"category"`
	Description               string    `json:"description" db:"description"`
	PriceIn                   float64   `json:"price_in" db:"price_in"`
	PriceOut                  float64   `json:"price_out" db:"price_out"`
	UsageType                 UsageType `json:"usage_type" db:"usage_type"`
	LowStockThresholdRetail   int       `json:"low_stock_threshold_retail" db:"low_stock_threshold_retail"`
	LowStockThresholdInternal int       `json:"low_stock_threshold_internal" db:"low_stock_threshold_internal"`
	CreatedAt                 time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt                 time.Time    `json:"updated_at" db:"updated_at"`
}

type ProductRepo interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, search string, category string, minPrice *float64, maxPrice *float64, usageType *UsageType, page, limit int) ([]*Product, int64, error)
}

type ProductUsecase interface {
	Create(ctx context.Context, product *Product) (*Product, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
	Update(ctx context.Context, id uuid.UUID, product *Product) (*Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, search string, category string, minPrice *float64, maxPrice *float64, usageType *UsageType, page, limit int) (*utils.PaginatedResponse[Product], error)
}
