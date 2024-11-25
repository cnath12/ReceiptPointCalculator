package handler

import (
    "encoding/json"
    "net/http"
    "github.com/go-chi/chi/v5"
    "ReceiptPointCalculator/internal/domain/model"
	"context"
	"log"
	"ReceiptPointCalculator/internal/validator"
)

type ReceiptService interface {
    ProcessReceipt(ctx context.Context, receipt *model.Receipt) (string, error)
    GetPoints(ctx context.Context, id string) (int64, error)
}

type ReceiptHandler struct {
    service ReceiptService
	validator *validator.ReceiptValidator
}

func NewReceiptHandler(service ReceiptService, validator *validator.ReceiptValidator) *ReceiptHandler {
    return &ReceiptHandler{
        service:   service,
        validator: validator,
    }
}


func (h *ReceiptHandler) ProcessReceipt(w http.ResponseWriter, r *http.Request) {
    var receipt model.Receipt
    if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
        log.Printf("Error decoding request: %v", err) // Add this line
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    // Add validation logging
    if err := h.validator.ValidateReceipt(&receipt); err != nil {
        log.Printf("Validation error: %v", err) // Add this line
        respondWithError(w, http.StatusBadRequest, err.Error())
        return
    }

    id, err := h.service.ProcessReceipt(r.Context(), &receipt)
    if err != nil {
        log.Printf("Error processing receipt: %v", err) // Add this line
        respondWithError(w, http.StatusInternalServerError, "Failed to process receipt")
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]string{"id": id})
}

func (h *ReceiptHandler) GetPoints(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    
    points, err := h.service.GetPoints(r.Context(), id)
    if err != nil {
        respondWithError(w, http.StatusNotFound, "Receipt not found")
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]int64{"points": points})
}
