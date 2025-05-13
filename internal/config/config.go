package config

import (
	"fmt"
	"os"
)

type Config struct {
	PostgresURL      string
	MongoURL         string
	MongoDB          string
	RabbitMQURL      string
	TransactionQueue string
	Port             string
}

// Load configuration from environment variables
func Load() (*Config, error) {
	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		postgresURL = "postgres://banking_user:banking_password@localhost:5432/banking_ledger"
	}

	mongoURL := os.Getenv("MONGO_URL")
	if mongoURL == "" {
		mongoURL = "mongodb://banking_user:banking_password@localhost:27017"
	}
	mongoDB := os.Getenv("MONGO_DB")
	if mongoDB == "" {
		mongoDB = "banking_ledger"
	}

	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}
	transactionQueue := os.Getenv("TRANSACTION_QUEUE")
	if transactionQueue == "" {
		transactionQueue = "transaction_queue"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	cfg := &Config{
		PostgresURL:      postgresURL,
		MongoURL:         mongoURL,
		MongoDB:          mongoDB,
		RabbitMQURL:      rabbitMQURL,
		TransactionQueue: transactionQueue,
		Port:             port,
	}

	return cfg, nil
}

// Configuration Validation
func (c *Config) Validate() error {
	if c.PostgresURL == "" {
		return fmt.Errorf("postgres URL is required")
	}
	if c.MongoURL == "" {
		return fmt.Errorf("mongo URL is required")
	}
	if c.RabbitMQURL == "" {
		return fmt.Errorf("RabbitMQ URL is required")
	}
	if c.TransactionQueue == "" {
		return fmt.Errorf("transaction queue name is required")
	}
	return nil
}
