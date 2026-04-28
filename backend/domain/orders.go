package domain

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
	GetByID(id string) (*Orders, error)
	Create(order *Orders) error
	Update(order *Orders) error
	Delete(id string) error
	ListByUserID(userId string) ([]*Orders, error)
	ListByBranchID(branchId string) ([]*Orders, error)
}
