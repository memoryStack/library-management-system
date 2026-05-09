package initializers

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDB opens the Postgres connection using env vars set by LoadEnv.
func ConnectDB(environment string) error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		pass := os.Getenv("DB_PASSWORD")
		name := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")
		ssl := os.Getenv("DB_SSLMODE")
		if ssl == "" {
			ssl = "disable"
		}
		if host == "" || user == "" || name == "" || port == "" {
			return fmt.Errorf("set DATABASE_URL or DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT")
		}
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			host, user, pass, name, port, ssl)
	}

	cfg := &gorm.Config{}

	switch environment {
	case EnvDevelopment:
		cfg.Logger = logger.Default.LogMode(logger.Info)
	default:
		cfg.Logger = logger.Default.LogMode(logger.Error)
	}

	sqlDB, err := gorm.Open(postgres.Open(dsn), cfg)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	pool, err := sqlDB.DB()
	if err != nil {
		return fmt.Errorf("sql db: %w", err)
	}

	maxOpen := 25
	if v := os.Getenv("DB_MAX_OPEN_CONNS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			maxOpen = n
		}
	}
	maxIdle := 5
	if v := os.Getenv("DB_MAX_IDLE_CONNS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			maxIdle = n
		}
	}
	pool.SetMaxOpenConns(maxOpen)
	pool.SetMaxIdleConns(maxIdle)
	if v := os.Getenv("DB_CONN_MAX_LIFETIME_MINUTES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			pool.SetConnMaxLifetime(time.Duration(n) * time.Minute)
		}
	}

	DB = sqlDB
	log.Println("database connection established")
	return nil
}
