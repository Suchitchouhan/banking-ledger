package models

import "banking-ledger/internal/constants"

type CreateAccountRequest struct {
	Name          string  `json:"name" binding:"required"`
	InitialAmount float64 `json:"initial_amount" binding:"min=0"`
}

type TransactionRequest struct {
	Amount      float64 `json:"amount" binding:"required,min=0.01"`
	Description string  `json:"description"`
}

type TransactionMessage struct {
	TransactionID string                    `json:"transaction_id"`
	AccountID     string                    `json:"account_id"`
	Type          constants.TransactionType `json:"type"`
	Amount        float64                   `json:"amount"`
	Description   string                    `json:"description"`
}
