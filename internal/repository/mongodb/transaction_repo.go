package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"banking-ledger/internal/constants"
	"banking-ledger/internal/domain"
)

type TransactionRepository struct {
	collection *mongo.Collection
}

func NewTransactionRepository(conn *Connection) *TransactionRepository {
	return &TransactionRepository{
		collection: conn.Database.Collection("transactions"),
	}
}

// Inserts a new transaction into the database
func (r *TransactionRepository) Create(transaction *domain.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, transaction)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %v", err)
	}
	return nil
}

// Retrieves a transaction by its ID
func (r *TransactionRepository) GetByID(id string) (*domain.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var transaction domain.Transaction
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&transaction)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to retrieve transaction: %v", err)
	}
	return &transaction, nil
}

// Retrieves all transactions for a specific account
func (r *TransactionRepository) ListByAccountID(accountID string) ([]*domain.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{
		"account_id": accountID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions: %v", err)
	}
	defer cursor.Close(ctx)

	var transactions []*domain.Transaction
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, fmt.Errorf("failed to decode transactions: %v", err)
	}

	return transactions, nil
}

// Updates the status of a transaction
func (r *TransactionRepository) UpdateStatus(id string, status constants.TransactionStatus) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"id": id}, update)
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("transaction not found")
	}

	return nil
}
