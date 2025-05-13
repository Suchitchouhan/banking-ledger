package service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"banking-ledger/internal/constants"
	"banking-ledger/internal/domain"
	"banking-ledger/internal/models"
	"banking-ledger/internal/queue"
)

type TransactionService struct {
	transactionRepo domain.TransactionRepository
	accountRepo     domain.AccountRepository
	producer        queue.Producer
}

func NewTransactionService(
	transactionRepo domain.TransactionRepository,
	accountRepo domain.AccountRepository,
	producer queue.Producer,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		producer:        producer,
	}
}

// creates a new deposit transaction
func (s *TransactionService) CreateDeposit(accountID string, amount float64, description string) (*domain.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("deposit amount must be positive")
	}

	_, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	transaction := &domain.Transaction{
		ID:          uuid.New().String(),
		AccountID:   accountID,
		Type:        constants.TransactionTypeDeposit,
		Amount:      amount,
		Status:      constants.TransactionStatusPending,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.transactionRepo.Create(transaction); err != nil {
		return nil, err
	}

	message := models.TransactionMessage{
		TransactionID: transaction.ID,
		AccountID:     accountID,
		Type:          constants.TransactionTypeDeposit,
		Amount:        amount,
		Description:   description,
	}

	if err := s.producer.PublishTransaction(message); err != nil {
		_ = s.transactionRepo.UpdateStatus(transaction.ID, constants.TransactionStatusFailed)
		return nil, err
	}

	return transaction, nil
}

// Creates a new withdrawal transaction
func (s *TransactionService) CreateWithdrawal(accountID string, amount float64, description string) (*domain.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("withdrawal amount must be positive")
	}

	account, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return nil, err
	}

	if account.Balance < amount {
		return nil, errors.New("insufficient funds")
	}

	now := time.Now()
	transaction := &domain.Transaction{
		ID:          uuid.New().String(),
		AccountID:   accountID,
		Type:        constants.TransactionTypeWithdrawal,
		Amount:      amount,
		Status:      constants.TransactionStatusPending,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.transactionRepo.Create(transaction); err != nil {
		return nil, err
	}

	message := models.TransactionMessage{
		TransactionID: transaction.ID,
		AccountID:     accountID,
		Type:          constants.TransactionTypeWithdrawal,
		Amount:        amount,
		Description:   description,
	}

	if err := s.producer.PublishTransaction(message); err != nil {
		_ = s.transactionRepo.UpdateStatus(transaction.ID, constants.TransactionStatusFailed)
		return nil, err
	}

	return transaction, nil
}

// Retrieves a transaction by ID
func (s *TransactionService) GetTransaction(id string) (*domain.Transaction, error) {
	return s.transactionRepo.GetByID(id)
}

// lists all transactions for an account
func (s *TransactionService) ListTransactionsByAccount(accountID string) ([]*domain.Transaction, error) {
	return s.transactionRepo.ListByAccountID(accountID)
}
