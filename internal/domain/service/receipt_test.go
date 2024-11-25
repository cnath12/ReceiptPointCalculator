package service

import (
    "testing"
    "ReceiptPointCalculator/internal/domain/model"
)

func TestCalculatePoints(t *testing.T) {
    tests := []struct {
        name     string
        receipt  model.Receipt
        expected int64
    }{
        {
            name: "Target Receipt",
            receipt: model.Receipt{
                Retailer:     "Target",
                PurchaseDate: "2022-01-01",
                PurchaseTime: "13:01",
                Items: []model.Item{
                    {ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
                    {ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
                    {ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
                    {ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
                    {ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
                },
                Total: "35.35",
            },
            expected: 38,
        },
        {
            name: "M&M Corner Market Receipt",
            receipt: model.Receipt{
                Retailer:     "M&M Corner Market",
                PurchaseDate: "2022-03-20",
                PurchaseTime: "14:33",
                Items: []model.Item{
                    {ShortDescription: "Gatorade", Price: "2.25"},
                    {ShortDescription: "Gatorade", Price: "2.25"},
                    {ShortDescription: "Gatorade", Price: "2.25"},
                    {ShortDescription: "Gatorade", Price: "2.25"},
                },
                Total: "9.00",
            },
            expected: 109,
        },
    }

    svc := NewReceiptService(nil)
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            points := svc.CalculatePoints(&tt.receipt)
            if points != tt.expected {
                t.Errorf("CalculatePoints() = %v, want %v", points, tt.expected)
            }
        })
    }
}