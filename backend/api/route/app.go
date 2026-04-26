package route

import (
	"backend/controller"

	"github.com/gin-gonic/gin"
)

// SetupProductRoutes sets up all product routes
func SetupProductRoutes(router *gin.RouterGroup, pc *controller.ProductController) {
	product := router.Group("/products")
	{
		product.GET("", pc.List)
		product.GET("/:id", pc.GetByID)
		product.POST("", pc.Create)
		product.PUT("/:id", pc.Update)
		product.DELETE("/:id", pc.Delete)
	}
}

// SetupBranchRoutes sets up all branch routes
func SetupBranchRoutes(router *gin.RouterGroup, bc *controller.BranchController) {
	branch := router.Group("/branches")
	{
		branch.POST("", bc.Create)       // Create a new branch
		branch.GET("", bc.List)          // List all branches with pagination
		branch.GET("/:id", bc.GetByID)   // Get a specific branch
		branch.PUT("/:id", bc.Update)    // Update a branch
		branch.DELETE("/:id", bc.Delete) // Delete a branch
	}
}

// SetupServiceRoutes sets up all service routes
func SetupServiceRoutes(router *gin.RouterGroup, sc *controller.ServiceController) {
	service := router.Group("/services")
	{
		service.POST("", sc.Create)                                   // Create a new service
		service.GET("", sc.List)                                      // List all services with pagination
		service.GET("/:id", sc.GetByID)                               // Get a specific service
		service.PUT("/:id", sc.Update)                                // Update a service
		service.DELETE("/:id", sc.Delete)                             // Delete a service
		service.GET("/category/:category", sc.ListServicesByCategory) // List services by category
	}
}

// SetupStylistRoutes sets up all stylist routes
func SetupStylistRoutes(router *gin.RouterGroup, stc *controller.StylistController) {
	stylist := router.Group("/stylists")
	{
		stylist.POST("", stc.Create)                                // Create a new stylist
		stylist.GET("", stc.List)                                   // List all stylists with pagination
		stylist.GET("/:id", stc.GetByID)                            // Get a specific stylist
		stylist.PUT("/:id", stc.Update)                             // Update a stylist
		stylist.DELETE("/:id", stc.Delete)                          // Delete a stylist
		stylist.GET("/branch/:branch_id", stc.ListStylistsByBranch) // List stylists by branch
	}
}
