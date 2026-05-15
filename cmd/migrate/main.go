package main

import (
	"log"

	"github.com/Lapnes/pos-kopitiam/internal/config"
	"github.com/Lapnes/pos-kopitiam/internal/models"
)

func main() {
	cfg := config.LoadConfig()
	db := config.InitDB(cfg)

	log.Println("Dropping old tables to prevent schema conflicts...")
	db.Migrator().DropTable(
		"menu_categories", // Drop M2M join table first
		&models.AuditLog{},
		&models.Reservation{},
		&models.DiscountApplication{},
		&models.Discount{},
		&models.PaymentSplit{},
		&models.Payment{},
		&models.OrderReturn{},
		&models.OrderDetail{},
		&models.Order{},
		&models.Shift{},
		&models.Employee{},
		&models.RecipeItem{},
		&models.Recipe{},
		&models.RawMaterial{},
		&models.Menu{},
		&models.Category{},
		&models.Table{},
		&models.Area{},
		&models.Branch{},
	)

	log.Println("Running Auto Migration...")

	err := db.AutoMigrate(
		&models.Branch{},
		&models.Area{},
		&models.Table{},
		&models.Category{},
		&models.RawMaterial{},
		&models.Menu{},       // GORM will auto-create menu_categories join table
		&models.Recipe{},
		&models.RecipeItem{},
		&models.Employee{},
		&models.Shift{},
		&models.Order{},
		&models.OrderDetail{},
		&models.OrderReturn{},
		&models.Payment{},
		&models.PaymentSplit{},
		&models.Discount{},
		&models.DiscountApplication{},
		&models.Reservation{},
		&models.AuditLog{},
	)

	if err != nil {
		log.Fatalf("Auto Migration Failed: %v", err)
	}

	log.Println("Auto Migration Completed Successfully!")
}
