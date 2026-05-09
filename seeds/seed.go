package main

import (
	"log"

	"github.com/Lapnes/pos-kopitiam/internal/config"
	"github.com/Lapnes/pos-kopitiam/internal/models"
	"github.com/Lapnes/pos-kopitiam/internal/utils"
)

func main() {
	cfg := config.LoadConfig()
	db := config.InitDB(cfg)

	log.Println("Seeding Database...")

	// Create Branch
	mainBranch := models.Branch{
		Name:    "Main Branch",
		Address: "Jl. Sudirman No. 1",
		Phone:   "08123456789",
		IsMain:  true,
	}
	db.FirstOrCreate(&mainBranch, models.Branch{Name: "Main Branch"})

	// Create Users
	hashedPassword, _ := utils.HashPassword("password123")
	users := []models.Employee{
		{BranchID: mainBranch.ID, Name: "Super Admin", Email: "admin@kopitiam.com", Password: hashedPassword, PINCode: "111111", Role: models.RoleSuperadmin},
		{BranchID: mainBranch.ID, Name: "Manager", Email: "manager@kopitiam.com", Password: hashedPassword, PINCode: "222222", Role: models.RoleManager},
		{BranchID: mainBranch.ID, Name: "Cashier 1", Email: "cashier1@kopitiam.com", Password: hashedPassword, PINCode: "333333", Role: models.RoleCashier},
		{BranchID: mainBranch.ID, Name: "Kitchen 1", Email: "kitchen1@kopitiam.com", Password: hashedPassword, PINCode: "444444", Role: models.RoleKitchen},
	}

	for _, user := range users {
		db.FirstOrCreate(&user, models.Employee{Email: user.Email})
		// Force update to save the new hashed passwords and PINs
		db.Model(&user).Updates(map[string]interface{}{
			"password": user.Password,
			"pin_code": user.PINCode,
		})
	}

	// Create Areas
	area := models.Area{BranchID: mainBranch.ID, Name: "Indoor"}
	db.FirstOrCreate(&area, models.Area{Name: "Indoor", BranchID: mainBranch.ID})

	// Create Tables
	tables := []models.Table{
		{AreaID: area.ID, Name: "Table 1", Capacity: 4},
		{AreaID: area.ID, Name: "Table 2", Capacity: 2},
	}
	for _, table := range tables {
		db.FirstOrCreate(&table, models.Table{Name: table.Name, AreaID: area.ID})
	}

	// Create Categories
	categoriesData := []struct {
		Name    string
		Station models.Station
	}{
		{"Coffee", models.StationBar},
		{"Tea", models.StationBar},
		{"Food", models.StationKitchen},
		{"Cold Beverages", models.StationBar},
		{"Snacks", models.StationKitchen},
		{"Desserts", models.StationKitchen},
		{"Juices", models.StationBar},
	}

	categoriesMap := make(map[string]models.Category)
	for _, c := range categoriesData {
		cat := models.Category{Name: c.Name, Station: c.Station}
		db.FirstOrCreate(&cat, models.Category{Name: c.Name})
		categoriesMap[c.Name] = cat
	}

	// Create Menus based on pos_kopitiam_full.sql
	menus := []models.Menu{
		{CategoryID: categoriesMap["Coffee"].ID, Name: "Milk Coffee", Price: 18000, CostPrice: 10000, DailyStock: 48, IsActive: true, Station: models.StationBar},
		{CategoryID: categoriesMap["Coffee"].ID, Name: "Black Coffee", Price: 15000, CostPrice: 8000, DailyStock: 50, IsActive: true, Station: models.StationBar},
		{CategoryID: categoriesMap["Tea"].ID, Name: "Pulled Tea", Price: 15000, CostPrice: 8000, DailyStock: 39, IsActive: true, Station: models.StationBar},
		{CategoryID: categoriesMap["Tea"].ID, Name: "Plain Tea", Price: 8000, CostPrice: 4000, DailyStock: 40, IsActive: true, Station: models.StationBar},
		{CategoryID: categoriesMap["Food"].ID, Name: "Fried Rice", Price: 35000, CostPrice: 20000, DailyStock: 20, IsActive: true, Station: models.StationKitchen},
		{CategoryID: categoriesMap["Food"].ID, Name: "Fried Noodles", Price: 30000, CostPrice: 15000, DailyStock: 20, IsActive: true, Station: models.StationKitchen},
		{CategoryID: categoriesMap["Cold Beverages"].ID, Name: "Iced Milk Coffee", Price: 20000, CostPrice: 12000, DailyStock: 30, IsActive: true, Station: models.StationBar},
		{CategoryID: categoriesMap["Snacks"].ID, Name: "French Fries", Price: 18000, CostPrice: 10000, DailyStock: 25, IsActive: true, Station: models.StationKitchen},
		{CategoryID: categoriesMap["Desserts"].ID, Name: "Chocolate Pudding", Price: 15000, CostPrice: 8000, DailyStock: 15, IsActive: true, Station: models.StationKitchen},
		{CategoryID: categoriesMap["Juices"].ID, Name: "Avocado Juice", Price: 22000, CostPrice: 12000, DailyStock: 20, IsActive: true, Station: models.StationBar},
	}
	for _, menu := range menus {
		db.FirstOrCreate(&menu, models.Menu{Name: menu.Name})
	}

	log.Println("Seeding Completed!")
}
