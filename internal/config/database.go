package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBConfig holds all database-related configuration
type DBConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	SlowThreshold   time.Duration
	LogLevel        logger.LogLevel
}

// DefaultDBConfig returns production-ready defaults
func DefaultDBConfig() *DBConfig {
	return &DBConfig{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            getEnvInt("DB_PORT", 3306),
		User:            getEnv("DB_USER", "root"),
		Password:        getEnv("DB_PASSWORD", ""),
		DBName:          getEnv("DB_NAME", "pos_kopitiam"),
		SSLMode:         getEnv("DB_SSLMODE", "disable"),
		MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),                // Max simultaneous connections
		MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),                // Keep warm in pool
		ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", time.Hour), // Recycle connections
		ConnMaxIdleTime: getEnvDuration("DB_CONN_MAX_IDLE_TIME", 10*time.Minute),
		SlowThreshold:   getEnvDuration("DB_SLOW_THRESHOLD", 200*time.Millisecond),
		LogLevel:        getEnvLogLevel("DB_LOG_LEVEL", logger.Warn), // Default: warn
	}
}

// DevelopmentDBConfig returns dev-friendly settings with SQL logging
func DevelopmentDBConfig() *DBConfig {
	cfg := DefaultDBConfig()
	// If DB_LOG_LEVEL is not set in env, use Info for development
	if os.Getenv("DB_LOG_LEVEL") == "" {
		cfg.LogLevel = logger.Info
	}
	cfg.SlowThreshold = 100 * time.Millisecond
	return cfg
}

// DSN builds MySQL connection string
func (c *DBConfig) DSN() string {
	if strings.HasPrefix(c.Host, "/") || strings.HasSuffix(c.Host, ".sock") {
		return fmt.Sprintf(
			"%s:%s@unix(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.User, c.Password, c.Host, c.DBName,
		)
	}
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName,
	)
}

// InitDB initializes database with connection pool, context support, and retry logic
func InitDB(cfg *DBConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	// Retry logic for transient connection failures
	maxRetries := 5
	retryDelay := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				logger.Config{
					SlowThreshold:             cfg.SlowThreshold,
					LogLevel:                  cfg.LogLevel,
					IgnoreRecordNotFoundError: true,
					Colorful:                  os.Getenv("APP_ENV") != "prod",
				},
			),
			// PrepareStmt: true, // Enable prepared statement cache for performance
		})
		if err == nil {
			break
		}
		log.Printf("[DB] Connection attempt %d/%d failed: %v. Retrying in %v...", i+1, maxRetries, err, retryDelay)
		time.Sleep(retryDelay)
		retryDelay *= 2 // Exponential backoff
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	// Verify connection with context (timeout support)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	log.Printf("[DB] Connected. Pool: max_open=%d, max_idle=%d, max_lifetime=%v",
		cfg.MaxOpenConns, cfg.MaxIdleConns, cfg.ConnMaxLifetime)

	return db, nil
}

// HealthCheck verifies DB connectivity with context
func HealthCheck(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return sqlDB.PingContext(ctx)
}

// GetDBStats returns connection pool statistics for monitoring
func GetDBStats(db *gorm.DB) (map[string]interface{}, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration_ms":     stats.WaitDuration.Milliseconds(),
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}, nil
}

// TransactionWithRetry executes a function within a transaction with deadlock retry
func TransactionWithRetry(db *gorm.DB, fn func(tx *gorm.DB) error, maxRetries int) error {
	var err error
	retryDelay := 50 * time.Millisecond

	for i := 0; i < maxRetries; i++ {
		err = db.Transaction(fn)
		if err == nil {
			return nil
		}

		// Check if it's a deadlock or lock timeout (retryable)
		errStr := err.Error()
		if !isRetryableError(errStr) {
			return err // Non-retryable error, fail fast
		}

		log.Printf("[TX] Retryable error (attempt %d/%d): %v. Retrying in %v...", i+1, maxRetries, err, retryDelay)
		time.Sleep(retryDelay)
		retryDelay *= 2 // Exponential backoff
	}

	return fmt.Errorf("transaction failed after %d retries: %w", maxRetries, err)
}

// isRetryableError checks if error is deadlock or lock timeout
func isRetryableError(err string) bool {
	retryable := []string{
		"deadlock detected",
		"lock timeout",
		"could not serialize access",
		"retry transaction",
		"connection reset by peer",
		"broken pipe",
	}
	for _, r := range retryable {
		if containsIgnoreCase(err, r) {
			return true
		}
	}
	return false
}

func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && (findSubstr(s, substr) >= 0 ||
		findSubstr(s, toUpper(substr)) >= 0 ||
		findSubstr(s, toLower(substr)) >= 0)
}

func findSubstr(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func toUpper(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			c = c - 'a' + 'A'
		}
		result[i] = c
	}
	return string(result)
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c = c - 'A' + 'a'
		}
		result[i] = c
	}
	return string(result)
}

// Helper functions for env vars

func getEnvInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return defaultVal
}

func getEnvLogLevel(key string, defaultVal logger.LogLevel) logger.LogLevel {
	if v := os.Getenv(key); v != "" {
		switch toLower(v) {
		case "silent":
			return logger.Silent
		case "error":
			return logger.Error
		case "warn":
			return logger.Warn
		case "info":
			return logger.Info
		}
	}
	return defaultVal
}
