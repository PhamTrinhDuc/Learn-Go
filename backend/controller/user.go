package controller

import (
	"backend/domain"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UserController struct {
	useCase  domain.UserUsecase
	validate *validator.Validate
}

func NewUserController(useCase domain.UserUsecase) *UserController {
	return &UserController{
		useCase:  useCase,
		validate: validator.New(),
	}
}

// CreateUser handles POST /users - Register a new user
func (uc *UserController) CreateUser(ctx *gin.Context) {
	var user domain.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	if err := uc.validate.Struct(&user); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	result, err := uc.useCase.Create(ctx.Request.Context(), &user)
	if err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	utils.RespondCreated(ctx, result)
}

// GetUser handles GET /users/:id - Get user details
func (uc *UserController) GetUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondBadRequest(ctx, "invalid user id format")
		return
	}

	user, err := uc.useCase.GetByID(ctx.Request.Context(), id)
	if err != nil {
		utils.RespondNotFound(ctx, err.Error())
		return
	}

	utils.RespondOK(ctx, user)
}

// UpdateUser handles PUT /users/:id - Update user information
func (uc *UserController) UpdateUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondBadRequest(ctx, "invalid user id format")
		return
	}

	// Get current user from context
	currentUserIDVal, exists := ctx.Get("user_id")
	currentRole, _ := ctx.Get("role")

	// User can only update themselves unless they're admin
	isAdmin := exists && (currentRole == string(domain.RoleAdmin))
	if exists && !isAdmin {
		currentUserID, _ := currentUserIDVal.(string)
		if currentUserID != idStr {
			utils.RespondBadRequest(ctx, "you can only update your own profile")
			return
		}
	}

	var user domain.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	result, err := uc.useCase.Update(ctx.Request.Context(), id, &user)
	if err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	utils.RespondOK(ctx, result)
}

// DeleteUser handles DELETE /users/:id - Delete a user
func (uc *UserController) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondBadRequest(ctx, "invalid user id format")
		return
	}

	if err := uc.useCase.Delete(ctx.Request.Context(), id); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	utils.RespondNoContent(ctx)
}

// ListUsers handles GET /users - List all users with pagination
func (uc *UserController) ListUsers(ctx *gin.Context) {
	page := utils.GetQueryInt(ctx, "page", 1)
	limit := utils.GetQueryInt(ctx, "limit", 10)

	result, err := uc.useCase.List(ctx.Request.Context(), page, limit)
	if err != nil {
		utils.RespondError(ctx, err)
		return
	}

	utils.RespondOK(ctx, result)
}

// ListUsersByRole handles GET /users/role/:role - List users by role
func (uc *UserController) ListUsersByRole(ctx *gin.Context) {
	roleStr := ctx.Param("role")
	role := domain.UserRole(roleStr)

	// Validate role
	validRoles := map[domain.UserRole]bool{
		domain.RoleAdmin:    true,
		domain.RoleStylist:  true,
		domain.RoleCustomer: true,
	}
	if !validRoles[role] {
		utils.RespondBadRequest(ctx, "invalid role")
		return
	}

	page := utils.GetQueryInt(ctx, "page", 1)
	limit := utils.GetQueryInt(ctx, "limit", 10)

	result, err := uc.useCase.ListByRole(ctx.Request.Context(), role, page, limit)
	if err != nil {
		utils.RespondError(ctx, err)
		return
	}

	utils.RespondOK(ctx, result)
}

// Login handles POST /auth/login - User authentication
func (uc *UserController) Login(ctx *gin.Context) {
	var req domain.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	if err := uc.validate.Struct(&req); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	user, token, err := uc.useCase.Authenticate(ctx.Request.Context(), req.Email, req.Password)
	if err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	utils.RespondOK(ctx, domain.LoginResponse{
		User:  user,
		Token: token,
	})
}

// GetMe handles GET /users/me - Get current user info
func (uc *UserController) GetMe(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		utils.RespondBadRequest(ctx, "user_id not found in context")
		return
	}

	userID := userIDVal.(string)
	id, err := uuid.Parse(userID)
	if err != nil {
		utils.RespondBadRequest(ctx, "invalid user_id format")
		return
	}

	user, err := uc.useCase.GetByID(ctx.Request.Context(), id)
	if err != nil {
		utils.RespondNotFound(ctx, err.Error())
		return
	}

	utils.RespondOK(ctx, user)
}

// Register handles POST /auth/register - User registration
func (uc *UserController) Register(ctx *gin.Context) {
	var req domain.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	if err := uc.validate.Struct(&req); err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	user := &domain.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Phone:    req.Phone,
		Role:     domain.RoleCustomer,
	}

	result, err := uc.useCase.Create(ctx.Request.Context(), user)
	if err != nil {
		utils.RespondBadRequest(ctx, err.Error())
		return
	}

	utils.RespondCreated(ctx, result)
}
