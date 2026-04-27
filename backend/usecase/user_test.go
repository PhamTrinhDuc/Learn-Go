package usecase

import (
	"context"
	"testing"

	"backend/domain"
	mock_domain "backend/mocks/backend/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserUsecase_Authenticate_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mock_domain.MockUserRepository)
	usecase := NewUserUsecase(mockRepo)

	email := "test@example.com"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	userID := uuid.New()

	user := &domain.User{
		ID:       userID,
		Email:    email,
		Password: string(hashedPassword),
		Name:     "Test User",
		Role:     domain.RoleCustomer,
	}

	mockRepo.On("GetByEmail", mock.Anything, email).Return(user, nil)
	mockRepo.On("Update", mock.Anything, userID, mock.Anything).Return(user, nil)

	// Act
	authenticatedUser, token, err := usecase.Authenticate(context.Background(), email, password)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, authenticatedUser)
	assert.NotEmpty(t, token)
	assert.Equal(t, user.ID, authenticatedUser.ID)
	assert.Equal(t, email, authenticatedUser.Email)
	mockRepo.AssertCalled(t, "GetByEmail", mock.Anything, email)
	mockRepo.AssertCalled(t, "Update", mock.Anything, userID, mock.Anything)
}

func TestUserUsecase_Authenticate_WrongPassword(t *testing.T) {
	// Arrange
	mockRepo := new(mock_domain.MockUserRepository)
	usecase := NewUserUsecase(mockRepo)

	email := "test@example.com"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	userID := uuid.New()

	user := &domain.User{
		ID:       userID,
		Email:    email,
		Password: string(hashedPassword),
	}

	mockRepo.On("GetByEmail", mock.Anything, email).Return(user, nil)

	// Act
	_, _, err := usecase.Authenticate(context.Background(), email, "wrongpassword")

	// Assert
	assert.Error(t, err)
}

func TestUserUsecase_Create_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mock_domain.MockUserRepository)
	usecase := NewUserUsecase(mockRepo)

	req := &domain.User{
		Email:    "newuser@example.com",
		Password: "password123",
		Name:     "New User",
		Phone:    "0123456789",
	}

	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.Email == req.Email && u.Name == req.Name && u.Role == domain.RoleCustomer
	})).Return(&domain.User{
		ID:    uuid.New(),
		Email: req.Email,
		Name:  req.Name,
		Role:  domain.RoleCustomer,
	}, nil)

	// Act
	result, err := usecase.Create(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, domain.RoleCustomer, result.Role)
	mockRepo.AssertCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestUserUsecase_Create_MissingEmail(t *testing.T) {
	// Arrange
	mockRepo := new(mock_domain.MockUserRepository)
	usecase := NewUserUsecase(mockRepo)

	req := &domain.User{
		Password: "password123",
		Name:     "New User",
	}

	// Act
	_, err := usecase.Create(context.Background(), req)

	// Assert
	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestUserUsecase_GetByID_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mock_domain.MockUserRepository)
	usecase := NewUserUsecase(mockRepo)

	userID := uuid.New()
	user := &domain.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}

	mockRepo.On("GetByID", mock.Anything, userID).Return(user, nil)

	// Act
	result, err := usecase.GetByID(context.Background(), userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, user, result)
	mockRepo.AssertCalled(t, "GetByID", mock.Anything, userID)
}

func TestUserUsecase_List_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mock_domain.MockUserRepository)
	usecase := NewUserUsecase(mockRepo)

	users := []*domain.User{
		{ID: uuid.New(), Email: "user1@example.com", Name: "User 1"},
		{ID: uuid.New(), Email: "user2@example.com", Name: "User 2"},
	}

	mockRepo.On("List", mock.Anything, 1, 10).Return(users, int64(2), nil)

	// Act
	result, err := usecase.List(context.Background(), 1, 10)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result.Items, 2)
	mockRepo.AssertCalled(t, "List", mock.Anything, 1, 10)
}

func TestUserUsecase_Delete_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mock_domain.MockUserRepository)
	usecase := NewUserUsecase(mockRepo)

	userID := uuid.New()
	mockRepo.On("Delete", mock.Anything, userID).Return(nil)

	// Act
	err := usecase.Delete(context.Background(), userID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "Delete", mock.Anything, userID)
}

func TestUserUsecase_Update_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mock_domain.MockUserRepository)
	usecase := NewUserUsecase(mockRepo)

	userID := uuid.New()
	existingUser := &domain.User{
		ID:    userID,
		Name:  "Old Name",
		Email: "test@example.com",
	}

	updateReq := &domain.User{
		Name: "Updated Name",
	}

	mockRepo.On("GetByID", mock.Anything, userID).Return(existingUser, nil)
	mockRepo.On("Update", mock.Anything, userID, mock.MatchedBy(func(u *domain.User) bool {
		return u.Name == "Updated Name"
	})).Return(&domain.User{
		ID:    userID,
		Name:  "Updated Name",
		Email: "test@example.com",
	}, nil)

	// Act
	result, err := usecase.Update(context.Background(), userID, updateReq)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated Name", result.Name)
	mockRepo.AssertCalled(t, "GetByID", mock.Anything, userID)
	mockRepo.AssertCalled(t, "Update", mock.Anything, userID, mock.Anything)
}
