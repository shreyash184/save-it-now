package main

import (
	"expense-tracker/db"
	"expense-tracker/models"
	"fmt"
	"log"
)

func main() {
	// Connect to DB
	db.InitDB()

	// Dummy expense
	exp := models.Expense{
		Amount:   123.45,
		Note:     "Testing DB insert",
		Category: "Testing",
	}

	// Insert into DB
	if err := db.DB.Create(&exp).Error; err != nil {
		log.Fatalf("❌ Failed to insert: %v", err)
	}

	fmt.Println("✅ Expense inserted successfully:", exp)
}
