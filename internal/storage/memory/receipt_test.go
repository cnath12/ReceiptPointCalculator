package memory

import (
    "context"
    "testing"
    "ReceiptPointCalculator/internal/domain/model"
)

func TestReceiptRepository(t *testing.T) {
    repo := NewReceiptRepository()
    ctx := context.Background()

    receipt := &model.Receipt{
        ID:       "test-id",
        Retailer: "Test Store",
    }

    t.Run("Save and Get Receipt", func(t *testing.T) {
        err := repo.Save(ctx, receipt)
        if err != nil {
            t.Errorf("Save() error = %v", err)
        }

        got, err := repo.GetByID(ctx, receipt.ID)
        if err != nil {
            t.Errorf("GetByID() error = %v", err)
        }
        if got.ID != receipt.ID {
            t.Errorf("GetByID() = %v, want %v", got.ID, receipt.ID)
        }
    })

    t.Run("Get Non-existent Receipt", func(t *testing.T) {
        _, err := repo.GetByID(ctx, "non-existent")
        if err != ErrNotFound {
            t.Errorf("GetByID() error = %v, want %v", err, ErrNotFound)
        }
    })
}