package route

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"tmdt-backend/controller"
	"tmdt-backend/repository"
	"tmdt-backend/services"
	"tmdt-backend/usecase"
)

// CORSMiddleware handles cross-origin requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientURL := os.Getenv("CLIENT_URL")
		origin := c.Request.Header.Get("Origin")

		if clientURL == "" || clientURL == "*" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else if origin == clientURL {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// Setup registers all API routes
func Setup(db *pgxpool.Pool, rClient *redis.Client, r *gin.Engine) {
	// Apply global middleware
	r.Use(CORSMiddleware())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	mailerService := services.NewMailerService()

	// Initialize usecases
	userUsecase := usecase.NewUserUsecase(userRepo, rClient)
	otpUsecase := usecase.NewOTPUsecase(rClient, mailerService)

	// Initialize controllers
	userCtrl := controller.NewUserController(userUsecase)
	otpCtrl := controller.NewOTPController(otpUsecase)

	// Setup API group
	api := r.Group("/api")
	{
		// User endpoints
		userGroup := api.Group("/user")
		{
			userGroup.GET("/", userCtrl.ListUsers)
			userGroup.POST("/register", userCtrl.Register)
			userGroup.POST("/login", userCtrl.Login)
			userGroup.POST("/google/login", userCtrl.GoogleLogin)
			userGroup.GET("/information/:id", userCtrl.GetProfile)
			userGroup.PUT("/information/:id", userCtrl.UpdateProfile)
			userGroup.GET("/role", userCtrl.GetRoles)
			userGroup.PUT("/role/:id", userCtrl.UpdateRole)
			userGroup.PUT("/status/:id", userCtrl.UpdateStatus)
			userGroup.POST("/check-password-validity", userCtrl.CheckPasswordValidity)
			userGroup.POST("/reset-password", userCtrl.ResetPassword)
			userGroup.POST("/check-email", userCtrl.CheckEmail)
			userGroup.PUT("/admin-reset-password/:id", userCtrl.AdminResetPassword)
		}

		// OTP endpoints
		otpGroup := api.Group("/otp")
		{
			otpGroup.POST("/send-otp", otpCtrl.SendOTP)
			otpGroup.POST("/verify-otp", otpCtrl.VerifyOTP)
		}
	}
}
