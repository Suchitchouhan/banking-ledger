package domain

import (
	"banking-ledger/internal/constants"
	"time"
)

type Transaction struct {
	ID          string                      `json:"id"`
	AccountID   string                      `json:"account_id"`
	Type        constants.TransactionType   `json:"type"`
	Amount      float64                     `json:"amount"`
	Status      constants.TransactionStatus `json:"status"`
	Description string                      `json:"description"`
	CreatedAt   time.Time                   `json:"created_at"`
	UpdatedAt   time.Time                   `json:"updated_at"`
}

type TransactionRepository interface {
	Create(transaction *Transaction) error
	GetByID(id string) (*Transaction, error)
	ListByAccountID(accountID string) ([]*Transaction, error)
	UpdateStatus(id string, status constants.TransactionStatus) error
}
