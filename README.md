# Pismo Backend Technical Challenge

A REST API built in Go for managing bank accounts and financial transactions. This service handles a simple account and transaction processing.

## Architecture & Design Decisions

### Project Structure

To balance simplicity with a clear separation of concerns, the project isolates business entities from the HTTP transport layer. This prevents business rules from being tangled with web request parsing.

```text
.
├── cmd/
│   └── api/                 # Entry point (DB connection, wiring, and startup)
│       └── main.go          
└── internal/
    ├── domain/              # Core business logic and entities
    │   ├── account.go       # Account entity and validation
    │   ├── money.go         # Custom Money type to handle amount (cents pattern)
    │   └── transaction.go   # Transaction entity and operation rules
    └── server/              # HTTP transport and database infrastructure
        ├── handlers.go      # Encapsulates HTTP handlers and request/response
        └── schemas.go       # Defines schemas and seed data for the database tables
```

### Architectural Decision: Pragmatic Separation of Concerns

In the context of a scoped microservice with only two core entities (Accounts and Transactions), adopting a rigid Clean Architecture (with multiple layers of interfaces, services, and repositories) would introduce unnecessary boilerplate and violate the Go philosophy of keeping things simple and readable.

Instead, I opted for a **pragmatic middle-ground**. The core business rules, validations, and financial precision logic are strictly isolated within the `domain` package, ensuring they remain pure and highly testable. However, to maximize development speed and code clarity, the HTTP transport and data access (SQL queries) remain coupled within the `server` handlers.

In a real-world, large-scale production environment where database portability or strict decoupling becomes a necessity, the natural next step would be to extract the data access layer using the **Repository Pattern** and **Dependency Inversion**. For this challenge, I prioritized a functional approach over premature abstraction.

### Database

SQLite was chosen for this project because it is exceptionally lightweight and easy to deploy, making it the ideal database for when simplicity and portability are priorities.

To maximize these benefits, it uses a CGO-free driver (modernc.org/sqlite). This allows the application to be compiled into a static binary without external C dependencies, ensuring the final Docker image remains as small as possible. Additionally, the `sqlx` library is used to keep database interactions clean and highly readable.

### Database Normalization: Operation Types

An architectural decision was to represent operation types as a dedicated database table rather than relying on application-level constants.

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
go test ./internal/... -v
```

### Check coverage percentage

```bash
go test ./internal/... -cover
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
    "document_number": "12345678900"
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

## Next Steps & Future Improvements

While this project was designed to be simple and fulfill the core requirements of the challenge, a production-ready evolution of this microservice would include the following improvements:

### 1. Architectural Evolution (Clean Architecture)

As the business rules grow, the current architecture should be refactored to a **Layered/Clean Architecture** (Domain, Service, Repository, Delivery). This would decouple the HTTP transport layer and the database driver from the business logic, heavily utilizing Dependency Inversion to make the system easier to test and maintain.

### 2. Database Migration & Containerization

- **PostgreSQL:** Migrate from SQLite to a robust relational database like PostgreSQL.
- **Migration Tool:** Introduce a formal migration tool (like `golang-migrate`) instead of running raw schema queries on startup.

### 3. Testing Strategy

- Add more robust **Integration Tests** for the HTTP handlers and database queries focusing on error handling.

### 4. Resiliency

- **Idempotency:** Implement idempotency keys for the `POST /transactions` endpoint to prevent double-spending or duplicate transactions in case of network retries.
- **Structured Logging & Observability:** Replace standard library logs with structured JSON logging (e.g., `slog` or `zap`) and add correlation IDs to trace requests across the system.

### 5. API Documentation

- Implement **Swagger/OpenAPI 3.0** specifications to automatically generate interactive documentation for the endpoints and payloads.

### 6. Evolving the Money Domain

- **Financial Operations (Value Object):** Expand the custom `Money` type into a fully-featured DDD Value Object. This includes implementing safe mathematical methods (e.g., `Add()`, `Subtract()`, `Percentage()`) directly on the type. This ensures financial calculations are encapsulated, preventing integer overflows and logic duplication across the system.
