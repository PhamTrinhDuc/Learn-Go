package usecase

import (
	"backend/domain"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository is a mock implementation of domain.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, id uuid.UUID, user *domain.User) (*domain.User, error) {
	args := m.Called(ctx, id, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, page, limit int) ([]*domain.User, int64, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*domain.User), int64(args.Int(1)), args.Error(2)
}

func (m *MockUserRepository) ListByRole(ctx context.Context, role domain.UserRole, page, limit int) ([]*domain.User, int64, error) {
	args := m.Called(ctx, role, page, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*domain.User), int64(args.Int(1)), args.Error(2)
}

// Test Authenticate with correct credentials
func TestAuthenticate_SuccessfulLogin(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := NewUserUsecase(mockRepo)

	userID := uuid.New()
	now := time.Now()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	storedUser := &domain.User{
		ID:        userID,
		Email:     "test@example.com",
		Password:  string(hashedPassword),
		Name:      "Test User",
		Phone:     "0123456789",
		Role:      domain.RoleCustomer,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(storedUser, nil)
	mockRepo.On("Update", mock.Anything, userID, mock.Anything).Return(storedUser, nil)

	ctx := context.Background()
	user, token, err := usecase.Authenticate(ctx, "test@example.com", "password123")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, domain.RoleCustomer, user.Role)

	mockRepo.AssertExpectations(t)
}

// Test Authenticate with wrong password
func TestAuthenticate_WrongPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := NewUserUsecase(mockRepo)

	userID := uuid.New()
	now := time.Now()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctPassword"), bcrypt.DefaultCost)

	storedUser := &domain.User{
		ID:        userID,
		Email:     "test@example.com",
		Password:  string(hashedPassword),
		Name:      "Test User",
		Phone:     "0123456789",
		Role:      domain.RoleCustomer,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(storedUser, nil)

	ctx := context.Background()
	user, token, err := usecase.Authenticate(ctx, "test@example.com", "wrongPassword")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)

	mockRepo.AssertExpectations(t)
}

// Test Authenticate with user not found
func TestAuthenticate_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := NewUserUsecase(mockRepo)

	mockRepo.On("GetByEmail", mock.Anything, "nonexistent@example.com").Return(nil, assert.AnError)

	ctx := context.Background()
	user, token, err := usecase.Authenticate(ctx, "nonexistent@example.com", "password123")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)

	mockRepo.AssertExpectations(t)
}

// Test Create user with valid data
func TestCreate_SuccessfulCreation(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := NewUserUsecase(mockRepo)

	newUser := &domain.User{
		Email:    "newuser@example.com",
		Password: "password123",
		Name:     "New User",
		Phone:    "0123456789",
	}

	userID := uuid.New()
	now := time.Now()
	createdUser := &domain.User{
		ID:        userID,
		Email:     "newuser@example.com",
		Password:  "hashed_password",
		Name:      "New User",
		Phone:     "0123456789",
		Role:      domain.RoleCustomer,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.Email == "newuser@example.com" &&
			u.Name == "New User" &&
			u.Role == domain.RoleCustomer &&
			u.IsActive == true
	})).Return(createdUser, nil)

	ctx := context.Background()
	result, err := usecase.Create(ctx, newUser)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "newuser@example.com", result.Email)
	assert.Equal(t, domain.RoleCustomer, result.Role)
	assert.True(t, result.IsActive)

	mockRepo.AssertExpectations(t)
}

// Test Create with missing email
func TestCreate_MissingEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := NewUserUsecase(mockRepo)

	newUser := &domain.User{
		Email:    "",
		Password: "password123",
		Name:     "Test User",
		Phone:    "0123456789",
	}

	ctx := context.Background()
	result, err := usecase.Create(ctx, newUser)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "email")

	mockRepo.AssertNotCalled(t, "Create")
}

// Test List users with pagination
func TestList_SuccessfulList(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := NewUserUsecase(mockRepo)

	users := []*domain.User{
		{
			ID:    uuid.New(),
			Email: "user1@example.com",
			Name:  "User 1",
			Role:  domain.RoleCustomer,
		},
		{
			ID:    uuid.New(),
			Email: "user2@example.com",
			Name:  "User 2",
			Role:  domain.RoleStylist,
		},
	}

	mockRepo.On("List", mock.Anything, 1, 10).Return(users, int64(2), nil)

	ctx := context.Background()
	result, err := usecase.List(ctx, 1, 10)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Items)
	assert.NotNil(t, result.Meta)
	assert.Equal(t, int64(2), result.Meta.Total)

	mockRepo.AssertExpectations(t)
}

// Test ListByRole with filter
func TestListByRole_SuccessfulList(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := NewUserUsecase(mockRepo)

	stylistUsers := []*domain.User{
		{
			ID:    uuid.New(),
			Email: "stylist1@example.com",
			Name:  "Stylist 1",
			Role:  domain.RoleStylist,
		},
		{
			ID:    uuid.New(),
			Email: "stylist2@example.com",
			Name:  "Stylist 2",
			Role:  domain.RoleStylist,
		},
	}

	mockRepo.On("ListByRole", mock.Anything, domain.RoleStylist, 1, 10).Return(stylistUsers, int64(2), nil)

	ctx := context.Background()
	result, err := usecase.ListByRole(ctx, domain.RoleStylist, 1, 10)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Items)
	assert.Equal(t, int64(2), result.Meta.Total)

	mockRepo.AssertExpectations(t)
}

// Test GetByID
func TestGetByID_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := NewUserUsecase(mockRepo)

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

	mockRepo.On("GetByID", mock.Anything, userID).Return(expectedUser, nil)

	ctx := context.Background()
	result, err := usecase.GetByID(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test@example.com", result.Email)

	mockRepo.AssertExpectations(t)
}

// Test Delete
func TestDelete_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := NewUserUsecase(mockRepo)

	userID := uuid.New()
	mockRepo.On("Delete", mock.Anything, userID).Return(nil)

	ctx := context.Background()
	err := usecase.Delete(ctx, userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
