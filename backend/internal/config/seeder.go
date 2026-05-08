package config

import (
	"crypto/md5"
	"encoding/hex"
	"log"

	"github.com/Lapnes/pos-kopitiam/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func hashMD5(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

func SeedDatabase(db *gorm.DB) {
	log.Println("🌱 Seeding database...")

	// ======================
	// Seed Categories
	// ======================
	categories := []models.Category{
		{CategoryName: "Coffee"},
		{CategoryName: "Tea"},
		{CategoryName: "Food"},
		{CategoryName: "Cold Beverages"},
		{CategoryName: "Snacks"},
		{CategoryName: "Desserts"},
		{CategoryName: "Juices"},
		{CategoryName: "Smoothies"},
		{CategoryName: "Milkshakes"},
		{CategoryName: "Special Menu"},
	}

	for i := range categories {
		var existing models.Category
		result := db.Where("category_name = ?", categories[i].CategoryName).First(&existing)
		if result.Error != nil {
			db.Create(&categories[i])
		} else {
			categories[i] = existing
		}
	}

	// ======================
	// Seed Menus
	// ======================
	menus := []models.Menu{
		{CategoryID: categories[0].CategoryID, MenuName: "Milk Coffee", Price: 18000, DailyStock: 50},
		{CategoryID: categories[0].CategoryID, MenuName: "Black Coffee", Price: 15000, DailyStock: 50},
		{CategoryID: categories[1].CategoryID, MenuName: "Pulled Tea", Price: 15000, DailyStock: 40},
		{CategoryID: categories[1].CategoryID, MenuName: "Plain Tea", Price: 8000, DailyStock: 40},
		{CategoryID: categories[2].CategoryID, MenuName: "Fried Rice", Price: 35000, DailyStock: 20},
		{CategoryID: categories[2].CategoryID, MenuName: "Fried Noodles", Price: 30000, DailyStock: 20},
		{CategoryID: categories[3].CategoryID, MenuName: "Iced Milk Coffee", Price: 20000, DailyStock: 30},
		{CategoryID: categories[4].CategoryID, MenuName: "French Fries", Price: 18000, DailyStock: 25},
		{CategoryID: categories[5].CategoryID, MenuName: "Chocolate Pudding", Price: 15000, DailyStock: 15},
		{CategoryID: categories[6].CategoryID, MenuName: "Avocado Juice", Price: 22000, DailyStock: 20},
	}

	for i := range menus {
		var existing models.Menu
		result := db.Where("menu_name = ?", menus[i].MenuName).First(&existing)
		if result.Error != nil {
			db.Create(&menus[i])
		} else {
			menus[i] = existing
		}
	}

	// ======================
	// Seed Employees
	// ======================
	employees := []models.Employee{
		{EmployeeName: "Admin Tenzly", PasswordHash: hashMD5("admin123")},
		{EmployeeName: "Kasir Budi", PasswordHash: hashMD5("budi123")},
		{EmployeeName: "Kasir Siti", PasswordHash: hashMD5("siti123")},
	}

	for i := range employees {
		var existing models.Employee
		result := db.Where("employee_name = ?", employees[i].EmployeeName).First(&existing)
		if result.Error != nil {
			db.Create(&employees[i])
		} else {
			employees[i] = existing
		}
	}

	// ======================
	// Seed Orders
	// ======================
	orders := []models.Order{
		{OrderID: uuid.New().String(), EmployeeID: &employees[0].EmployeeID},
		{OrderID: uuid.New().String(), EmployeeID: &employees[1].EmployeeID},
		{OrderID: uuid.New().String(), EmployeeID: &employees[2].EmployeeID},
	}

	for i := range orders {
		db.Create(&orders[i])
	}

	// ======================
	// Seed Order Details
	// ======================
	orderDetails := []models.OrderDetail{
		{OrderID: orders[0].OrderID, MenuID: menus[0].MenuID, Quantity: 2, UnitPrice: menus[0].Price},
		{OrderID: orders[0].OrderID, MenuID: menus[5].MenuID, Quantity: 1, UnitPrice: menus[5].Price},

		{OrderID: orders[1].OrderID, MenuID: menus[1].MenuID, Quantity: 2, UnitPrice: menus[1].Price},

		{OrderID: orders[2].OrderID, MenuID: menus[2].MenuID, Quantity: 1, UnitPrice: menus[2].Price},
	}

	for i := range orderDetails {
		db.Create(&orderDetails[i])
	}

	// ======================
	// Update Total
	// ======================
	for i := range orders {
		var total int64
		db.Model(&models.OrderDetail{}).
			Where("order_id = ?", orders[i].OrderID).
			Select("SUM(subtotal)").Scan(&total)

		db.Model(&orders[i]).Update("total_price", total)
	}

	// ======================
	// Create Triggers
	// ======================
	SetupDatabaseLogic(db)

	log.Println("🎉 Seeding completed!")
}
