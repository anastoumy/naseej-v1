package graph

import (
	"github.com/naseej/server/config"
	"github.com/naseej/server/internal/repository"
	"github.com/naseej/server/internal/service"
	"github.com/naseej/server/internal/websocket"
	"gorm.io/gorm"
)

// Resolver is the root resolver dependency injection container.
type Resolver struct {
	DB               *gorm.DB
	Config           *config.Config
	WSHub            *websocket.Hub
	AuthService      *service.AuthService
	ProductService   *service.ProductService
	InventoryService *service.InventoryService
	OrderService     *service.OrderService
	CustomerService  *service.CustomerService
	PaymentService   *service.PaymentService
	StaffService     *service.StaffService
	MerchantRepo     *repository.MerchantRepository
}

// NewResolver creates a fully wired Resolver with all service dependencies.
func NewResolver(db *gorm.DB, cfg *config.Config, wsHub *websocket.Hub) *Resolver {
	// Repositories
	merchantRepo := repository.NewMerchantRepository(db)
	productRepo := repository.NewProductRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	staffRepo := repository.NewStaffRepository(db)

	// Services
	authService := service.NewAuthService(staffRepo, cfg)
	productService := service.NewProductService(productRepo)
	inventoryService := service.NewInventoryService(inventoryRepo, wsHub)
	orderService := service.NewOrderService(orderRepo, inventoryRepo, wsHub)
	customerService := service.NewCustomerService(customerRepo)
	paymentService := service.NewPaymentService(paymentRepo, orderRepo)
	staffService := service.NewStaffService(staffRepo)

	return &Resolver{
		DB:               db,
		Config:           cfg,
		WSHub:            wsHub,
		AuthService:      authService,
		ProductService:   productService,
		InventoryService: inventoryService,
		OrderService:     orderService,
		CustomerService:  customerService,
		PaymentService:   paymentService,
		StaffService:     staffService,
		MerchantRepo:     merchantRepo,
	}
}
