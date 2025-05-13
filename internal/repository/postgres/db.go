package postgres

import (
	"fmt"
	"log"

	"banking-ledger/internal/repository/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(url string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("database migration failed: %v", err)
	}

	return db, nil
}

// DB migrations
func runMigrations(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Account{}); err != nil {
		return fmt.Errorf("failed to migrate accounts table: %v", err)
	}

	return nil
}

// Close the database connection
func Close(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Error getting underlying SQL DB: %v", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}
}
