package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"tmdt-backend/domain"
)

type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *pgxpool.Pool) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) scanUser(row pgx.Row) (*domain.User, error) {
	var user domain.User
	err := row.Scan(
		&user.ID,
		&user.Password,
		&user.FullName,
		&user.Email,
		&user.NumPhone,
		&user.Role,
		&user.IsLock,
		&user.JoinedDate,
		&user.Gender,
		&user.DOB,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FetchAll(ctx context.Context) ([]domain.User, error) {
	rows, err := r.db.Query(ctx, "SELECT id, password, full_name, email, num_phone, role, is_lock, joined_date, gender, dob FROM users ORDER BY joined_date ASC")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all users: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		user, err := r.scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		if user != nil {
			users = append(users, *user)
		}
	}

	return users, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	row := r.db.QueryRow(ctx, "SELECT id, password, full_name, email, num_phone, role, is_lock, joined_date, gender, dob FROM users WHERE id = $1 LIMIT 1", id)
	user, err := r.scanUser(row)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return user, nil
}

func (r *userRepository) GetByEmailOrPhone(ctx context.Context, username string) (*domain.User, error) {
	row := r.db.QueryRow(ctx, "SELECT id, password, full_name, email, num_phone, role, is_lock, joined_date, gender, dob FROM users WHERE num_phone = $1 OR email = $2 LIMIT 1", username, username)
	user, err := r.scanUser(row)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email/phone: %w", err)
	}
	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := r.db.QueryRow(ctx, "SELECT id, password, full_name, email, num_phone, role, is_lock, joined_date, gender, dob FROM users WHERE email = $1 LIMIT 1", email)
	user, err := r.scanUser(row)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

func (r *userRepository) GetByPhone(ctx context.Context, phone string) (*domain.User, error) {
	row := r.db.QueryRow(ctx, "SELECT id, password, full_name, email, num_phone, role, is_lock, joined_date, gender, dob FROM users WHERE num_phone = $1 LIMIT 1", phone)
	user, err := r.scanUser(row)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by phone: %w", err)
	}
	return user, nil
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (password, full_name, email, num_phone, role, is_lock, joined_date, gender, dob)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, joined_date
	`
	joinedDate := time.Now()
	err := r.db.QueryRow(
		ctx,
		query,
		user.Password,
		user.FullName,
		user.Email,
		user.NumPhone,
		user.Role,
		user.IsLock,
		joinedDate,
		user.Gender,
		user.DOB,
	).Scan(&user.ID, &user.JoinedDate)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET full_name = $1, email = $2, num_phone = $3, dob = $4, gender = $5
		WHERE id = $6
	`
	cmd, err := r.db.Exec(ctx, query, user.FullName, user.Email, user.NumPhone, user.DOB, user.Gender, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user profile: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("user not found for update")
	}
	return nil
}

func (r *userRepository) UpdateRole(ctx context.Context, id int, role string) error {
	query := `
		UPDATE users
		SET role = $1
		WHERE id = $2
	`
	cmd, err := r.db.Exec(ctx, query, role, id)
	if err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("user not found for role update")
	}
	return nil
}

func (r *userRepository) UpdateStatus(ctx context.Context, id int, isLock bool) error {
	query := `
		UPDATE users
		SET is_lock = $1
		WHERE id = $2
	`
	cmd, err := r.db.Exec(ctx, query, isLock, id)
	if err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("user not found for status update")
	}
	return nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, email string, hashedPassword string) error {
	query := `
		UPDATE users
		SET password = $1
		WHERE email = $2
	`
	cmd, err := r.db.Exec(ctx, query, hashedPassword, email)
	if err != nil {
		return fmt.Errorf("failed to update user password: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("user not found for password update")
	}
	return nil
}

func (r *userRepository) GetDistinctRoles(ctx context.Context) ([]string, int, error) {
	rows, err := r.db.Query(ctx, "SELECT DISTINCT role FROM users")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch distinct roles: %w", err)
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, 0, fmt.Errorf("failed to scan role row: %w", err)
		}
		roles = append(roles, role)
	}

	var count int
	err = r.db.QueryRow(ctx, "SELECT COUNT(DISTINCT role) FROM users").Scan(&count)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count distinct roles: %w", err)
	}

	return roles, count, nil
}
