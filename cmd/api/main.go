package main

import (
	_ "github.com/Lapnes/pos-kopitiam-v3nf/docs"
	"os"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/config"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/fatih/color"
)

// @title POS KopiTiam API
// @version 3.0
// @description Enterprise POS system for coffee shop chain
// @host localhost:8080
// @BasePath /api/v1
func main() {
	cfg := config.LoadConfig()
	dbCfg := config.DefaultDBConfig()
	db, err := config.InitDB(dbCfg)
	if err != nil {
		color.Red("Failed to initialize database: %v", err)
		os.Exit(1)
	}
	
	// Temporarily initialize Redis only if needed by routes (as we dropped Redis logic in some refactors)
	// But since routes.SetupRoutes still requires it, we initialize it here.
	redisClient := config.InitRedis(cfg)

	router := gin.Default()
	router.Use(cors.Default())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "env": cfg.AppEnv})
	})

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	routes.SetupRoutes(router, db, redisClient, cfg)

	port := cfg.AppPort
	if port == "" {
		port = "8080"
	}
	color.Green("Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		color.Red("Failed to start server: %v", err)
		os.Exit(1)
	}
}
