package config

import (
	"fmt"
	"log"

	"github.com/Lapnes/pos-kopitiam/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		GetEnv("DB_USER", "root"),
		GetEnv("DB_PASSWORD", ""),
		GetEnv("DB_HOST", "127.0.0.1"),
		GetEnv("DB_PORT", "3306"),
		GetEnv("DB_NAME", "pos_kopitiam"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		// DisableForeignKeyConstraintWhenMigrating: true, // Penting!
	})

	if err != nil {
		log.Fatalf("[NO] Gagal koneksi DB: %v", err)
	}

	log.Println("[YES] Database Connected!")

	// Migrasi semua tabel
	err = db.AutoMigrate(
		&models.Category{},
		&models.Employee{},
		&models.Menu{},
		&models.Order{},
		&models.OrderDetail{},
	)
	if err != nil {
		log.Fatalf("[NO] AutoMigrate gagal: %v", err)
	}

	log.Println("[SUCCESS] AutoMigrate Success!")
	return db
}
