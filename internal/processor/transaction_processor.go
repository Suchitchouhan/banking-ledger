package processor

import (
	"banking-ledger/internal/constants"
	"banking-ledger/internal/domain"
	"banking-ledger/internal/models"
	"context"
	"log"
)

type TransactionProcessor struct {
	accountRepo     domain.AccountRepository
	transactionRepo domain.TransactionRepository
}

func NewTransactionProcessor(
	accountRepo domain.AccountRepository,
	transactionRepo domain.TransactionRepository,
) *TransactionProcessor {
	return &TransactionProcessor{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
	}
}

func (p *TransactionProcessor) ProcessTransaction(ctx context.Context, msg models.TransactionMessage) error {
	transaction, err := p.transactionRepo.GetByID(msg.TransactionID)
	if err != nil {
		log.Printf("Failed to retrieve transaction %s: %v", msg.TransactionID, err)
		return err
	}

	account, err := p.accountRepo.GetByID(msg.AccountID)
	if err != nil {
		log.Printf("Failed to retrieve account %s: %v", msg.AccountID, err)
		// Mark transaction as failed
		_ = p.transactionRepo.UpdateStatus(transaction.ID, constants.TransactionStatusFailed)
		return err
	}

	// Calculate new balance based on transaction type
	var newBalance float64
	switch transaction.Type {
	case constants.TransactionTypeDeposit:
		newBalance = account.Balance + transaction.Amount
	case constants.TransactionTypeWithdrawal:
		// Double check balance sufficiency
		if account.Balance < transaction.Amount {
			log.Printf("Insufficient funds in account %s for transaction %s", account.ID, transaction.ID)
			// Mark transaction as failed
			_ = p.transactionRepo.UpdateStatus(transaction.ID, constants.TransactionStatusFailed)
			return nil // Don't retry
		}
		newBalance = account.Balance - transaction.Amount
	default:
		log.Printf("Unknown transaction type: %s", transaction.Type)
		// Mark transaction as failed
		_ = p.transactionRepo.UpdateStatus(transaction.ID, constants.TransactionStatusFailed)
		return nil // Don't retry
	}

	// Update account balance
	if err := p.accountRepo.UpdateBalance(account.ID, newBalance); err != nil {
		log.Printf("Failed to update balance for account %s: %v", account.ID, err)
		// Mark transaction as failed
		_ = p.transactionRepo.UpdateStatus(transaction.ID, constants.TransactionStatusFailed)
		return err
	}

	// Mark transaction as completed
	if err := p.transactionRepo.UpdateStatus(transaction.ID, constants.TransactionStatusCompleted); err != nil {
		log.Printf("Failed to update transaction status: %v", err)
		// This is problematic as the balance was updated but the transaction status wasn't
		//  we'd need to handle this inconsistency
		return err
	}

	log.Printf("Successfully processed transaction %s for account %s", transaction.ID, account.ID)
	return nil
}
