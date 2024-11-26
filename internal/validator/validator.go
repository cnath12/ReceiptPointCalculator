package validator

import (
    "regexp"
    "time"
    "github.com/go-playground/validator/v10"
	"math"
	"ReceiptPointCalculator/internal/domain/model"
	"strconv"
)

type ReceiptValidator struct {
    validate *validator.Validate
}

func NewReceiptValidator() *ReceiptValidator {
    v := validator.New()
    
    // Register custom validation functions
    v.RegisterValidation("retailer", validateRetailer)
	v.RegisterValidation("shortDescription", validateShortDescription)
    v.RegisterValidation("date", validateDate)
    v.RegisterValidation("time", validateTime)
    v.RegisterValidation("price", validatePrice)
    v.RegisterValidation("totalMatch", validateTotal)

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

func validateShortDescription(fl validator.FieldLevel) bool {
    desc := fl.Field().String()
    matched, _ := regexp.MatchString(`^[\w\s\-]+$`, desc)
    return matched
}

func validateTotal(fl validator.FieldLevel) bool {
    receipt, ok := fl.Parent().Interface().(model.Receipt)
    if !ok {
        return false
    }

    total, err := strconv.ParseFloat(receipt.Total, 64)
    if err != nil {
        return false
    }

    var sum float64
    for _, item := range receipt.Items {
        price, err := strconv.ParseFloat(item.Price, 64)
        if err != nil {
            return false
        }
        sum += price
    }

    // Round to 2 decimal places to avoid floating point precision issues
    sum = math.Round(sum*100) / 100
    total = math.Round(total*100) / 100

    return sum == total
}