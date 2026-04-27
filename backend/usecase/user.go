package usecase

import (
	"backend/domain"
	"backend/internal/utils"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	repo domain.UserRepository
}

func NewUserUsecase(repo domain.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

// Create validates and creates a new user
func (u *UserUsecase) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	if user == nil {
		return nil, utils.ErrNilInput
	}

	if err := utils.CombineErrors(
		utils.ValidateRequired("email", user.Email),
		utils.ValidateRequired("password", user.Password),
		utils.ValidateRequired("name", user.Name),
		utils.ValidateRequired("phone", user.Phone),
	); err != nil {
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)
	user.IsActive = true

	// Default role if not specified
	if user.Role == "" {
		user.Role = domain.RoleCustomer
	}

	return u.repo.Create(ctx, user)
}

// GetByID retrieves a user by ID
func (u *UserUsecase) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	if id == uuid.Nil {
		return nil, utils.ErrInvalidID
	}
	return u.repo.GetByID(ctx, id)
}

// Update validates and updates a user
func (u *UserUsecase) Update(ctx context.Context, id uuid.UUID, user *domain.User) (*domain.User, error) {
	if id == uuid.Nil {
		return nil, utils.ErrInvalidID
	}
	if user == nil {
		return nil, utils.ErrNilInput
	}

	// Fetch existing user
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update allowed fields only
	if user.Name != "" {
		existing.Name = user.Name
	}
	if user.Phone != "" {
		existing.Phone = user.Phone
	}
	if user.Email != "" {
		existing.Email = user.Email
	}
	if user.Birthday != nil {
		existing.Birthday = user.Birthday
	}
	if user.Address != nil {
		existing.Address = user.Address
	}
	if user.PreferredBranchID != nil {
		existing.PreferredBranchID = user.PreferredBranchID
	}
	existing.LoyaltyPoints = user.LoyaltyPoints
	existing.IsActive = user.IsActive

	// Hash password if provided
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		existing.Password = string(hashedPassword)
	}

	return u.repo.Update(ctx, id, existing)
}

// Delete deletes a user
func (u *UserUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return utils.ErrInvalidID
	}
	return u.repo.Delete(ctx, id)
}

// List retrieves a list of users with pagination
func (u *UserUsecase) List(ctx context.Context, page, limit int) (*utils.PaginatedResponse[domain.User], error) {
	page, limit = utils.NormalizePagination(page, limit)
	users, total, err := u.repo.List(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	return utils.CreatePaginatedResponse(users, total, page, limit), nil
}

// ListByRole retrieves users filtered by role with pagination
func (u *UserUsecase) ListByRole(ctx context.Context, role domain.UserRole, page, limit int) (*utils.PaginatedResponse[domain.User], error) {
	
	page, limit = utils.NormalizePagination(page, limit)
	users, total, err := u.repo.ListByRole(ctx, role, page, limit)
	if err != nil {
		return nil, err
	}
	return utils.CreatePaginatedResponse(users, total, page, limit), nil
}

// Authenticate validates credentials and returns user with JWT token
func (u *UserUsecase) Authenticate(ctx context.Context, email, password string) (*domain.User, string, error) {
	if err := utils.CombineErrors(
		utils.ValidateRequired("email", email),
		utils.ValidateRequired("password", password),
	); err != nil {
		return nil, "", err
	}

	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := generateJWTToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Update last visit
	user.LastVisitAt = timePtr(time.Now())
	u.repo.Update(ctx, user.ID, user)

	return user, token, nil
}

// generateJWTToken creates a JWT token for the user
func generateJWTToken(user *domain.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-this-in-production"
	}

	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Helper function to get time pointer
func timePtr(t time.Time) *time.Time {
	return &t
}
