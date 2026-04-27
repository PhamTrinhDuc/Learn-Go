package controller

import (
	"backend/domain"
	"backend/internal/utils"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserUsecase is a mock implementation of domain.UserUsecase interface
type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserUsecase) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserUsecase) Update(ctx context.Context, id uuid.UUID, user *domain.User) (*domain.User, error) {
	args := m.Called(ctx, id, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserUsecase) List(ctx context.Context, page, limit int) (*utils.PaginatedResponse[domain.User], error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*utils.PaginatedResponse[domain.User]), args.Error(1)
}

func (m *MockUserUsecase) ListByRole(ctx context.Context, role domain.UserRole, page, limit int) (*utils.PaginatedResponse[domain.User], error) {
	args := m.Called(ctx, role, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*utils.PaginatedResponse[domain.User]), args.Error(1)
}

func (m *MockUserUsecase) Authenticate(ctx context.Context, email, password string) (*domain.User, string, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).(*domain.User), args.String(1), args.Error(2)
}

// Test CreateUser with valid data
func TestCreateUser_Success(t *testing.T) {
	mockUC := new(MockUserUsecase)

	userID := uuid.New()
	now := time.Now()
	expectedUser := &domain.User{
		ID:        userID,
		Email:     "test@example.com",
		Name:      "Test User",
		Phone:     "0123456789",
		Role:      domain.RoleCustomer,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockUC.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.Email == "test@example.com" && u.Name == "Test User"
	})).Return(expectedUser, nil)

	uc := &UserController{
		useCase:  mockUC,
		validate: validator.New(),
	}

	req := domain.User{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
		Phone:    "0123456789",
	}

	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq

	uc.CreateUser(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockUC.AssertExpectations(t)
}

// Test Login with correct credentials
func TestLogin_Success(t *testing.T) {
	mockUC := new(MockUserUsecase)

	userID := uuid.New()
	now := time.Now()
	expectedUser := &domain.User{
		ID:        userID,
		Email:     "test@example.com",
		Name:      "Test User",
		Phone:     "0123456789",
		Role:      domain.RoleCustomer,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockUC.On("Authenticate", mock.Anything, "test@example.com", "password123").
		Return(expectedUser, "token123", nil)

	uc := &UserController{
		useCase:  mockUC,
		validate: validator.New(),
	}

	req := domain.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq

	uc.Login(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

// Test GetMe with user context
func TestGetMe_Success(t *testing.T) {
	mockUC := new(MockUserUsecase)

	userID := uuid.New()
	now := time.Now()
	expectedUser := &domain.User{
		ID:        userID,
		Email:     "test@example.com",
		Name:      "Test User",
		Phone:     "0123456789",
		Role:      domain.RoleCustomer,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockUC.On("GetByID", mock.Anything, userID).Return(expectedUser, nil)

	uc := &UserController{
		useCase:  mockUC,
		validate: validator.New(),
	}

	httpReq, _ := http.NewRequest("GET", "/users/me", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq
	c.Set("user_id", userID.String())

	uc.GetMe(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

// Test Register - new user signup
func TestRegister_Success(t *testing.T) {
	mockUC := new(MockUserUsecase)

	userID := uuid.New()
	now := time.Now()
	expectedUser := &domain.User{
		ID:        userID,
		Email:     "newuser@example.com",
		Name:      "New User",
		Phone:     "0123456789",
		Role:      domain.RoleCustomer,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockUC.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.Email == "newuser@example.com" && u.Role == domain.RoleCustomer
	})).Return(expectedUser, nil)

	uc := &UserController{
		useCase:  mockUC,
		validate: validator.New(),
	}

	req := domain.RegisterRequest{
		Email:    "newuser@example.com",
		Password: "password123",
		Name:     "New User",
		Phone:    "0123456789",
	}

	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq

	uc.Register(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockUC.AssertExpectations(t)
}
