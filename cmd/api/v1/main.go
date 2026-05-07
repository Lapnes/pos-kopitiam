package main

import (
	"log"
	"net/http"

	"github.com/Lapnes/pos-kopitiam/internal/config"
	"github.com/Lapnes/pos-kopitiam/internal/dto"
	"github.com/Lapnes/pos-kopitiam/internal/models"
	"github.com/Lapnes/pos-kopitiam/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	config.InitConfig()
	db := config.ConnectDB()

	// Seed hanya boleh di dev
	if config.GetEnv("APP_ENV", "dev") == "dev" {
		config.SeedDatabase(db)
	}

	r := gin.Default()

	// ======================
	// CORS — izinkan semua origin (untuk dev)
	// ======================
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	// ======================
	// STATIC FILES — serve folder /static
	// ======================
	r.Static("/static", "./static")

	// ======================
	// HEALTH CHECK
	// ======================
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "POS KopiTiam API is running",
		})
	})

	// ======================
	// ORDER ENDPOINTS
	// ======================
	r.POST("/api/v1/orders", createOrderHandler(db))
	r.GET("/api/v1/orders", getAllOrdersHandler(db))
	r.GET("/api/v1/orders/report", getOrderReportHandler(db)) // harus sebelum :id
	r.GET("/api/v1/orders/:id", getOrderByIDHandler(db))
	r.DELETE("/api/v1/orders/:id", deleteOrderHandler(db))

	// ======================
	// MASTER DATA
	// ======================
	r.GET("/api/v1/menus", getMenusHandler(db))
	r.GET("/api/v1/categories", getCategoriesHandler(db))
	r.GET("/api/v1/employees", getEmployeesHandler(db))

	port := config.GetEnv("PORT", "3400")
	log.Printf("🚀 Server running on http://localhost:%s", port)
	log.Printf("🌐 Kasir Web: http://localhost:%s/static/pos-kopitiam-kasir.html", port)
	r.Run(":" + port)
}

//////////////////////////////////////////////////
// HANDLERS
//////////////////////////////////////////////////

func createOrderHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req struct {
			EmployeeID *uint64 `json:"employee_id"`
			Items      []struct {
				MenuID   uint64 `json:"menu_id" binding:"required"`
				Quantity int    `json:"quantity" binding:"required,min=1"`
			} `json:"items" binding:"required,min=1"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var items []struct {
			MenuID uint64
			Qty    int
		}

		for _, i := range req.Items {
			items = append(items, struct {
				MenuID uint64
				Qty    int
			}{
				MenuID: i.MenuID,
				Qty:    i.Quantity,
			})
		}

		order, err := service.CreateOrder(db, req.EmployeeID, items)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status":   "success",
			"order_id": order.OrderID,
		})
	}
}

func deleteOrderHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("id")

		// Cek order ada
		var order models.Order
		if err := db.First(&order, "order_id = ?", orderID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		// Hapus order_details dulu (FK constraint), lalu order
		if err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("order_id = ?", orderID).Delete(&models.OrderDetail{}).Error; err != nil {
				return err
			}
			if err := tx.Where("order_id = ?", orderID).Delete(&models.Order{}).Error; err != nil {
				return err
			}
			return nil
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Order deleted",
		})
	}
}

func getAllOrdersHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var orders []models.Order

		if err := db.Preload("OrderDetails.Menu").
			Preload("Employee").
			Find(&orders).Error; err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"count":  len(orders),
			"data":   orders,
		})
	}
}

func getOrderByIDHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("id")

		var order models.Order
		if err := db.Preload("OrderDetails.Menu").
			Preload("Employee").
			First(&order, "order_id = ?", orderID).Error; err != nil {

			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   order,
		})
	}
}

func getOrderReportHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reports []dto.OrderReportDTO

		err := db.Table("orders").
			Select(`
				orders.order_id,
				COALESCE(employees.employee_name, 'Walk-in') as employee_name,
				menus.menu_name,
				order_details.quantity,
				order_details.unit_price,
				order_details.subtotal,
				orders.total_price
			`).
			Joins("LEFT JOIN employees ON employees.employee_id = orders.employee_id").
			Joins("JOIN order_details ON order_details.order_id = orders.order_id").
			Joins("JOIN menus ON menus.menu_id = order_details.menu_id").
			Order("orders.order_id DESC").
			Scan(&reports).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"count":  len(reports),
			"data":   reports,
		})
	}
}

func getMenusHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var menus []models.Menu

		if err := db.Preload("Category").Find(&menus).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"count":  len(menus),
			"data":   menus,
		})
	}
}

func getCategoriesHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var categories []models.Category

		if err := db.Preload("Menus").Find(&categories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"count":  len(categories),
			"data":   categories,
		})
	}
}

func getEmployeesHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var employees []models.Employee

		if err := db.Find(&employees).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"count":  len(employees),
			"data":   employees,
		})
	}
}