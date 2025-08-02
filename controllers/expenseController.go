package controllers

import (
	"expense-tracker/db"
	"expense-tracker/models"
	"net/http"
	"time"
	"testing"
	"github.com/gin-gonic/gin"
)

func TestInsertExpense(t *testing.T) {
	exp := models.Expense{
		Amount:   100,
		Note:     "Test expense",
		Category: "Test",
	}

	err := db.DB.Create(&exp).Error
	if err != nil {
		t.Errorf("Failed to insert expense: %v", err)
	}
}

func AddExpense(c *gin.Context) {
	var expense models.Expense

	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set CreatedAt if not set
	if expense.CreatedAt.IsZero() {
		expense.CreatedAt = time.Now()
	}

	if err := db.DB.Create(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expense"})
		return
	}
	// log.Println("Received expense:", expense)

	c.JSON(http.StatusOK, expense)
}

func GetExpenses(c *gin.Context) {
	var expenses []models.Expense
	if err := db.DB.Find(&expenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve expenses"})
		return
	}
	c.JSON(http.StatusOK, expenses)
}
