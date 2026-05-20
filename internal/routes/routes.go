package routes

import (
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/config"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/handler"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/middleware"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/repository"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB, redisClient *redis.Client, cfg *config.Config) {
	// Repositories
	userRepo := repository.NewUserRepo(db)
	categoryRepo := repository.NewCategoryRepository(db)
	menuRepo := repository.NewMenuRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	orderRepo := repository.NewOrderRepo(db)
	paymentRepo := repository.NewPaymentRepository(db)
	returnRepo := repository.NewReturnRepository(db)
	shiftRepo := repository.NewShiftRepository(db)
	cashMovementRepo := repository.NewCashMovementRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)
	stockRepo := repository.NewStockRepository(db)
	taxRepo := repository.NewTaxRepository(db)
	discountRepo := repository.NewDiscountRepository(db)
	reservationRepo := repository.NewReservationRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, cfg)
	employeeService := service.NewEmployeeService(db, userRepo)
	categoryService := service.NewCategoryService(db, categoryRepo)
	menuService := service.NewMenuService(db, menuRepo)
	customerService := service.NewCustomerService(db, customerRepo)
	orderService := service.NewOrderService(db, orderRepo, inventoryRepo, shiftRepo, taxRepo)
	paymentService := service.NewPaymentService(db, paymentRepo, orderRepo, orderService, cfg)
	returnService := service.NewReturnService(db, returnRepo, orderRepo, orderService)
	shiftService := service.NewShiftService(db, shiftRepo, cashMovementRepo)
	inventoryService := service.NewInventoryService(db, inventoryRepo)
	stockService := service.NewStockService(db, stockRepo, inventoryRepo)
	taxService := service.NewTaxService(db, taxRepo)
	discountService := service.NewDiscountService(db, discountRepo)
	reservationService := service.NewReservationService(db, reservationRepo)
	analyticsService := service.NewAnalyticsService(db, analyticsRepo)
	printerService := service.NewPrinterService()

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	employeeHandler := handler.NewEmployeeHandler(employeeService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	menuHandler := handler.NewMenuHandler(menuService)
	customerHandler := handler.NewCustomerHandler(customerService)
	orderHandler := handler.NewOrderHandler(orderService)
	paymentHandler := handler.NewPaymentHandler(paymentService)
	returnHandler := handler.NewReturnHandler(returnService)
	shiftHandler := handler.NewShiftHandler(shiftService)
	cashMovementHandler := handler.NewCashMovementHandler(shiftService)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)
	stockHandler := handler.NewStockHandler(stockService)
	taxHandler := handler.NewTaxHandler(taxService)
	discountHandler := handler.NewDiscountHandler(discountService)
	reservationHandler := handler.NewReservationHandler(reservationService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)
	printerHandler := handler.NewPrinterHandler(printerService)
	printerHandler.SetOrderService(orderService) // FIX: Inject orderService
	catalogHandler := handler.NewCatalogHandler(db, menuService, analyticsService, orderService, paymentService)

	// API v1
	v1 := router.Group("/api/v1")
	v1.Use(middleware.AuditMiddleware(db))

	// Dynamic Midtrans config for frontend client
	v1.GET("/config/midtrans", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"client_key":  cfg.MidtransClientKey,
			"environment": cfg.MidtransEnv,
		})
	})

	// Auth
	auth := v1.Group("/auth")
	auth.POST("/login", authHandler.Login)
	auth.POST("/login-pin", authHandler.LoginPIN)

	// Employees
	employees := v1.Group("/employees")
	employees.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	employees.GET("", employeeHandler.GetAll)
	employees.GET("/:id", employeeHandler.GetByID)
	employees.POST("", middleware.RequireRole(models.RoleManager), employeeHandler.Create)
	employees.PUT("/:id", middleware.RequireRole(models.RoleManager), employeeHandler.Update)
	employees.DELETE("/:id", middleware.RequireRole(models.RoleManager), employeeHandler.Delete)

	// Categories
	categories := v1.Group("/categories")
	categories.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	categories.GET("", categoryHandler.GetAll)
	categories.GET("/:id", categoryHandler.GetByID)
	categories.POST("", middleware.RequireRole(models.RoleManager), categoryHandler.Create)
	categories.PUT("/:id", middleware.RequireRole(models.RoleManager), categoryHandler.Update)
	categories.DELETE("/:id", middleware.RequireRole(models.RoleManager), categoryHandler.Delete)

	// Menu
	menus := v1.Group("/menus")
	menus.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	menus.GET("", menuHandler.GetActive)
	menus.GET("/all", menuHandler.GetAll)
	menus.GET("/:id", menuHandler.GetByID)
	menus.POST("", middleware.RequireRole(models.RoleManager), menuHandler.Create)
	menus.PUT("/:id", middleware.RequireRole(models.RoleManager, models.RoleKitchen), menuHandler.Update)
	menus.DELETE("/:id", middleware.RequireRole(models.RoleManager), menuHandler.Delete)

	// Customers
	customers := v1.Group("/customers")
	customers.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	customers.GET("", customerHandler.GetAll)
	customers.GET("/:id", customerHandler.GetByID)
	customers.POST("", customerHandler.Create)
	customers.PUT("/:id", customerHandler.Update)
	customers.DELETE("/:id", middleware.RequireRole(models.RoleManager), customerHandler.Delete)

	// Orders
	orders := v1.Group("/orders")
	orders.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	orders.GET("", orderHandler.GetAll)
	orders.GET("/by-number/:order_number", orderHandler.GetByNumber) // FIX: Tambah sebelum /:id
	orders.GET("/:id", orderHandler.GetByID)
	orders.POST("", orderHandler.Create)
	orders.PATCH("/:id/notes", orderHandler.UpdateNotes)
	orders.PATCH("/:id/confirm", orderHandler.Confirm)
	orders.PATCH("/:id/cancel", orderHandler.Cancel)
	orders.PATCH("/:id/void-item", orderHandler.VoidItem)
	orders.PATCH("/:id/void", orderHandler.VoidOrder)
	orders.PATCH("/:id/status", orderHandler.UpdateStatus)

	// Payments
	payments := v1.Group("/payments")
	payments.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	payments.POST("", paymentHandler.Process)
	payments.GET("/:id", paymentHandler.GetByID)
	payments.GET("/order/:order_id", paymentHandler.GetByOrderID)

	// Returns
	returns := v1.Group("/returns")
	returns.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	returns.POST("", returnHandler.ProcessReturn)

	// Shifts
	shifts := v1.Group("/shifts")
	shifts.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	shifts.POST("/open", shiftHandler.Open)
	shifts.GET("/current", shiftHandler.GetCurrent)
	shifts.PATCH("/close", shiftHandler.Close)
	shifts.POST("/:id/cash-movement", cashMovementHandler.Add)
	shifts.GET("/:id/cash-movements", cashMovementHandler.GetByShift)

	// Inventory
	inventory := v1.Group("/inventory")
	inventory.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	inventory.GET("/raw-materials", inventoryHandler.GetAllRawMaterials)
	inventory.GET("/raw-materials/:id", inventoryHandler.GetRawMaterial)
	inventory.POST("/raw-materials", inventoryHandler.CreateRawMaterial)
	inventory.PUT("/raw-materials/:id", inventoryHandler.UpdateRawMaterial)
	inventory.DELETE("/raw-materials/:id", inventoryHandler.DeleteRawMaterial)
	inventory.GET("/recipes/:menu_id", inventoryHandler.GetRecipeByMenu)
	inventory.POST("/recipes", inventoryHandler.UpsertRecipe)
	inventory.PUT("/recipes/:id", inventoryHandler.UpdateRecipe)
	inventory.DELETE("/recipes/:id", inventoryHandler.DeleteRecipe)
	inventory.GET("/stock-estimation", inventoryHandler.StockEstimation)

	// Stock
	stocks := v1.Group("/stock")
	stocks.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	stocks.POST("/adjust", stockHandler.Adjust)
	stocks.GET("/adjustments/:id", stockHandler.GetByRawMaterial)

	// Tax
	taxes := v1.Group("/taxes")
	taxes.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	taxes.GET("", taxHandler.GetActive)
	taxes.GET("/:id", taxHandler.GetByID)
	taxes.POST("", middleware.RequireRole(models.RoleManager), taxHandler.Create)
	taxes.PUT("/:id", middleware.RequireRole(models.RoleManager), taxHandler.Update)
	taxes.DELETE("/:id", middleware.RequireRole(models.RoleManager), taxHandler.Delete)

	// Discounts
	discounts := v1.Group("/discounts")
	discounts.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	discounts.GET("", discountHandler.GetActive)
	discounts.GET("/:id", discountHandler.GetByID)
	discounts.POST("", middleware.RequireRole(models.RoleManager), discountHandler.Create)
	discounts.PUT("/:id", middleware.RequireRole(models.RoleManager), discountHandler.Update)
	discounts.DELETE("/:id", middleware.RequireRole(models.RoleManager), discountHandler.Delete)
	discounts.POST("/apply", discountHandler.Apply)

	// Reservations
	reservations := v1.Group("/reservations")
	reservations.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	reservations.GET("", reservationHandler.GetByBranch)
	reservations.GET("/:id", reservationHandler.GetByID)
	reservations.POST("", reservationHandler.Create)
	reservations.PUT("/:id", reservationHandler.Update)
	reservations.DELETE("/:id", reservationHandler.Delete)

	v1.GET("/best-sellers", middleware.AuthMiddleware(cfg.JWTSecret), analyticsHandler.BestSellers)

	// Analytics — Legacy (on-demand computation)
	analytics := v1.Group("/analytics")
	analytics.Use(middleware.AuthMiddleware(cfg.JWTSecret), middleware.RequireRole(models.RoleManager))
	analytics.GET("/sales-summary", analyticsHandler.SalesSummary)
	analytics.GET("/best-sellers", analyticsHandler.BestSellers)
	analytics.GET("/return-impact", analyticsHandler.ReturnImpact)
	analytics.GET("/cogs", analyticsHandler.COGS)

	// Analytics — NEW: Materialized View Endpoints (instant dashboard data)
	analytics.GET("/daily-sales", analyticsHandler.DailySales)
	analytics.GET("/weekly-sales", analyticsHandler.WeeklySales)
	analytics.GET("/payment-summary", analyticsHandler.PaymentSummary)
	analytics.GET("/menu-profitability", analyticsHandler.MenuProfitability)
	analytics.GET("/shift-reconciliation", analyticsHandler.ShiftReconciliation)
	analytics.GET("/void-return-logs", analyticsHandler.VoidAndReturnLogs)
	analytics.GET("/cashier-performance", analyticsHandler.CashierPerformance)
	analytics.GET("/audit-logs", analyticsHandler.AuditLogs)

	// Printer
	printer := v1.Group("/printer")
	printer.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	printer.GET("/receipt/:order_id", printerHandler.PrintReceipt)
	printer.GET("/kitchen-ticket/:order_id", printerHandler.PrintKitchen)

	// Catalog Page under basic auth
	catalogGroup := router.Group("/catalog", gin.BasicAuth(gin.Accounts{
		"catalog": "kopitiam123",
	}))

	// Serving the static files directly from Gin
	catalogGroup.StaticFile("", "./FE/kopitiam-pos-fixed/catalog.html")
	catalogGroup.StaticFile("/index.html", "./FE/kopitiam-pos-fixed/catalog.html")
	catalogGroup.StaticFile("/catalog.js", "./FE/kopitiam-pos-fixed/catalog.js")
	catalogGroup.StaticFile("/global.css", "./FE/kopitiam-pos-fixed/global.css")
	catalogGroup.StaticFile("/midtrans-loader.js", "./FE/kopitiam-pos-fixed/midtrans-loader.js")

	// API routes for catalog
	catalogAPI := catalogGroup.Group("/api")
	catalogAPI.GET("/menus", catalogHandler.GetMenus)
	catalogAPI.POST("/checkout", catalogHandler.Checkout)
}
