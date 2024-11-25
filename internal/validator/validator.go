package validator

import (
    "regexp"
    "time"
    "github.com/go-playground/validator/v10"
)

type ReceiptValidator struct {
    validate *validator.Validate
}

func NewReceiptValidator() *ReceiptValidator {
    v := validator.New()
    
    // Register custom validation functions
    v.RegisterValidation("retailer", validateRetailer)
    v.RegisterValidation("date", validateDate)
    v.RegisterValidation("time", validateTime)
    v.RegisterValidation("price", validatePrice)
    
    return &ReceiptValidator{validate: v}
}

func validateRetailer(fl validator.FieldLevel) bool {
    retailer := fl.Field().String()
    matched, _ := regexp.MatchString(`^[\w\s\-&]+$`, retailer)
    return matched
}

func validateDate(fl validator.FieldLevel) bool {
    date := fl.Field().String()
    _, err := time.Parse("2006-01-02", date)
    return err == nil
}

func validateTime(fl validator.FieldLevel) bool {
    timeStr := fl.Field().String()
    _, err := time.Parse("15:04", timeStr)
    return err == nil
}

func validatePrice(fl validator.FieldLevel) bool {
    price := fl.Field().String()
    matched, _ := regexp.MatchString(`^\d+\.\d{2}$`, price)
    return matched
}

func (v *ReceiptValidator) ValidateReceipt(receipt interface{}) error {
    return v.validate.Struct(receipt)
}