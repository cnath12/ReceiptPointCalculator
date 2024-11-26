// internal/tests/integration_test.go
package integration

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/go-chi/chi/v5"
    "ReceiptPointCalculator/internal/api/handler"
    "ReceiptPointCalculator/internal/domain/model"
    "ReceiptPointCalculator/internal/domain/service"
    "ReceiptPointCalculator/internal/storage/memory"
    "ReceiptPointCalculator/internal/validator"
)

func setupTestServer() *chi.Mux {
    repo := memory.NewReceiptRepository()
    v := validator.NewReceiptValidator()
    svc := service.NewReceiptService(repo)
    handler := handler.NewReceiptHandler(svc, v)

    r := chi.NewRouter()
    r.Post("/receipts/process", handler.ProcessReceipt)
    r.Get("/receipts/{id}/points", handler.GetPoints)
    return r
}

func TestFullReceiptFlow(t *testing.T) {
    r := setupTestServer()
    
    // Test cases from the requirements
    tests := []struct {
        name          string
        receipt      model.Receipt
        wantPoints   int64
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
            wantPoints: 38,
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
            wantPoints: 109,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Step 1: Process Receipt
            receiptJSON, _ := json.Marshal(tt.receipt)
            req := httptest.NewRequest("POST", "/receipts/process", bytes.NewBuffer(receiptJSON))
            req.Header.Set("Content-Type", "application/json")
            rr := httptest.NewRecorder()
            r.ServeHTTP(rr, req)

            if rr.Code != http.StatusOK {
                t.Errorf("Process receipt failed: got status %d", rr.Code)
            }

            var processResp map[string]string
            json.NewDecoder(rr.Body).Decode(&processResp)
            receiptID := processResp["id"]

            // Step 2: Get Points
            req = httptest.NewRequest("GET", "/receipts/"+receiptID+"/points", nil)
            rr = httptest.NewRecorder()
            r.ServeHTTP(rr, req)

            if rr.Code != http.StatusOK {
                t.Errorf("Get points failed: got status %d", rr.Code)
            }

            var pointsResp map[string]int64
            json.NewDecoder(rr.Body).Decode(&pointsResp)

            if pointsResp["points"] != tt.wantPoints {
                t.Errorf("Points = %v, want %v", pointsResp["points"], tt.wantPoints)
            }
        })
    }
}

func TestInvalidRequests(t *testing.T) {
    r := setupTestServer()
    
    tests := []struct {
        name       string
        receipt    interface{}
        wantStatus int
    }{
        {
            name: "Missing Required Fields",
            receipt: map[string]interface{}{
                "retailer": "Target",
                // missing other required fields
            },
            wantStatus: http.StatusBadRequest,
        },
        {
            name: "Invalid Date Format",
            receipt: model.Receipt{
                Retailer:     "Target",
                PurchaseDate: "2022-13-01", // invalid month
                PurchaseTime: "13:01",
                Items: []model.Item{{ShortDescription: "Test", Price: "1.00"}},
                Total: "1.00",
            },
            wantStatus: http.StatusBadRequest,
        },
        {
            name: "Invalid Time Format",
            receipt: model.Receipt{
                Retailer:     "Target",
                PurchaseDate: "2022-01-01",
                PurchaseTime: "25:01", // invalid hour
                Items: []model.Item{{ShortDescription: "Test", Price: "1.00"}},
                Total: "1.00",
            },
            wantStatus: http.StatusBadRequest,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            receiptJSON, _ := json.Marshal(tt.receipt)
            req := httptest.NewRequest("POST", "/receipts/process", bytes.NewBuffer(receiptJSON))
            req.Header.Set("Content-Type", "application/json")
            rr := httptest.NewRecorder()
            r.ServeHTTP(rr, req)

            if rr.Code != tt.wantStatus {
                t.Errorf("got status %d, want %d", rr.Code, tt.wantStatus)
            }
        })
    }
}