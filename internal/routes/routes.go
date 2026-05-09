package routes

import (
	"net/http"

	"github.com/Lapnes/pos-kopitiam/internal/config"
	"github.com/Lapnes/pos-kopitiam/internal/handler"
	"github.com/Lapnes/pos-kopitiam/internal/middleware"
	"github.com/Lapnes/pos-kopitiam/internal/repository"
	"github.com/Lapnes/pos-kopitiam/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func InitRoutes(router *gin.Engine, db *gorm.DB, redisClient *redis.Client, cfg *config.Config) {
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
	orderService := service.NewOrderService(db, orderRepo, inventoryRepo, redisClient)
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
	menuHandler := handler.NewMenuHandler(menuService)
	employeeHandler := handler.NewEmployeeHandler(employeeService)

	// Health Check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "up",
			"redis":  redisClient.Ping(c).Val() == "PONG",
			"db":     db.Error == nil,
		})
	})

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
	
	orders := protected.Group("/orders")
	{
		orders.POST("", orderHandler.CreateOrder)
		orders.GET("", orderHandler.GetOrders)
		orders.GET("/:id", orderHandler.GetOrder)
		orders.PUT("/:id", orderHandler.UpdateOrder)
		orders.PUT("/:id/confirm", orderHandler.ConfirmOrder)
		orders.PUT("/:id/cancel", orderHandler.CancelOrder)
		
		// Manager only routes
		managerRoutes := orders.Group("/")
		managerRoutes.Use(middleware.RoleMiddleware("superadmin", "manager"))
		managerRoutes.POST("/:id/void-item", orderHandler.VoidItem)
		managerRoutes.POST("/:id/void", orderHandler.VoidOrder)

		// Cashier allowed routes (returns and lookup)
		cashierRoutes := orders.Group("/")
		cashierRoutes.Use(middleware.RoleMiddleware("superadmin", "manager", "cashier"))
		cashierRoutes.POST("/:id/returns", returnHandler.ProcessReturn)
		cashierRoutes.GET("/by-number/:number", orderHandler.GetOrderByNumber)
	}

	// Shifts
	shifts := protected.Group("/shifts")
	shifts.Use(middleware.RoleMiddleware("cashier", "manager", "superadmin"))
	{
		shifts.POST("/open", shiftHandler.OpenShift)
		shifts.GET("/current", shiftHandler.GetCurrentShift)
		shifts.POST("/close", shiftHandler.CloseShift)
	}

	// Master Data
	master := protected.Group("/master")
	master.Use(middleware.RoleMiddleware("superadmin", "admin", "manager"))
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
	employees.Use(middleware.RoleMiddleware("superadmin", "manager"))
	{
		employees.GET("", employeeHandler.GetEmployees)
		employees.POST("", employeeHandler.CreateEmployee)
		employees.PUT("/:id", employeeHandler.UpdateEmployee)
		employees.DELETE("/:id", employeeHandler.DeleteEmployee)
	}

	// Analytics
	analytics := protected.Group("/analytics")
	analytics.Use(middleware.RoleMiddleware("superadmin", "manager"))
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

	// Additional feature groups (payments, KDS, reservations, etc.) would be added here similarly.
}
