package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/domain"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserUsecase manually implements domain.UserUsecase
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

func setupTestRouter(uc domain.UserUsecase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	ctrl := NewUserController(uc)

	r.POST("/users", ctrl.CreateUser)
	r.GET("/users/:id", ctrl.GetUser)
	r.PUT("/users/:id", func(c *gin.Context) {
		// Mock auth middleware injecting values
		if c.GetHeader("X-Mock-User-ID") != "" {
			c.Set("user_id", c.GetHeader("X-Mock-User-ID"))
			c.Set("role", c.GetHeader("X-Mock-Role"))
		}
		ctrl.UpdateUser(c)
	})
	r.DELETE("/users/:id", ctrl.DeleteUser)
	r.GET("/users", ctrl.ListUsers)
	r.GET("/users/role/:role", ctrl.ListUsersByRole)
	r.POST("/auth/login", ctrl.Login)
	r.GET("/users/me", func(c *gin.Context) {
		if c.GetHeader("X-Mock-User-ID") != "" {
			c.Set("user_id", c.GetHeader("X-Mock-User-ID"))
		}
		ctrl.GetMe(c)
	})
	r.POST("/auth/register", ctrl.Register)

	return r
}

func TestUserController_CreateUser(t *testing.T) {
	mockUC := new(MockUserUsecase)
	router := setupTestRouter(mockUC)

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		user := &domain.User{
			Email: "test@example.com",
			Name:  "Test User",
			Phone: "1234567890",
			Role:  domain.RoleCustomer,
		}

		mockUC.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
			return u.Email == user.Email
		})).Return(&domain.User{
			ID:    userID,
			Email: user.Email,
			Name:  user.Name,
			Phone: user.Phone,
			Role:  user.Role,
		}, nil).Once()

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("Bind_Error", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Validate_Error", func(t *testing.T) {
		// Missing fields that validator requires
		user := &domain.User{
			Email: "invalid-email",
		}
		body, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Usecase_Error", func(t *testing.T) {
		user := &domain.User{
			Email: "test@example.com",
			Name:  "Test User",
			Phone: "1234567890",
			Role:  domain.RoleCustomer,
		}

		mockUC.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("database error")).Once()

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		mockUC.AssertExpectations(t)
	})
}

func TestUserController_GetUser(t *testing.T) {
	mockUC := new(MockUserUsecase)
	router := setupTestRouter(mockUC)

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		user := &domain.User{
			ID:    userID,
			Email: "test@example.com",
			Name:  "Test User",
		}

		mockUC.On("GetByID", mock.Anything, userID).Return(user, nil).Once()

		req, _ := http.NewRequest("GET", "/users/"+userID.String(), nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("Invalid_ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/invalid-uuid", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		userID := uuid.New()
		mockUC.On("GetByID", mock.Anything, userID).Return(nil, errors.New("user not found")).Once()

		req, _ := http.NewRequest("GET", "/users/"+userID.String(), nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockUC.AssertExpectations(t)
	})
}

func TestUserController_UpdateUser(t *testing.T) {
	mockUC := new(MockUserUsecase)
	router := setupTestRouter(mockUC)

	t.Run("Success_Self", func(t *testing.T) {
		userID := uuid.New()
		user := &domain.User{
			Name:  "Updated Name",
			Phone: "0987654321",
		}

		mockUC.On("Update", mock.Anything, userID, mock.Anything).Return(&domain.User{
			ID:    userID,
			Email: "test@example.com",
			Name:  user.Name,
			Phone: user.Phone,
		}, nil).Once()

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest("PUT", "/users/"+userID.String(), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-User-ID", userID.String())
		req.Header.Set("X-Mock-Role", string(domain.RoleCustomer))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("Success_Admin_Updating_Other", func(t *testing.T) {
		userID := uuid.New()
		adminID := uuid.New()
		user := &domain.User{
			Name:  "Updated Name",
			Phone: "0987654321",
		}

		mockUC.On("Update", mock.Anything, userID, mock.Anything).Return(&domain.User{
			ID:    userID,
			Email: "test@example.com",
			Name:  user.Name,
			Phone: user.Phone,
		}, nil).Once()

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest("PUT", "/users/"+userID.String(), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-User-ID", adminID.String())
		req.Header.Set("X-Mock-Role", string(domain.RoleAdmin))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("Unauthorized_Update_Other", func(t *testing.T) {
		userID := uuid.New()
		otherUserID := uuid.New()
		user := &domain.User{
			Name: "Updated Name",
		}

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest("PUT", "/users/"+userID.String(), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-User-ID", otherUserID.String())
		req.Header.Set("X-Mock-Role", string(domain.RoleCustomer))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

func TestUserController_DeleteUser(t *testing.T) {
	mockUC := new(MockUserUsecase)
	router := setupTestRouter(mockUC)

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		mockUC.On("Delete", mock.Anything, userID).Return(nil).Once()

		req, _ := http.NewRequest("DELETE", "/users/"+userID.String(), nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNoContent, resp.Code)
		mockUC.AssertExpectations(t)
	})
}

func TestUserController_ListUsers(t *testing.T) {
	mockUC := new(MockUserUsecase)
	router := setupTestRouter(mockUC)

	t.Run("Success", func(t *testing.T) {
		users := []*domain.User{
			{ID: uuid.New(), Email: "u1@example.com"},
			{ID: uuid.New(), Email: "u2@example.com"},
		}

		paginatedResponse := &utils.PaginatedResponse[domain.User]{
			Items: users,
			Meta: utils.PaginationMeta{
				Total:      2,
				Page:       1,
				Limit:      10,
				TotalPages: 1,
			},
		}

		mockUC.On("List", mock.Anything, 1, 10).Return(paginatedResponse, nil).Once()

		req, _ := http.NewRequest("GET", "/users?page=1&limit=10", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockUC.AssertExpectations(t)
	})
}

func TestUserController_ListUsersByRole(t *testing.T) {
	mockUC := new(MockUserUsecase)
	router := setupTestRouter(mockUC)

	t.Run("Success", func(t *testing.T) {
		users := []*domain.User{
			{ID: uuid.New(), Email: "u1@example.com", Role: domain.RoleCustomer},
		}

		paginatedResponse := &utils.PaginatedResponse[domain.User]{
			Items: users,
			Meta: utils.PaginationMeta{
				Total:      2,
				Page:       1,
				Limit:      10,
				TotalPages: 1,
			},
		}

		mockUC.On("ListByRole", mock.Anything, domain.RoleCustomer, 1, 10).Return(paginatedResponse, nil).Once()

		req, _ := http.NewRequest("GET", "/users/role/customer?page=1&limit=10", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("Invalid_Role", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/role/invalidrole", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

func TestUserController_Login(t *testing.T) {
	mockUC := new(MockUserUsecase)
	router := setupTestRouter(mockUC)

	t.Run("Success", func(t *testing.T) {
		reqBody := domain.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		expectedUser := &domain.User{
			ID:    uuid.New(),
			Email: "test@example.com",
			Role:  domain.RoleCustomer,
		}
		expectedToken := "test-jwt-token"

		mockUC.On("Authenticate", mock.Anything, reqBody.Email, reqBody.Password).Return(expectedUser, expectedToken, nil).Once()

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockUC.AssertExpectations(t)
	})
}

func TestUserController_GetMe(t *testing.T) {
	mockUC := new(MockUserUsecase)
	router := setupTestRouter(mockUC)

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		user := &domain.User{
			ID:    userID,
			Email: "me@example.com",
		}

		mockUC.On("GetByID", mock.Anything, userID).Return(user, nil).Once()

		req, _ := http.NewRequest("GET", "/users/me", nil)
		req.Header.Set("X-Mock-User-ID", userID.String())
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("Missing_Context_ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/me", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

func TestUserController_Register(t *testing.T) {
	mockUC := new(MockUserUsecase)
	router := setupTestRouter(mockUC)

	t.Run("Success", func(t *testing.T) {
		reqBody := domain.RegisterRequest{
			Email:    "new@example.com",
			Password: "password123",
			Name:     "New User",
			Phone:    "1234567890",
		}

		expectedUser := &domain.User{
			ID:    uuid.New(),
			Email: reqBody.Email,
			Name:  reqBody.Name,
			Phone: reqBody.Phone,
			Role:  domain.RoleCustomer,
		}

		mockUC.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
			return u.Email == reqBody.Email && u.Role == domain.RoleCustomer
		})).Return(expectedUser, nil).Once()

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		mockUC.AssertExpectations(t)
	})
}
