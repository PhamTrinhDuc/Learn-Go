// bootstrap/container.go
package bootstrap

import (
	"backend/controller"
	"backend/internal/observability"
	"backend/repository"
	"backend/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	// Repositories
	BranchRepo          *repository.BranchRepository
	UserRepo            *repository.UserRepository
	StylistRepo         *repository.StylistRepository
	StylistScheduleRepo *repository.StylistScheduleRepository
	ServicesRepo        *repository.ServiceRepository
	ProductRepo         *repository.ProductRepository
	OrderRepo           *repository.OrderRepository
	OrderItemRepo       *repository.OrderItemsRepository

	// Usecases
	BranchUC          *usecase.BranchUsecase
	UserUC            *usecase.UserUsecase
	StylistUC         *usecase.StylistUsecase
	StylistScheduleUC *usecase.StylistScheduleUsecase
	ServicesUC        *usecase.ServiceUsecase
	ProductUC         *usecase.ProductUseCase
	OrderUC           *usecase.OrderUsecase
	OrderItemUC       *usecase.OrderItemsUsecase

	// Controllers
	BranchCtl          *controller.BranchController
	UserCtl            *controller.UserController
	StylistCtl         *controller.StylistController
	StylistScheduleCtl *controller.StylistScheduleController
	ServicesCtl        *controller.ServiceController
	ProductCtl         *controller.ProductController
	OrderCtl           *controller.OrderController
	OrderItemCtl       *controller.OrderItemsController
}

func NewContainer(pool *pgxpool.Pool, telemetry *observability.Telemetry) *Container {
	// Repositories
	branchRepo := repository.NewBranchRepository(pool)
	userRepo := repository.NewUserRepository(pool)
	stylistRepo := repository.NewStylistRepository(pool)
	stylistScheduleRepo := repository.NewStylistScheduleRepository(pool)
	productRepo := repository.NewProductRepository(pool)
	servicesRepo := repository.NewServiceRepository(pool)
	orderRepo := repository.NewOrderRepository(pool)
	orderItemRepo := repository.NewOrderItemsRepository(pool)

	// Usecases
	branchUC := usecase.NewBranchUsecase(branchRepo)
	userUC := usecase.NewUserUsecase(userRepo)
	stylistUC := usecase.NewStylistUsecase(stylistRepo)
	stylistScheduleUC := usecase.NewStylistScheduleUsecase(stylistScheduleRepo)
	productUC := usecase.NewProductUseCase(productRepo)
	servicesUC := usecase.NewServiceUsecase(servicesRepo)
	orderUC := usecase.NewOrderUsecase(orderRepo)
	orderItemUC := usecase.NewOrderItemsUsecase(orderItemRepo)

	// Controllers
	branchCtl := controller.NewBranchController(branchUC, telemetry)
	userCtl := controller.NewUserController(userUC)
	stylistCtl := controller.NewStylistController(stylistUC, telemetry)
	stylistScheduleCtl := controller.NewStylistScheduleController(stylistScheduleUC)
	productCtl := controller.NewProductController(productUC)
	serviceCtl := controller.NewServiceController(servicesUC)
	orderCtl := controller.NewOrderController(orderUC)
	orderItemCtl := controller.NewOrderItemsController(orderItemUC)

	return &Container{
		BranchRepo:          branchRepo,
		UserRepo:            userRepo,
		StylistRepo:         stylistRepo,
		StylistScheduleRepo: stylistScheduleRepo,
		ProductRepo:         productRepo,
		ServicesRepo:        servicesRepo,
		OrderRepo:           orderRepo,
		OrderItemRepo:       orderItemRepo,

		BranchUC:          branchUC,
		UserUC:            userUC,
		StylistUC:         stylistUC,
		StylistScheduleUC: stylistScheduleUC,
		ProductUC:         productUC,
		ServicesUC:        servicesUC,
		OrderUC:           orderUC,
		OrderItemUC:       orderItemUC,

		BranchCtl:          branchCtl,
		UserCtl:            userCtl,
		StylistCtl:         stylistCtl,
		StylistScheduleCtl: stylistScheduleCtl,
		ProductCtl:         productCtl,
		ServicesCtl:        serviceCtl,
		OrderCtl:           orderCtl,
		OrderItemCtl:       orderItemCtl,
	}
}
