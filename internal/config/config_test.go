package config

import (
	"os"
	"reflect"
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid configuration",
			config: Config{
				PostgresURL:      "postgres://user:password@localhost:5432/db",
				MongoURL:         "mongodb://user:password@localhost:27017",
				RabbitMQURL:      "amqp://guest:guest@localhost:5672/",
				TransactionQueue: "transaction_queue",
			},
			wantErr: false,
		},
		{
			name: "Missing PostgresURL",
			config: Config{
				MongoURL:         "mongodb://user:password@localhost:27017",
				RabbitMQURL:      "amqp://guest:guest@localhost:5672/",
				TransactionQueue: "transaction_queue",
			},
			wantErr: true,
			errMsg:  "postgres URL is required",
		},
		{
			name: "Missing MongoURL",
			config: Config{
				PostgresURL:      "postgres://user:password@localhost:5432/db",
				RabbitMQURL:      "amqp://guest:guest@localhost:5672/",
				TransactionQueue: "transaction_queue",
			},
			wantErr: true,
			errMsg:  "mongo URL is required",
		},
		{
			name: "Missing RabbitMQURL",
			config: Config{
				PostgresURL:      "postgres://user:password@localhost:5432/db",
				MongoURL:         "mongodb://user:password@localhost:27017",
				TransactionQueue: "transaction_queue",
			},
			wantErr: true,
			errMsg:  "RabbitMQ URL is required",
		},
		{
			name: "Missing TransactionQueue",
			config: Config{
				PostgresURL: "postgres://user:password@localhost:5432/db",
				MongoURL:    "mongodb://user:password@localhost:27017",
				RabbitMQURL: "amqp://guest:guest@localhost:5672/",
			},
			wantErr: true,
			errMsg:  "transaction queue name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name           string
		envVars        map[string]string
		expectedConfig *Config
	}{
		{
			name: "All environment variables set",
			envVars: map[string]string{
				"POSTGRES_URL":      "postgres://user:password@localhost:5432/db",
				"MONGO_URL":         "mongodb://user:password@localhost:27017",
				"MONGO_DB":          "custom_db",
				"RABBITMQ_URL":      "amqp://guest:guest@localhost:5672/",
				"TRANSACTION_QUEUE": "custom_queue",
				"PORT":              "9090",
			},
			expectedConfig: &Config{
				PostgresURL:      "postgres://user:password@localhost:5432/db",
				MongoURL:         "mongodb://user:password@localhost:27017",
				MongoDB:          "custom_db",
				RabbitMQURL:      "amqp://guest:guest@localhost:5672/",
				TransactionQueue: "custom_queue",
				Port:             "9090",
			},
		},
		{
			name:    "No environment variables set, use defaults",
			envVars: map[string]string{},
			expectedConfig: &Config{
				PostgresURL:      "postgres://banking_user:banking_password@localhost:5432/banking_ledger",
				MongoURL:         "mongodb://banking_user:banking_password@localhost:27017",
				MongoDB:          "banking_ledger",
				RabbitMQURL:      "amqp://guest:guest@localhost:5672/",
				TransactionQueue: "transaction_queue",
				Port:             "8080",
			},
		},
		{
			name: "Some environment variables set",
			envVars: map[string]string{
				"POSTGRES_URL": "postgres://user:password@localhost:5432/db",
				"MONGO_DB":     "custom_db",
			},
			expectedConfig: &Config{
				PostgresURL:      "postgres://user:password@localhost:5432/db",
				MongoURL:         "mongodb://banking_user:banking_password@localhost:27017",
				MongoDB:          "custom_db",
				RabbitMQURL:      "amqp://guest:guest@localhost:5672/",
				TransactionQueue: "transaction_queue",
				Port:             "8080",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			defer func() {
				for key := range tt.envVars {
					os.Unsetenv(key)
				}
			}()

			cfg, err := Load()
			if err != nil {
				t.Fatalf("Load() returned an error: %v", err)
			}

			if !reflect.DeepEqual(cfg, tt.expectedConfig) {
				t.Errorf("Load() = %v, want %v", cfg, tt.expectedConfig)
			}
		})
	}
}
