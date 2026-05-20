package main

import (
	"os"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/config"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/fatih/color"
)

func main() {
	// Load config to parse .env file
	_ = config.LoadConfig()

	dbCfg := config.DefaultDBConfig()
	db, err := config.InitDB(dbCfg)
	if err != nil {
		color.Red("InitDB Failed: %v", err)
		os.Exit(1)
	}

	color.Cyan("Dropping old tables (dependency-aware)...")

	// DROP ORDER: Junction tables first → Child → Parent
	// Use string names for junction tables (GORM many2many tables)
	db.Migrator().DropTable(
		"menu_categories",
		"recipe_ingredients", // FIX: Explicitly drop junction/child tables
		&models.DiscountApplication{},
		&models.PaymentSplit{},
		&models.OrderReturn{},
		&models.OrderDetail{},
		&models.AuditLog{},
		&models.Reservation{},
		&models.StockAdjustment{},
		&models.CashMovement{},
		&models.Payment{},
		&models.Order{},
		&models.Discount{},
		&models.Recipe{},
		&models.RecipeIngredient{}, // FIX: Ensure this is dropped
		&models.RawMaterial{},
		&models.Menu{},
		&models.Category{},
		&models.TaxConfig{},
		&models.Customer{},
		&models.Shift{},
		&models.User{},
		&models.Branch{},
		&models.Branch{},
	)

	color.Cyan("Running Auto Migration (CREATE order: Parent → Child → Junction)...")

	err = db.AutoMigrate(
		// Parent tables
		&models.Branch{},
		&models.User{},
		&models.Customer{},
		&models.Shift{},
		&models.TaxConfig{},
		&models.Category{},
		&models.Menu{},
		// Child tables - FIX: RecipeIngredient BEFORE Recipe so FK works
		&models.RawMaterial{},
		&models.RecipeIngredient{}, // FIX: Create ingredients table first
		&models.Recipe{},
		&models.Order{},
		&models.OrderDetail{},
		&models.Payment{},
		&models.PaymentSplit{},
		&models.OrderReturn{},
		&models.Discount{},
		&models.DiscountApplication{},
		&models.Reservation{},
		&models.StockAdjustment{},
		&models.CashMovement{},
		&models.AuditLog{},
	)

	if err != nil {
		color.Red("Auto Migration Failed: %v", err)
		os.Exit(1)
	}

	// FIX: Explicitly create junction tables that GORM many2many might miss
	// Ensure recipe_ingredients table exists with proper schema
	if !db.Migrator().HasTable("recipe_ingredients") {
		color.Yellow("Creating recipe_ingredients table explicitly...")
		if err := db.Migrator().CreateTable(&models.RecipeIngredient{}); err != nil {
			color.Red("Warning: Could not create recipe_ingredients: %v", err)
		}
	}

	// FIX: Ensure menu_categories junction table exists
	if !db.Migrator().HasTable("menu_categories") {
		color.Yellow("Creating menu_categories junction table explicitly...")
		// GORM creates this automatically for many2many, but let's be safe
		type MenuCategory struct {
			MenuID     string `gorm:"primaryKey"`
			CategoryID string `gorm:"primaryKey"`
		}
		if err := db.Migrator().CreateTable(&MenuCategory{}); err != nil {
			color.Red("Warning: Could not create menu_categories: %v", err)
		}
	}

	color.Green("Auto Migration Completed Successfully!")
	color.Yellow("Next: go run seeds/seed.go")
}
