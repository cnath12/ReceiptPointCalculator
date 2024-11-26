# Receipt Point Calculator

A Go service that processes receipts and calculates points based on rules.

## Prerequisites

- Go 1.21 or higher
- Git

## Setup & Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/cnath12/ReceiptPointCalculator.git
   cd ReceiptPointCalculator
   ```

2. Install dependencies:
   ```bash
   go mod init ReceiptPointCalculator
   go get -u github.com/go-chi/chi/v5
   go get -u github.com/google/uuid
   go get -u github.com/go-playground/validator/v10
   ```

3. Start the server:
   ```bash
   go run cmd/server/main.go
   ```

## Testing

Run all tests:
```bash
go test ./... -v
```

Test specific packages:
```bash
# Test service package
go test ./internal/domain/service -v

# Test handler package
go test ./internal/api/handler -v

# Test validator package
go test ./internal/validator -v

# Test memory storage
go test ./internal/storage/memory -v
```

Run integration tests:
```bash
go test ./internal/tests -v
```