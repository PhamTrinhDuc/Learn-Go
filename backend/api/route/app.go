package route

import (
	"backend/api/middleware"
	"backend/controller"
	"backend/domain"

	"github.com/gin-gonic/gin"
)

// SetupProductRoutes sets up all product routes
func SetupProductRoutes(router *gin.RouterGroup, pc *controller.ProductController, authMiddleware *middleware.AuthMiddleware) {
	product := router.Group("/products")
	{
		// Read operations: authenticated users can view
		product.GET("", authMiddleware.OptionalHandler(), pc.List)
		product.GET("/:id", authMiddleware.OptionalHandler(), pc.GetByID)

		// Write operations: admin/manager only
		product.POST("", authMiddleware.Handler(), authMiddleware.RequireRole(string(domain.RoleAdmin), string(domain.RoleManager)), pc.Create)
		product.PUT("/:id", authMiddleware.Handler(), authMiddleware.RequireRole(string(domain.RoleAdmin), string(domain.RoleManager)), pc.Update)
		product.DELETE("/:id", authMiddleware.Handler(), authMiddleware.RequireRole(string(domain.RoleAdmin), string(domain.RoleManager)), pc.Delete)
	}
}

// SetupBranchRoutes sets up all branch routes
func SetupBranchRoutes(router *gin.RouterGroup, bc *controller.BranchController, authMiddleware *middleware.AuthMiddleware) {
	branch := router.Group("/branches")
	{
		// Read operations: authenticated users can view
		branch.GET("", authMiddleware.OptionalHandler(), bc.List)        // List all branches with pagination
		branch.GET("/:id", authMiddleware.OptionalHandler(), bc.GetByID) // Get a specific branch

		// Write operations: admin/manager only
		branch.POST("", authMiddleware.Handler(), authMiddleware.RequireRole(string(domain.RoleAdmin), string(domain.RoleManager)), bc.Create)       // Create a new branch
		branch.PUT("/:id", authMiddleware.Handler(), authMiddleware.RequireRole(string(domain.RoleAdmin), string(domain.RoleManager)), bc.Update)    // Update a branch
		branch.DELETE("/:id", authMiddleware.Handler(), authMiddleware.RequireRole(string(domain.RoleAdmin), string(domain.RoleManager)), bc.Delete) // Delete a branch
	}
}

// SetupServiceRoutes sets up all service routes
func SetupServiceRoutes(router *gin.RouterGroup, sc *controller.ServiceController, authMiddleware *middleware.AuthMiddleware) {
	service := router.Group("/services")
	{
		// Read operations: authenticated users can view
		service.GET("", authMiddleware.OptionalHandler(), sc.List)                                      // List all services with pagination
		service.GET("/:id", authMiddleware.OptionalHandler(), sc.GetByID)                               // Get a specific service
		service.GET("/category/:category", authMiddleware.OptionalHandler(), sc.ListServicesByCategory) // List services by category

		// Write operations: admin/manager only
		service.POST("", authMiddleware.Handler(), authMiddleware.RequireRole(string(domain.RoleAdmin), string(domain.RoleManager)), sc.Create)       // Create a new service
		service.PUT("/:id", authMiddleware.Handler(), authMiddleware.RequireRole(string(domain.RoleAdmin), string(domain.RoleManager)), sc.Update)    // Update a service
		service.DELETE("/:id", authMiddleware.Handler(), authMiddleware.RequireRole(string(domain.RoleAdmin), string(domain.RoleManager)), sc.Delete) // Delete a service
	}
}

// SetupStylistRoutes sets up all stylist routes
func SetupStylistRoutes(router *gin.RouterGroup, stc *controller.StylistController, authMiddleware *middleware.AuthMiddleware) {
	stylist := router.Group("/stylists")
	{
		// Read operations: authenticated users can view
		stylist.GET("", authMiddleware.OptionalHandler(), stc.List)                                   // List all stylists with pagination
		stylist.GET("/:id", authMiddleware.OptionalHandler(), stc.GetByID)                            // Get a specific stylist
		stylist.GET("/branch/:branch_id", authMiddleware.OptionalHandler(), stc.ListStylistsByBranch) // List stylists by branch

		// Write operations: admin/manager only
		stylist.POST("", authMiddleware.Handler(), authMiddleware.RequireRole(string(domain.RoleAdmin), string(domain.RoleManager)), stc.Create)       // Create a new stylist
		stylist.PUT("/:id", authMiddleware.Handler(), authMiddleware.RequireRole(string(domain.RoleAdmin), string(domain.RoleManager)), stc.Update)    // Update a stylist
		stylist.DELETE("/:id", authMiddleware.Handler(), authMiddleware.RequireRole(string(domain.RoleAdmin), string(domain.RoleManager)), stc.Delete) // Delete a stylist
	}
}

// SetupUserRoutes sets up all user routes
func SetupUserRoutes(router *gin.RouterGroup, uc *controller.UserController, authMiddleware *middleware.AuthMiddleware) {
	// Public auth routes (no authentication required)
	auth := router.Group("/auth")
	{
		auth.POST("/login", uc.Login)       // User login
		auth.POST("/register", uc.Register) // User registration
	}

	// Protected user routes (authentication required)
	user := router.Group("/users")
	user.Use(authMiddleware.Handler()) // Require authentication
	{
		// User can access their own info
		user.GET("/me", uc.GetMe) // Get current user info

		// Admin/Manager only - can manage all users
		user.GET("", authMiddleware.RequireRole(string(domain.RoleAdmin), string(domain.RoleManager)), uc.ListUsers)
		user.POST("", authMiddleware.RequireRole(string(domain.RoleAdmin)), uc.CreateUser)
		user.GET("/:id", uc.GetUser)                                                             // Can view any user (for customer service)
		user.PUT("/:id", uc.UpdateUser)                                                          // Can update own or admin can update any
		user.DELETE("/:id", authMiddleware.RequireRole(string(domain.RoleAdmin)), uc.DeleteUser) // Admin only

		// Admin/Manager only - view by role
		user.GET("/role/:role", authMiddleware.RequireRole(string(domain.RoleAdmin), string(domain.RoleManager)), uc.ListUsersByRole)
	}
}
