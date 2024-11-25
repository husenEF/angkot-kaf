package database

import (
	"log"
	"os"
	"path/filepath"

	"github.com/robzlabz/angkot-kaf/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	// Ensure database directory exists
	dbDir := "database"
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatal("Failed to create database directory:", err)
	}

	dbPath := filepath.Join(dbDir, "angkot.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto Migrate the schema
	err = db.AutoMigrate(&models.Driver{}, &models.Trip{}, &models.Passenger{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	DB = db
}
