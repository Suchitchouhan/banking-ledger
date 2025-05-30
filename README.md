# Banking Ledger Service

## Overview

The Banking Ledger Service is a scalable and reliable system built with Golang. It is designed to efficiently manage bank accounts and transactions, even under high load.

## Architecture

The service follows a microservices-inspired architecture with the following components:

- **API Service**: Handles HTTP requests for managing accounts and transactions.
- **Transaction Processor**: Processes transactions asynchronously.
- **Database Layer**:
    - **PostgreSQL**: Stores account balances and basic account details.
    - **MongoDB**: Maintains detailed transaction logs.
- **Message Queue**: Uses RabbitMQ for reliable transaction processing.

## Key Features

- Scalable to handle high transaction volumes.
- Asynchronous transaction processing for better performance.
- Ensures transaction consistency with ACID-like principles.
- Clear separation of responsibilities across components.
- Robust error handling for reliability.

## System Design Principles

- **Single Responsibility**: Each component focuses on a specific task.
- **Loose Coupling**: Services communicate via message queues for flexibility.
- **Scalability**: Designed to grow with increasing demand.
- **Reliability**: Maintains data integrity and consistency.

## Prerequisites

To run the service, ensure you have the following installed:

- Docker
- Docker Compose
- Go (version 1.23)

## Setup and Running

### Local Development

1. Clone the repository:
     ```bash
     git clone https://github.com/Suchitchouhan/banking-ledger.git
     cd banking-ledger
     ```

2. Build and start the services:
     ```bash
        make start-db
        make run
     ```

### API Endpoints

#### Accounts
- **Create Account**: `POST /accounts`
    - Request body:
        ```json
        { 
            "name": "suchit chouhan", 
            "initial_amount": 1000.00 
        }
        ```

- **List Accounts**: `GET /accounts`

- **Get Account Details**: `GET /accounts/{id}`

#### Transactions
- **Deposit Funds**: `POST /accounts/{id}/deposit`
    - Request body:
        ```json
        { 
            "amount": 500.00, 
            "description": "Salary" 
        }
        ```

- **Withdraw Funds**: `POST /accounts/{id}/withdraw`
    - Request body:
        ```json
        { 
            "amount": 200.00, 
            "description": "Groceries" 
        }
        ```

## Transaction Flow

1. The client sends a transaction request to the API.
2. The API validates the request.
3. The transaction is logged and sent to RabbitMQ.
4. The Transaction Processor picks up the message.
5. The processor updates the account balance in PostgreSQL.
6. The transaction log is stored in MongoDB.

## Testing

Run the tests using the following command:
```bash
make test
```

## Coverage

Run the Coverage using the following command:
make coverage
make cover-html


## Performance Considerations

- Designed for horizontal scalability using a microservices architecture.
- Asynchronous processing with RabbitMQ improves throughput.
- Separate databases optimize data access patterns.

## Potential Improvements

- Add authentication and authorization for security.
- Enhance logging for better debugging.
- Integrate detailed monitoring and metrics.
- Develop a frontend application for user interaction.
- Support advanced transaction types.

## Troubleshooting

- Ensure all services are running.
- Check logs for error messages.
- Verify network connectivity between services.

## Contact

For any questions or issues, feel free to reach out at:  
**Email**: suchitchouhan@outlook.com
