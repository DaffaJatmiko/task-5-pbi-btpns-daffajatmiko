package database

import (
	"fmt"
	"log"
	"os"
	"task-5-pbi-btpns-daffajatmiko/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

// ConnectDB initializes the database connection and performs migrations.
func ConnectDB() *gorm.DB {
	godotenv.Load(".env")

	DB_HOST := os.Getenv("DB_HOST")
	DB_USER := os.Getenv("DB_USER")
	DB_DRIVER := os.Getenv("DB_DRIVER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")
	DB_PORT := os.Getenv("DB_PORT")

	// Construct the Data Source Name (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
	db, err := gorm.Open("mysql", dsn) // Connect to the MySQL database

	if err != nil {
		fmt.Printf("Cannot connect to %s database", DB_DRIVER)
		log.Fatal(err)
	}

	// Perform auto migrations to create or update database tables
	err = db.Debug().AutoMigrate(&models.User{}, &models.Photo{}).Error
	if err != nil {
		log.Fatalf("Migrating table error: %v", err)
	}

	// Add foreign key constraint for Photo model
	err = db.Debug().Model(&models.Photo{}).AddForeignKey("user_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("Error while attaching foreign key: %v", err)
	}

	return db
}
