package db

import (
	"expense-tracker/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"fmt"
)

var DB *gorm.DB

func InitDB() {
	dsn := os.Getenv("SUPABASE_DB_URL")
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to Supabase/Postgres database")
	}
	fmt.Println("DSN:", dsn)

	database.AutoMigrate(&models.Expense{})
	DB = database
}
