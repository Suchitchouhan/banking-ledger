package domain

import (
	"time"
)

type Account struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AccountRepository interface {
	Create(account *Account) error
	GetByID(id string) (*Account, error)
	UpdateBalance(id string, newBalance float64) error
	List() ([]*Account, error)
}
