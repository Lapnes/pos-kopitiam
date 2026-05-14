package main

import (
	"log"
	"time"

	"github.com/Lapnes/pos-kopitiam/internal/config"
	"github.com/Lapnes/pos-kopitiam/internal/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// @title           KopiTiam POS API
// @version         1.0
// @description     This is a POS API for KopiTiam Enterprise.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	cfg := config.LoadConfig()
	db := config.InitDB(cfg)
	redisClient := config.InitRedis(cfg)

	if cfg.AppEnv == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// CORS Middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
			"http://localhost:3001",   // Next.js dev fallback port
			"http://127.0.0.1:3001",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Init Routes
	routes.InitRoutes(router, db, redisClient, cfg)

	log.Printf("Starting Server on port %s", cfg.AppPort)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
