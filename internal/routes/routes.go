package routes

import (
	"net/http"

	"github.com/Lapnes/pos-kopitiam/internal/config"
	"github.com/Lapnes/pos-kopitiam/internal/handler"
	"github.com/Lapnes/pos-kopitiam/internal/middleware"
	"github.com/Lapnes/pos-kopitiam/internal/repository"
	"github.com/Lapnes/pos-kopitiam/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	_ "github.com/Lapnes/pos-kopitiam/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRoutes(router *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// Initialize Redis for health check
	redisClient := config.InitRedis(cfg)

	// Repositories
	userRepo := repository.NewUserRepo(db)
	orderRepo := repository.NewOrderRepo(db)
	shiftRepo := repository.NewShiftRepository(db)
	returnRepo := repository.NewReturnRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db)
	menuRepo := repository.NewMenuRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, cfg)
	// FIX: Pass shiftRepo to OrderService for ShiftID assignment during confirm
	orderService := service.NewOrderService(db, orderRepo, inventoryRepo, shiftRepo)
	shiftService := service.NewShiftService(shiftRepo)
	returnService := service.NewReturnService(db, returnRepo, orderRepo)
	inventoryService := service.NewInventoryService(inventoryRepo)
	analyticsService := service.NewAnalyticsService(analyticsRepo)
	printerService := service.NewPrinterService(orderRepo, userRepo)
	menuService := service.NewMenuService(menuRepo)
	employeeService := service.NewEmployeeService(userRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	orderHandler := handler.NewOrderHandler(orderService)
	shiftHandler := handler.NewShiftHandler(shiftService)
	returnHandler := handler.NewReturnHandler(returnService)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)
	printerHandler := handler.NewPrinterHandler(printerService)
	paymentHandler := handler.NewPaymentHandler(db)
	menuHandler := handler.NewMenuHandler(menuService)
	employeeHandler := handler.NewEmployeeHandler(employeeService)

	// Health Check — includes Redis status
	// HealthCheck godoc
	// @Summary      System Health Check
	// @Description  Returns status of DB and Redis connections.
	// @Tags         System
	// @Produce      json
	// @Success      200 {object} map[string]interface{} "System is up"
	// @Router       /health [get]
	router.GET("/health", func(c *gin.Context) {
		redisStatus := true
		if redisClient != nil {
			ctx := c.Request.Context()
			if err := redisClient.Ping(ctx).Err(); err != nil {
				redisStatus = false
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "up",
			"db":     db.Error == nil,
			"redis":  redisStatus,
		})
	})

	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")

	// Public routes
	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/login-pin", authHandler.LoginPIN)
		auth.POST("/refresh", authHandler.Refresh)
		auth.POST("/logout", authHandler.Logout)
	}

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	menus := protected.Group("/menus")
	{
		menus.GET("", menuHandler.GetActiveMenus)
	}

	// Orders — flat structure with inline role checks per route
	orders := protected.Group("/orders")
	{
		orders.POST("", orderHandler.CreateOrder)
		orders.GET("", orderHandler.GetOrders)
		orders.GET("/:id", orderHandler.GetOrder)
		orders.PUT("/:id", orderHandler.UpdateOrder)
		orders.PUT("/:id/confirm", orderHandler.ConfirmOrder)
		orders.PUT("/:id/cancel", orderHandler.CancelOrder)
		orders.POST("/:id/payments", paymentHandler.ProcessPayment)

		// Manager-only: void and void-item
		orders.POST("/:id/void", middleware.RoleMiddleware("manager"), orderHandler.VoidOrder)
		orders.POST("/:id/void-item", middleware.RoleMiddleware("manager"), orderHandler.VoidItem)

		// Cashier+Manager: returns and lookup
		orders.POST("/:id/returns", middleware.RoleMiddleware("cashier", "manager"), returnHandler.ProcessReturn)
		orders.GET("/by-number/:number", middleware.RoleMiddleware("cashier", "manager"), orderHandler.GetOrderByNumber)
	}

	// Shifts
	shifts := protected.Group("/shifts")
	shifts.Use(middleware.RoleMiddleware("cashier", "manager"))
	{
		shifts.POST("/open", shiftHandler.OpenShift)
		shifts.GET("/current", shiftHandler.GetCurrentShift)
		shifts.POST("/close", shiftHandler.CloseShift)
	}

	// Master Data
	master := protected.Group("/master")
	master.Use(middleware.RoleMiddleware("manager"))
	{
		master.POST("/raw-materials", inventoryHandler.CreateRawMaterial)
		master.PUT("/raw-materials/:id", inventoryHandler.UpdateRawMaterial)
		master.GET("/raw-materials/:id", inventoryHandler.GetRawMaterial)
		master.DELETE("/raw-materials/:id", inventoryHandler.DeleteRawMaterial)

		master.POST("/recipes", inventoryHandler.CreateRecipe)
		master.PUT("/recipes/:id", inventoryHandler.UpdateRecipe)
		master.GET("/recipes/:menu_id", inventoryHandler.GetRecipe)
		master.DELETE("/recipes/:id", inventoryHandler.DeleteRecipe)
	}

	// Employees
	employees := protected.Group("/employees")
	employees.Use(middleware.RoleMiddleware("manager"))
	{
		employees.GET("", employeeHandler.GetEmployees)
		employees.POST("", employeeHandler.CreateEmployee)
		employees.PUT("/:id", employeeHandler.UpdateEmployee)
		employees.DELETE("/:id", employeeHandler.DeleteEmployee)
	}

	// Analytics
	analytics := protected.Group("/analytics")
	analytics.Use(middleware.RoleMiddleware("manager"))
	{
		analytics.GET("/sales-summary", analyticsHandler.GetSalesSummary)
		analytics.GET("/best-sellers", analyticsHandler.GetBestSellers)
		analytics.GET("/return-impact", analyticsHandler.GetReturnImpact)
	}

	// Printers
	printers := protected.Group("/printers")
	{
		printers.GET("/receipt/:id", printerHandler.GetReceiptPayload)
		printers.GET("/kitchen/:id", printerHandler.GetKitchenTicketPayload)
	}
}
