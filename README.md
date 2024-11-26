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
The base url is: http://localhost:8080 

## API Endpoint Usage

1. Process Receipt
```bash
POST http://localhost:8080/receipts/process
```
# Request Body Example:
```bash
{
  "retailer": "M&M Corner Market",
  "purchaseDate": "2022-03-20",
  "purchaseTime": "14:33",
  "items": [
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    }
  ],
  "total": "9.00"
}

# Response:
{
  "id": "7fb1377b-b223-49d9-a31a-5a02701dd310"
}
```

2. Get Points
```bash
GET http://localhost:8080/receipts/{id}/points

# Response:
{
  "points": 109
}
```
## Validation Rules

Receipt Fields

Retailer: Required, alphanumeric characters, spaces, '-', and '&' only
purchaseDate: Required, format YYYY-MM-DD
purchaseTime: Required, 24-hour format HH:MM
items: Required, at least one item
total: Required, format XX.XX (dollars.cents)

Item Fields

shortDescription: Required, alphanumeric characters, spaces, and '-' only
price: Required, format XX.XX (dollars.cents)

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