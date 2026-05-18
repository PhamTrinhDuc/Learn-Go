package domain

import "context"

type PaymentMethod string

const (
	PaymentMethodCash   PaymentMethod = "cash"
	PaymentMethodCard   PaymentMethod = "card"
	PaymentMethodOnline PaymentMethod = "online"
)

type Orders struct {
	ID            string        `json:"id" db:"id"`
	UserID        string        `json:"user_id" db:"user_id"`
	BranchID      string        `json:"branch_id" db:"branch_id"`
	TotalAmount   float64       `json:"total_amount" db:"total_amount"`
	PaymentMethod PaymentMethod `json:"payment_method" db:"payment_method"`
	Status        string        `json:"status" db:"status"`
	CreatedAt     string        `json:"created_at" db:"created_at"`
	UpdatedAt     string        `json:"updated_at" db:"updated_at"`
}

type OrdersRepo interface {
	Get(ctx context.Context) ([]*Orders, error)
	GetByID(ctx context.Context, id string) (*Orders, error)
	Create(ctx context.Context, order *Orders) error
	Update(ctx context.Context, order *Orders) error
	Delete(ctx context.Context, id string) error
	ListByUserID(ctx context.Context, userId string) ([]*Orders, error)
	ListByBranchID(ctx context.Context, branchId string) ([]*Orders, error)
}

type OrdersUsecase interface {
	Get(ctx context.Context) ([]*Orders, error)
	GetByID(ctx context.Context, id string) (*Orders, error)
	Create(ctx context.Context, order *Orders) error
	Update(ctx context.Context, id string, order *Orders) error
	Delete(ctx context.Context, id string) error
	ListByUserID(ctx context.Context, userId string) ([]*Orders, error)
	ListByBranchID(ctx context.Context, branchId string) ([]*Orders, error)
}
