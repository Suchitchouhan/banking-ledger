package models

import "time"

type Account struct {
	ID        string    `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	Balance   float64   `gorm:"type:decimal(20,2);default:0.00;not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
