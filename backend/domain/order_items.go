package domain

type OrderItems struct {
	ID        string  `json:"id" db:"id"`
	OrderID   string  `json:"order_id" db:"order_id"`
	ProductID string  `json:"product_id" db:"product_id"`
	Quantity  int     `json:"quantity" db:"quantity"`
	UnitPrice float64 `json:"unit_price" db:"unit_price"`
}

type OrderItemsRepo interface {
	GetByID(id string) (*OrderItems, error)
	ListByOrderID(orderId string) ([]*OrderItems, error)
	Create(orderItem *OrderItems) error
	Update(orderItem *OrderItems) error
	Delete(id string) error
}
