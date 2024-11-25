package validator

import (
    "testing"
    "ReceiptPointCalculator/internal/domain/model"
)

func TestReceiptValidator(t *testing.T) {
    v := NewReceiptValidator()

    tests := []struct {
        name      string
        receipt   model.Receipt
        wantError bool
    }{
        {
            name: "Valid Receipt",
            receipt: model.Receipt{
                Retailer:     "Target",
                PurchaseDate: "2022-01-01",
                PurchaseTime: "13:01",
                Items: []model.Item{
                    {ShortDescription: "Test Item", Price: "12.34"},
                },
                Total: "12.34",
            },
            wantError: false,
        },
        {
            name: "Invalid Retailer",
            receipt: model.Receipt{
                Retailer:     "Target!!!",
                PurchaseDate: "2022-01-01",
                PurchaseTime: "13:01",
                Items: []model.Item{
                    {ShortDescription: "Test Item", Price: "12.34"},
                },
                Total: "12.34",
            },
            wantError: true,
        },
        {
            name: "Invalid Date",
            receipt: model.Receipt{
                Retailer:     "Target",
                PurchaseDate: "2022-13-01",
                PurchaseTime: "13:01",
                Items: []model.Item{
                    {ShortDescription: "Test Item", Price: "12.34"},
                },
                Total: "12.34",
            },
            wantError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := v.ValidateReceipt(&tt.receipt)
            if (err != nil) != tt.wantError {
                t.Errorf("ValidateReceipt() error = %v, wantError %v", err, tt.wantError)
            }
        })
    }
}