package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/domain"
	mock_domain "backend/mocks/backend/domain"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestGinContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(w)
	return ctx, w
}

func TestUserController_CreateUser_Success(t *testing.T) {
	// Arrange
	mockUC := new(mock_domain.MockUserUsecase)
	controller := &UserController{
		useCase:  mockUC,
		validate: validator.New(),
	}

	ctx, w := setupTestGinContext()
	reqBody := gin.H{
		"email":    "test@example.com",
		"password": "password123",
		"name":     "Test User",
		"phone":    "0123456789",
	}

	bodyBytes, _ := json.Marshal(reqBody)
	ctx.Request = httptest.NewRequest("POST", "/users", bytes.NewReader(bodyBytes))
	ctx.Request.Header.Set("Content-Type", "application/json")

	mockUC.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.Email == "test@example.com" && u.Name == "Test User"
	})).Return(&domain.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
		Role:  domain.RoleCustomer,
	}, nil)

	// Act
	controller.CreateUser(ctx)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	mockUC.AssertCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestUserController_Login_Success(t *testing.T) {
	// Arrange
	mockUC := new(mock_domain.MockUserUsecase)
	controller := &UserController{
		useCase:  mockUC,
		validate: validator.New(),
	}

	ctx, w := setupTestGinContext()
	reqBody := gin.H{
		"email":    "test@example.com",
		"password": "password123",
	}

	bodyBytes, _ := json.Marshal(reqBody)
	ctx.Request = httptest.NewRequest("POST", "/auth/login", bytes.NewReader(bodyBytes))
	ctx.Request.Header.Set("Content-Type", "application/json")

	user := &domain.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
		Role:  domain.RoleCustomer,
	}

	mockUC.On("Authenticate", mock.Anything, "test@example.com", "password123").
		Return(user, "fake-token", nil)

	// Act
	controller.Login(ctx)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertCalled(t, "Authenticate", mock.Anything, "test@example.com", "password123")
}

func TestUserController_GetMe_Success(t *testing.T) {
	// Arrange
	mockUC := new(mock_domain.MockUserUsecase)
	controller := &UserController{
		useCase:  mockUC,
		validate: validator.New(),
	}

	userID := uuid.New()
	ctx, w := setupTestGinContext()
	ctx.Request = httptest.NewRequest("GET", "/users/me", nil)
	ctx.Set("user_id", userID.String())

	user := &domain.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}

	mockUC.On("GetByID", mock.Anything, userID).Return(user, nil)

	// Act
	controller.GetMe(ctx)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertCalled(t, "GetByID", mock.Anything, userID)
}

func TestUserController_Register_Success(t *testing.T) {
	// Arrange
	mockUC := new(mock_domain.MockUserUsecase)
	controller := &UserController{
		useCase:  mockUC,
		validate: validator.New(),
	}

	ctx, w := setupTestGinContext()
	reqBody := gin.H{
		"email":    "newuser@example.com",
		"password": "password123",
		"name":     "New User",
		"phone":    "0123456789",
	}

	bodyBytes, _ := json.Marshal(reqBody)
	ctx.Request = httptest.NewRequest("POST", "/auth/register", bytes.NewReader(bodyBytes))
	ctx.Request.Header.Set("Content-Type", "application/json")

	newUser := &domain.User{
		ID:    uuid.New(),
		Email: "newuser@example.com",
		Name:  "New User",
		Role:  domain.RoleCustomer,
	}

	mockUC.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.Email == "newuser@example.com"
	})).Return(newUser, nil)

	// Act
	controller.Register(ctx)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	mockUC.AssertCalled(t, "Create", mock.Anything, mock.Anything)
}
