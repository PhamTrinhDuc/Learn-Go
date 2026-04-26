package domain

type Inventory struct {
	ID               string `json:"id" db:"id"`
	ProductId        string `json:"product_id" db:"product_id"`
	BranchId         string `json:"branch_id" db:"branch_id"`
	QuantityTotal    int    `json:"quantity_total" db:"quantity_total"`
	QuantityRetail   int    `json:"quantity_retail" db:"quantity_retail"`
	QuantityInternal int    `json:"quantity_internal" db:"quantity_internal"`
	UpdatedAt        string `json:"updated_at" db:"updated_at"`
}

type InventoryRepo interface {
	GetByID(id string) (*Inventory, error)
	GetByProductAndBranch(productId, branchId string) (*Inventory, error)
	Create(inventory *Inventory) error
	Update(inventory *Inventory) error
	Delete(id string) error
	ListByBranchID(branchId string) ([]*Inventory, error)
}
