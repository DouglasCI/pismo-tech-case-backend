# Pismo Backend Technical Challenge

A REST API built in Go for managing bank accounts and financial transactions. This service handles a simple account and transaction processing.

## Architecture & Design Decisions

### Project Structure

To focus on simplicity, the code was structured in two folders and a flat hierarchy of files:

```text
.
├── cmd/
│   └── api/                 # Entry point (DB connection and server startup)
│       └── main.go          
└── internal/
    └── server/              # The core of the application
        ├── handlers.go      # Encapsulates handlers for each HTTP endpoint
        ├── models.go        # Contains models definitions and validation rules
        └── schemas.go       # Defines schemas and seed data for the database tables
```

#### Architectural Decision: Pragmatism vs. Over-engineering

In the context of a scoped microservice with only two core entities (Accounts and Transactions), adopting a rigid Clean Architecture with multiple layers of interfaces, use cases, and repositories would introduce unnecessary boilerplate and violate the Go philosophy of keeping things simple and readable.

The current "flat" structure couples data access directly within the handlers to maximize development speed and code clarity. In a real-world, large-scale production environment where database portability or strict decoupling becomes a necessity, the next step would be to extract the data access layer using the **Repository Pattern** and **Dependency Inversion**. For this challenge, I prioritized a functional, testable, and pragmatic approach over premature abstraction.

### Database

SQLite was chosen for this project because it is exceptionally lightweight and easy to deploy, making it the ideal database for when simplicity and portability are priorities.

To maximize these benefits, it uses a CGO-free driver (modernc.org/sqlite). This allows the application to be compiled into a static binary without external C dependencies, ensuring the final Docker image remains as small as possible. Additionally, the `sqlx` library is used to keep database interactions clean and highly readable.

### Database Normalization: Operation Types

A architectural decision was to represent operation types as a dedicated database table rather than relying on application-level constants.

In a real-world financial system, the rules governing how transactions are classified often change or expand. By moving these to the database, we achieve:

- Referential Integrity: Every transaction is tied to a valid operation type via a Foreign Key. This makes it impossible for the database to accept a transaction with an "unknown" operation ID, providing a second layer of defense beyond application validation.

- Extensibility: New types of operations (e.g., "Cashback," "Reversal," "Interest") can be added via a simple SQL migration or seed script without requiring a recompilation or redeployment of the Go binary.

- Auditability: Standardizes descriptions across the system. This allows for easier reporting and joins when building dashboards or financial statements, as the "source of truth" for what an ID means lives right next to the data.

#### Design Trade-off: The "Join Tax"

By separating `operation_types` into a lookup table, retrieving the `description` of a `operation_type` in a transaction requires a JOIN. In high-volume systems (millions of reads per second), this can introduce overhead compared to a denormalized "flat" table.

### Tech Stack

- Language: Go 1.25+
- Database: SQLite (SQLX for orchestration)
- Containerization: Docker (Multi-stage builds)
- Testing: Native Go testing package.

## Getting Started

### Prerequisites

- Docker & Docker Compose (Recommended)
- Go 1.25 (for native execution)

### Option 1: Running with Docker (Recommended)

The easiest way to get the service up and running is using the provided shell script:

```bash
chmod +x run.sh
./run.sh
```

The server will be available at <http://localhost:8080>.

### Option 2: Running Natively

```bash
go mod tidy
go run ./cmd/api
```

The server will be available at <http://localhost:8080>.

## Testing & Coverage

The suite includes unit tests for business rules and integration tests for API handlers using in-memory SQLite.

### Run all tests

```bash
go test ./internal/server -v
```

### Check coverage percentage

```bash
go test ./internal/server -cover
```

## API Endpoints

### Accounts

- `POST /accounts`: Create a new account.
- Payload:

  ```json
  {
    "document_number": "12345678900"
  }
  ```

- `GET /accounts/{accountId}`: Retrieve account details.
- Response:

  ```json
  {
    "account_id": 1,
    "document_number": "123456"
  }
  ```

### Transactions

- `POST /transactions`: Create a new transaction.
- Payload:

  ```json
    {
      "account_id": 1,
      "operation_type_id": 2,
      "amount": -10.45
    } 
  ```

- Operation Types reference table.

  | ID | Description                | Amount Sign  |
  |----|----------------------------|--------------|
  | 1  | NORMAL PURCHASE            | Negative (-) |
  | 2  | PURCHASE WITH INSTALLMENTS | Negative (-) |
  | 3  | WITHDRAWAL                 | Negative (-) |
  | 4  | CREDIT VOUCHER             | Positive (+) |
