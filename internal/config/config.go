package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv     string
	AppPort    string
	AppSecret  string
	
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       string
	
	BigcapitalBaseURL string
	BigcapitalAPIKey  string
	BigcapitalOrgID   string
	
	JWTSecret         string
	JWTExpiry         string
	RefreshTokenExpiry string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading it. Using environment variables.")
	}

	return &Config{
		AppEnv:     getEnv("APP_ENV", "dev"),
		AppPort:    getEnv("APP_PORT", "8080"),
		AppSecret:  getEnv("APP_SECRET", "secret"),
		
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "pos_kopitiam"),
		
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnv("REDIS_DB", "0"),
		
		BigcapitalBaseURL: getEnv("BIGCAPITAL_BASE_URL", "https://api.bigcapital.ly"),
		BigcapitalAPIKey:  getEnv("BIGCAPITAL_API_KEY", ""),
		BigcapitalOrgID:   getEnv("BIGCAPITAL_ORG_ID", ""),
		
		JWTSecret:         getEnv("JWT_SECRET", "super-secret-jwt-key"),
		JWTExpiry:         getEnv("JWT_EXPIRY", "24h"),
		RefreshTokenExpiry: getEnv("REFRESH_TOKEN_EXPIRY", "168h"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
