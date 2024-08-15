package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "sqlite" // Default to SQLite if not specified
	}

	var dsn string
	var dialector gorm.Dialector

	switch dbType {
	case "postgres":
		dsn = os.Getenv("POSTGRES_DSN")
		if dsn == "" {
			return fmt.Errorf("POSTGRES_DSN environment variable is not set")
		}
		dialector = postgres.Open(dsn)
	case "sqlite":
		dsn = os.Getenv("SQLITE_DSN")
		if dsn == "" {
			dsn = "test.db" // Default SQLite database file
		}
		dialector = sqlite.Open(dsn)
	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}

	var err error
	DB, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	return nil
}
