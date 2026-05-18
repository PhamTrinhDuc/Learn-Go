package domain

import "context"

type OrderItems struct {
	ID        string  `json:"id" db:"id"`
	OrderID   string  `json:"order_id" db:"order_id"`
	ProductID string  `json:"product_id" db:"product_id"`
	Quantity  int     `json:"quantity" db:"quantity"`
	UnitPrice float64 `json:"unit_price" db:"unit_price"`
}

type OrderItemsRepo interface {
	GetByID(ctx context.Context, id string) (*OrderItems, error)
	ListByOrderID(ctx context.Context, orderId string) ([]*OrderItems, error)
	Create(ctx context.Context, orderItem *OrderItems) error
	Update(ctx context.Context, orderItem *OrderItems) error
	Delete(ctx context.Context, id string) error
}

type OrderItemsUsecase interface {
	GetByID(ctx context.Context, id string) (*OrderItems, error)
	ListByOrderID(ctx context.Context, orderId string) ([]*OrderItems, error)
	Create(ctx context.Context, orderItem *OrderItems) error
	Update(ctx context.Context, id string, orderItem *OrderItems) error
	Delete(ctx context.Context, id string) error
}
