package main

import (
	"log"

	"github.com/Lapnes/pos-kopitiam/internal/config"
)

func main() {
	config.InitConfig()
	db := config.ConnectDB()
	config.SeedDatabase(db)

	// // Tambah FK constraints
	// if err := config.AddForeignKeys(db); err != nil {
	// 	log.Fatalf("[NO] Gagal menambah FK: %v", err)
	// }
	// log.Println("[YES] Foreign Key Constraints added!")

	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	log.Println("[SUCCESS] Migration & Seed Success!")
}
