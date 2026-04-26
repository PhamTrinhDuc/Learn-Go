package domain

import (
	"backend/internal/utils"
	"context"
)

type Stylist struct {
	ID        string `json:"id" db:"id"`
	BranchID  string `json:"branch_id" db:"branch_id"`
	Name      string `json:"name" db:"name"`
	Phone     string `json:"phone" db:"phone"`
	IsActive  bool   `json:"is_active" db:"is_active"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

type StylistRepository interface {
	Create(ctx context.Context, stylist *Stylist) (*Stylist, error)
	GetByID(ctx context.Context, id string) (*Stylist, error)
	Update(ctx context.Context, id string, stylist *Stylist) (*Stylist, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) ([]*Stylist, int64, error)
	ListByBranch(ctx context.Context, branchID string, page, limit int) ([]*Stylist, int64, error)
}

type StylistUsecase interface {
	Create(ctx context.Context, stylist *Stylist) (*Stylist, error)
	GetByID(ctx context.Context, id string) (*Stylist, error)
	Update(ctx context.Context, id string, stylist *Stylist) (*Stylist, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) (*utils.PaginatedResponse[Stylist], error)
	ListByBranch(ctx context.Context, branchID string, page, limit int) (*utils.PaginatedResponse[Stylist], error)
}
