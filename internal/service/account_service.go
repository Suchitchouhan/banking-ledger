package service

import (
	"errors"
	"time"

	"banking-ledger/internal/domain"

	"github.com/google/uuid"
)

type AccountService struct {
	accountRepo domain.AccountRepository
}

func NewAccountService(accountRepo domain.AccountRepository) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
	}
}

// Creates a new account with initial balance
func (s *AccountService) CreateAccount(name string, initialBalance float64) (*domain.Account, error) {
	if initialBalance < 0 {
		return nil, errors.New("initial balance cannot be negative")
	}

	now := time.Now()
	account := &domain.Account{
		ID:        uuid.New().String(),
		Name:      name,
		Balance:   initialBalance,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.accountRepo.Create(account); err != nil {
		return nil, err
	}
	return account, nil
}

// Retrieves an account by ID
func (s *AccountService) GetAccount(id string) (*domain.Account, error) {
	return s.accountRepo.GetByID(id)
}

// lists all accounts
func (s *AccountService) ListAccounts() ([]*domain.Account, error) {
	return s.accountRepo.List()
}
