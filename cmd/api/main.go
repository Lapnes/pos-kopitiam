package main

import (
	"log"

	"github.com/Lapnes/pos-kopitiam/internal/config"
	"github.com/Lapnes/pos-kopitiam/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	db := config.InitDB(cfg)
	redisClient := config.InitRedis(cfg)

	if cfg.AppEnv == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Use Middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Init Routes
	routes.InitRoutes(router, db, redisClient, cfg)

	log.Printf("Starting Server on port %s", cfg.AppPort)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
