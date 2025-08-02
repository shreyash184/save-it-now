package models

import "time"

type Expense struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Amount    float64   `json:"amount"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	Category  string    `json:"category"`
}
