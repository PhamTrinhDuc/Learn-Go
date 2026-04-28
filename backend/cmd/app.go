package main

import (
	"backend/api/middleware"
	"backend/api/route"
	"backend/bootstrap"
	"backend/internal/utils"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	// Khởi tạo Gin engine
	r := gin.Default()

	// Tạo context lắng nghe tín hiệu từ hệ điều hành
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := bootstrap.NewDB(ctx, bootstrap.DBConfig{
		Host:     utils.GetEnvString("DB_HOST", "localhost"),
		Port:     utils.GetEnvInt("DB_PORT", 5433),
		User:     utils.GetEnvString("DB_USER", "mcp_user"), // Use app_user for RLS enforcement
		Password: utils.GetEnvString("DB_PASSWORD", "mcp_password"),
		DBName:   utils.GetEnvString("DB_NAME", "salon_chain"),
		SSLMode:  utils.GetEnvString("DB_SSLMODE", "disable"),
		MaxConns: int32(utils.GetEnvInt("MAX_CONNS", 10)),
		MinConns: int32(utils.GetEnvInt("MAX_CONNS", 2)),
	})
	if err != nil {
		log.Fatalf("Không thể kết nối DB: %v", err)
	}

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware() // Khởi tạo middleware của bạn

	//  Group các API lại (ví dụ prefix /api/v1)
	container := bootstrap.NewContainer(db.GetPool())

	v1 := r.Group("/api/v1")
	route.SetupBranchRoutes(v1, container.BranchCtl, authMiddleware)
	route.SetupUserRoutes(v1, container.UserCtl, authMiddleware)
	route.SetupStylistRoutes(v1, container.StylistCtl, authMiddleware)

	// 6. Chạy Server
	r.Run(":8080") // Server sẽ chạy tại http://localhost:8080
}
