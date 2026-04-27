package repository

import (
	"backend/domain"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user
func (r *UserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `
		INSERT INTO users (id, email, password, name, phone, birthday, address, role, loyalty_points, preferred_branch_id, last_visit_at, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, email, name, phone, birthday, address, role, loyalty_points, preferred_branch_id, last_visit_at, is_active, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		user.ID, user.Email, user.Password, user.Name, user.Phone,
		user.Birthday, user.Address, user.Role, user.LoyaltyPoints,
		user.PreferredBranchID, user.LastVisitAt, user.IsActive,
		user.CreatedAt, user.UpdatedAt,
	).Scan(
		&user.ID, &user.Email, &user.Name, &user.Phone,
		&user.Birthday, &user.Address, &user.Role, &user.LoyaltyPoints,
		&user.PreferredBranchID, &user.LastVisitAt, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user := &domain.User{}

	query := `
		SELECT id, email, password, name, phone, birthday, address, role, loyalty_points, preferred_branch_id, last_visit_at, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name, &user.Phone,
		&user.Birthday, &user.Address, &user.Role, &user.LoyaltyPoints,
		&user.PreferredBranchID, &user.LastVisitAt, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}

	query := `
		SELECT id, email, password, name, phone, birthday, address, role, loyalty_points, preferred_branch_id, last_visit_at, is_active, created_at, updated_at
		FROM users
		WHERE email = $1 AND is_active = true
	`

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name, &user.Phone,
		&user.Birthday, &user.Address, &user.Role, &user.LoyaltyPoints,
		&user.PreferredBranchID, &user.LastVisitAt, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, id uuid.UUID, user *domain.User) (*domain.User, error) {
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users
		SET email = $2, name = $3, phone = $4, birthday = $5, address = $6, role = $7, 
		    loyalty_points = $8, preferred_branch_id = $9, last_visit_at = $10, is_active = $11, updated_at = $12
		WHERE id = $1
		RETURNING id, email, name, phone, birthday, address, role, loyalty_points, preferred_branch_id, last_visit_at, is_active, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		id, user.Email, user.Name, user.Phone, user.Birthday, user.Address, user.Role,
		user.LoyaltyPoints, user.PreferredBranchID, user.LastVisitAt, user.IsActive, user.UpdatedAt,
	).Scan(
		&user.ID, &user.Email, &user.Name, &user.Phone,
		&user.Birthday, &user.Address, &user.Role, &user.LoyaltyPoints,
		&user.PreferredBranchID, &user.LastVisitAt, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// Delete soft deletes a user
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET is_active = false, updated_at = $2 WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// List retrieves a list of users with pagination
func (r *UserRepository) List(ctx context.Context, page, limit int) ([]*domain.User, int64, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM users WHERE is_active = true`
	err := r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Get paginated results
	query := `
		SELECT id, email, name, phone, birthday, address, role, loyalty_points, preferred_branch_id, last_visit_at, is_active, created_at, updated_at
		FROM users
		WHERE is_active = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	users := make([]*domain.User, 0, limit)
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.Name, &user.Phone,
			&user.Birthday, &user.Address, &user.Role, &user.LoyaltyPoints,
			&user.PreferredBranchID, &user.LastVisitAt, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating users: %w", err)
	}

	return users, total, nil
}

// ListByRole retrieves users filtered by role with pagination
func (r *UserRepository) ListByRole(ctx context.Context, role domain.UserRole, page, limit int) ([]*domain.User, int64, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM users WHERE role = $1 AND is_active = true`
	err := r.db.QueryRow(ctx, countQuery, role).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Get paginated results
	query := `
		SELECT id, email, name, phone, birthday, address, role, loyalty_points, preferred_branch_id, last_visit_at, is_active, created_at, updated_at
		FROM users
		WHERE role = $1 AND is_active = true
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, role, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users by role: %w", err)
	}
	defer rows.Close()

	users := make([]*domain.User, 0, limit)
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.Name, &user.Phone,
			&user.Birthday, &user.Address, &user.Role, &user.LoyaltyPoints,
			&user.PreferredBranchID, &user.LastVisitAt, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating users: %w", err)
	}

	return users, total, nil
}
