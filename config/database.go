package config

import (
	"Test_Fleetify/domain/models"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func ConnectDB() {
	_ = godotenv.Load()

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	DB = database
	log.Println("Database connected successfully")

	// Migrate tables in correct order
	err = migrateTables(DB)
	if err != nil {
		log.Fatalf("Failed to migrate tables: %v", err)
	}
	log.Println("Tables migrated successfully")
}

func migrateTables(db *gorm.DB) error {
	// Migrate tables in correct order to avoid foreign key issues
	tables := []interface{}{
		&models.Department{},
		&models.Employee{},
		&models.Attendance{},
		&models.AttendanceHistory{},
	}

	for _, table := range tables {
		err := db.AutoMigrate(table)
		if err != nil {
			return fmt.Errorf("failed to migrate table: %v", err)
		}
	}

	return nil
}
