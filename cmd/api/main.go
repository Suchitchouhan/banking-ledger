package main

import (
	"fmt"
	"log"

	"banking-ledger/internal/api"
	"banking-ledger/internal/config"
	"banking-ledger/internal/queue"
	"banking-ledger/internal/repository/mongodb"
	"banking-ledger/internal/repository/postgres"
	"banking-ledger/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	fmt.Println(cfg)
	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Connecting to PostgreSQL
	postgresDB, err := postgres.NewConnection(cfg.PostgresURL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer postgres.Close(postgresDB)

	// Connecting to MongoDB
	mongoDB, err := mongodb.NewConnection(cfg.MongoURL, cfg.MongoDB)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Disconnect()

	// RabbitMQ producer
	producer, err := queue.NewRabbitMQProducer(cfg.RabbitMQURL, cfg.TransactionQueue)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ producer: %v", err)
	}
	defer producer.Close()

	accountRepo := postgres.NewAccountRepository(postgresDB)
	transactionRepo := mongodb.NewTransactionRepository(mongoDB)

	accountService := service.NewAccountService(accountRepo)
	transactionService := service.NewTransactionService(
		transactionRepo,
		accountRepo,
		producer,
	)

	handler := api.NewHandler(accountService, transactionService)

	router := handler.CreateRouter()

	if err := api.StartServer(router, cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
