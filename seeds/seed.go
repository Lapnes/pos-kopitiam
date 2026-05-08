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
	category := models.Category{Name: "Makanan", Station: models.StationKitchen}
	db.FirstOrCreate(&category, models.Category{Name: "Makanan"})

	// Create Menus
	menus := []models.Menu{
		{CategoryID: category.ID, Name: "Nasi Goreng Spesial", Price: 35000, CostPrice: 20000, DailyStock: 50, Station: models.StationKitchen},
		{CategoryID: category.ID, Name: "Mie Goreng", Price: 30000, CostPrice: 15000, DailyStock: 50, Station: models.StationKitchen},
	}
	for _, menu := range menus {
		db.FirstOrCreate(&menu, models.Menu{Name: menu.Name})
	}

	log.Println("Seeding Completed!")
}
