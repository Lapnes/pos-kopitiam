package config

import (
	"os"

	"github.com/joho/godotenv"
)

func InitConfig() {
	if os.Getenv("APP_ENV") != "production" {
		_ = godotenv.Load()
	}
}

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
