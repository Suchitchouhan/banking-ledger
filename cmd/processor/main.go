package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"banking-ledger/internal/config"
	"banking-ledger/internal/models"
	"banking-ledger/internal/processor"
	"banking-ledger/internal/queue"
	"banking-ledger/internal/repository/mongodb"
	"banking-ledger/internal/repository/postgres"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to PostgreSQL
	postgresDB, err := postgres.NewConnection(cfg.PostgresURL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer postgres.Close(postgresDB)

	// Connect to MongoDB
	mongoDB, err := mongodb.NewConnection(cfg.MongoURL, cfg.MongoDB)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Disconnect()

	// Create repositories
	accountRepo := postgres.NewAccountRepository(postgresDB)
	transactionRepo := mongodb.NewTransactionRepository(mongoDB)

	// Create consumer
	consumer, err := queue.NewRabbitMQConsumer(cfg.RabbitMQURL, cfg.TransactionQueue)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ consumer: %v", err)
	}
	defer consumer.Close()

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Process transactions
	processor := processor.NewTransactionProcessor(accountRepo, transactionRepo)
	err = consumer.Consume(ctx, func(data []byte) error {
		var msg models.TransactionMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return err
		}

		return processor.ProcessTransaction(ctx, msg)
	})
	if err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down transaction processor...")
	cancel()
}
