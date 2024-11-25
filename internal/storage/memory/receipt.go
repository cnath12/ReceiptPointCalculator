package memory

import (
    "context"
    "sync"
    "ReceiptPointCalculator/internal/domain/model"
	"errors"
)
var ErrNotFound = errors.New("receipt not found")

type ReceiptRepository struct {
    mu      sync.RWMutex
    storage map[string]*model.Receipt
}

func NewReceiptRepository() *ReceiptRepository {
    return &ReceiptRepository{
        storage: make(map[string]*model.Receipt),
    }
}

func (r *ReceiptRepository) Save(_ context.Context, receipt *model.Receipt) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    r.storage[receipt.ID] = receipt
    return nil
}

func (r *ReceiptRepository) GetByID(_ context.Context, id string) (*model.Receipt, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    receipt, ok := r.storage[id]
    if !ok {
        return nil, ErrNotFound
    }
    return receipt, nil
}
