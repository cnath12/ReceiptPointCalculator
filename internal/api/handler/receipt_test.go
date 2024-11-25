package handler

import (
    "bytes"
    "context"
    "encoding/json"
    "errors"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/go-chi/chi/v5"
    "ReceiptPointCalculator/internal/domain/model"
    "ReceiptPointCalculator/internal/validator"
)

type mockService struct {
    processFunc func(ctx context.Context, receipt *model.Receipt) (string, error)
    pointsFunc  func(ctx context.Context, id string) (int64, error)
}

func (m *mockService) ProcessReceipt(ctx context.Context, receipt *model.Receipt) (string, error) {
    return m.processFunc(ctx, receipt)
}

func (m *mockService) GetPoints(ctx context.Context, id string) (int64, error) {
    return m.pointsFunc(ctx, id)
}

var ErrNotFound = errors.New("receipt not found")

func TestProcessReceipt(t *testing.T) {
    tests := []struct {
        name           string
        receipt       model.Receipt
        mockID        string
        mockErr       error
        expectedCode  int
    }{
        {
            name: "Success",
            receipt: model.Receipt{
                Retailer:     "Target",
                PurchaseDate: "2022-01-01", 
                PurchaseTime: "13:01",        
                Items: []model.Item{
                    {
                        ShortDescription: "Test Item",
                        Price:           "1.00",
                    },
                },
                Total: "1.00",
            },
            mockID:       "test-id",
            mockErr:      nil,
            expectedCode: http.StatusOK,
        },
        {
            name: "Invalid Receipt - Missing Required Fields",
            receipt: model.Receipt{
                Retailer: "Target",
                Total:   "1.00",
            },
            mockID:       "",
            mockErr:      nil,
            expectedCode: http.StatusBadRequest,
        },
        {
            name: "Invalid Receipt - Invalid Date Format",
            receipt: model.Receipt{
                Retailer:     "Target",
                PurchaseDate: "2022-13-01",    // Invalid month
                PurchaseTime: "13:01",
                Items: []model.Item{
                    {
                        ShortDescription: "Test Item",
                        Price:           "1.00",
                    },
                },
                Total: "1.00",
            },
            mockID:       "",
            mockErr:      nil,
            expectedCode: http.StatusBadRequest,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockSvc := &mockService{
                processFunc: func(ctx context.Context, receipt *model.Receipt) (string, error) {
                    return tt.mockID, tt.mockErr
                },
            }
            
            v := validator.NewReceiptValidator()
            handler := NewReceiptHandler(mockSvc, v)

            body, _ := json.Marshal(tt.receipt)
            req := httptest.NewRequest(http.MethodPost, "/receipts/process", bytes.NewBuffer(body))
            w := httptest.NewRecorder()

            handler.ProcessReceipt(w, req)

            if w.Code != tt.expectedCode {
                t.Errorf("ProcessReceipt() status = %v, want %v", w.Code, tt.expectedCode)
            }

            // For successful cases, verify response
            if tt.expectedCode == http.StatusOK {
                var response map[string]string
                json.NewDecoder(w.Body).Decode(&response)
                if response["id"] != tt.mockID {
                    t.Errorf("ProcessReceipt() id = %v, want %v", response["id"], tt.mockID)
                }
            }
        })
    }
}

func TestGetPoints(t *testing.T) {
    tests := []struct {
        name          string
        id            string
        mockPoints    int64
        mockErr       error
        expectedCode  int
    }{
        {
            name:         "Success",
            id:          "test-id",
            mockPoints:  100,
            mockErr:     nil,
            expectedCode: http.StatusOK,
        },
        {
            name:         "Not Found",
            id:          "non-existent",
            mockPoints:  0,
            mockErr:     ErrNotFound,
            expectedCode: http.StatusNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockSvc := &mockService{
                pointsFunc: func(ctx context.Context, id string) (int64, error) {
                    return tt.mockPoints, tt.mockErr
                },
            }

            v := validator.NewReceiptValidator()
            handler := NewReceiptHandler(mockSvc, v)

            req := httptest.NewRequest(http.MethodGet, "/receipts/{id}/points", nil)
            rctx := chi.NewRouteContext()
            rctx.URLParams.Add("id", tt.id)
            req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
            w := httptest.NewRecorder()

            handler.GetPoints(w, req)

            if w.Code != tt.expectedCode {
                t.Errorf("GetPoints() status = %v, want %v", w.Code, tt.expectedCode)
            }
        })
    }
}