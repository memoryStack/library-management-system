package initializers

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"library-management-system/backend/models"
)

var DB *gorm.DB

// ConnectDB opens Postgres using DATABASE_URL from the loaded env file.
func ConnectDB(environment string) error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL must be set in .env.%s", environment)
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

	pool.SetMaxOpenConns(25)
	pool.SetMaxIdleConns(5)
	pool.SetConnMaxLifetime(30 * time.Minute)

	DB = sqlDB
	log.Println("database connection established")
	return nil
}

func SyncDB() error {
	// to create/sync the tables in the database
	return DB.AutoMigrate(&models.Book{})
}
