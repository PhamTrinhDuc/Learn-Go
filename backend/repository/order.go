package repository

import (
	"backend/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Get(ctx context.Context) ([]*domain.Orders, error) {
	query := `
		SELECT id, user_id, branch_id, total_amount, payment_method, status, created_at, updated_at
		FROM orders
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Orders
	for rows.Next() {
		var order domain.Orders
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.BranchID,
			&order.TotalAmount,
			&order.PaymentMethod,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}
	return orders, nil
}

// GetByID retrieves an order by its ID
func (r *OrderRepository) GetByID(ctx context.Context, id string) (*domain.Orders, error) {
	query := `
		SELECT id, user_id, branch_id, total_amount, payment_method, status, created_at, updated_at
		FROM orders
		WHERE id = $1
	`
	order := &domain.Orders{}
	var paymentMethod string
	err := r.db.QueryRow(ctx, query, id).Scan(
		&order.ID,
		&order.UserID,
		&order.BranchID,
		&order.TotalAmount,
		&paymentMethod,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	order.PaymentMethod = domain.PaymentMethod(paymentMethod)
	return order, nil
}

// Create inserts a new order
func (r *OrderRepository) Create(ctx context.Context, order *domain.Orders) error {
	if order.ID == "" {
		order.ID = uuid.New().String()
	}

	query := `
		INSERT INTO orders (id, user_id, branch_id, total_amount, payment_method, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		order.ID,
		order.UserID,
		order.BranchID,
		order.TotalAmount,
		string(order.PaymentMethod),
		order.Status,
	).Scan(&order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}

// Update updates an existing order
func (r *OrderRepository) Update(ctx context.Context, order *domain.Orders) error {
	query := `
		UPDATE orders
		SET user_id = $1, branch_id = $2, total_amount = $3, payment_method = $4, status = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at
	`
	err := r.db.QueryRow(ctx, query,
		order.UserID,
		order.BranchID,
		order.TotalAmount,
		string(order.PaymentMethod),
		order.Status,
		order.ID,
	).Scan(&order.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}
	return nil
}

// Delete removes an order by ID
func (r *OrderRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM orders WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("order not found")
	}
	return nil
}

// ListByUserID retrieves orders for a specific user
func (r *OrderRepository) ListByUserID(ctx context.Context, userId string) ([]*domain.Orders, error) {
	query := `
		SELECT id, user_id, branch_id, total_amount, payment_method, status, created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders by user: %w", err)
	}
	defer rows.Close()

	var orders []*domain.Orders
	for rows.Next() {
		order := &domain.Orders{}
		var paymentMethod string
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.BranchID,
			&order.TotalAmount,
			&paymentMethod,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		order.PaymentMethod = domain.PaymentMethod(paymentMethod)
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

// ListByBranchID retrieves orders for a specific branch
func (r *OrderRepository) ListByBranchID(ctx context.Context, branchId string) ([]*domain.Orders, error) {
	query := `
		SELECT id, user_id, branch_id, total_amount, payment_method, status, created_at, updated_at
		FROM orders
		WHERE branch_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, branchId)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders by branch: %w", err)
	}
	defer rows.Close()

	var orders []*domain.Orders
	for rows.Next() {
		order := &domain.Orders{}
		var paymentMethod string
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.BranchID,
			&order.TotalAmount,
			&paymentMethod,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		order.PaymentMethod = domain.PaymentMethod(paymentMethod)
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

// OrderItemsRepository handles database operations for order items
type OrderItemsRepository struct {
	db *pgxpool.Pool
}

func NewOrderItemsRepository(db *pgxpool.Pool) *OrderItemsRepository {
	return &OrderItemsRepository{db: db}
}

// GetByID retrieves an order item by its ID
func (r *OrderItemsRepository) GetByID(ctx context.Context, id string) (*domain.OrderItems, error) {
	query := `
		SELECT id, order_id, product_id, quantity, unit_price
		FROM order_items
		WHERE id = $1
	`
	item := &domain.OrderItems{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&item.OrderID,
		&item.ProductID,
		&item.Quantity,
		&item.UnitPrice,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get order item: %w", err)
	}
	return item, nil
}

// ListByOrderID retrieves all items for a specific order
func (r *OrderItemsRepository) ListByOrderID(ctx context.Context, orderId string) ([]*domain.OrderItems, error) {
	query := `
		SELECT id, order_id, product_id, quantity, unit_price
		FROM order_items
		WHERE order_id = $1
	`
	rows, err := r.db.Query(ctx, query, orderId)
	if err != nil {
		return nil, fmt.Errorf("failed to list order items: %w", err)
	}
	defer rows.Close()

	var items []*domain.OrderItems
	for rows.Next() {
		item := &domain.OrderItems{}
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.UnitPrice,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// Create inserts a new order item
func (r *OrderItemsRepository) Create(ctx context.Context, item *domain.OrderItems) error {
	if item.ID == "" {
		item.ID = uuid.New().String()
	}
	query := `
		INSERT INTO order_items (id, order_id, product_id, quantity, unit_price)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query,
		item.ID,
		item.OrderID,
		item.ProductID,
		item.Quantity,
		item.UnitPrice,
	)
	if err != nil {
		return fmt.Errorf("failed to create order item: %w", err)
	}
	return nil
}

// Update updates an existing order item
func (r *OrderItemsRepository) Update(ctx context.Context, item *domain.OrderItems) error {
	query := `
		UPDATE order_items
		SET order_id = $1, product_id = $2, quantity = $3, unit_price = $4
		WHERE id = $5
	`
	_, err := r.db.Exec(ctx, query,
		item.OrderID,
		item.ProductID,
		item.Quantity,
		item.UnitPrice,
		item.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update order item: %w", err)
	}
	return nil
}

// Delete removes an order item by ID
func (r *OrderItemsRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM order_items WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete order item: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("order item not found")
	}
	return nil
}
