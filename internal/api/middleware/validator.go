package middleware

import (
	"bytes"
	"io"
    "context"
    "encoding/json"
    "net/http"
    "ReceiptPointCalculator/internal/domain/model"
    "ReceiptPointCalculator/internal/validator"
    validatorv10 "github.com/go-playground/validator/v10"  
)

type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

type ErrorResponse struct {
    Status  int               `json:"status"`
    Message string            `json:"message"`
    Errors  []ValidationError `json:"errors,omitempty"`
}

func ValidateRequest(v *validator.ReceiptValidator) func(next http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if r.Method == http.MethodPost {
                // Read the body
                body, err := io.ReadAll(r.Body)
                if err != nil {
                    respondWithError(w, http.StatusBadRequest, "Failed to read request body")
                    return
                }
                
                // Create new reader with the same body
                r.Body = io.NopCloser(bytes.NewBuffer(body))

                var receipt model.Receipt
                if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&receipt); err != nil {
                    respondWithError(w, http.StatusBadRequest, "Invalid request body")
                    return
                }

                if err := v.ValidateReceipt(&receipt); err != nil {
                    validationErrors := []ValidationError{}
                    if validationErrs, ok := err.(validatorv10.ValidationErrors); ok {
                        for _, err := range validationErrs {
                            validationErrors = append(validationErrors, ValidationError{
                                Field:   err.Field(),
                                Message: getErrorMessage(err),
                            })
                        }
                    }

                    response := ErrorResponse{
                        Status:  http.StatusBadRequest,
                        Message: "Validation failed",
                        Errors:  validationErrors,
                    }

                    w.Header().Set("Content-Type", "application/json")
                    w.WriteHeader(http.StatusBadRequest)
                    json.NewEncoder(w).Encode(response)
                    return
                }

                ctx := context.WithValue(r.Context(), "receipt", receipt)
                next.ServeHTTP(w, r.WithContext(ctx))
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

func getErrorMessage(err validatorv10.FieldError) string {
    switch err.Tag() {
    case "required":
        return "This field is required"
    case "retailer":
        return "Retailer name must contain only alphanumeric characters, spaces, '-', and '&'"
    case "date":
        return "Must be in format YYYY-MM-DD"
    case "time":
        return "Must be in format HH:MM (24-hour)"
    case "price":
        return "Must be in format XX.XX"
    default:
        return "Invalid value"
    }
}

func respondWithError(w http.ResponseWriter, code int, message string) {
    response := map[string]string{"error": message}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(response)
}