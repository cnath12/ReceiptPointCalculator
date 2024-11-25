package service

import (
    "context"
    "math"
    "strconv"
    "strings"
    "time"
    "unicode"
    "ReceiptPointCalculator/internal/domain/model"
	"ReceiptPointCalculator/pkg/utils"
)

type ReceiptRepository interface {
    Save(ctx context.Context, receipt *model.Receipt) error
    GetByID(ctx context.Context, id string) (*model.Receipt, error)
}

type ReceiptService struct {
    repo ReceiptRepository
}

func NewReceiptService(repo ReceiptRepository) *ReceiptService {
    return &ReceiptService{repo: repo}
}

func (s *ReceiptService) ProcessReceipt(ctx context.Context, receipt *model.Receipt) (string, error) {
    id := utils.GenerateID()
    receipt.ID = id
    
    if err := s.repo.Save(ctx, receipt); err != nil {
        return "", err
    }
    
    return id, nil
}

func (s *ReceiptService) GetPoints(ctx context.Context, id string) (int64, error) {
    receipt, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return 0, err
    }
    
    return s.CalculatePoints(receipt), nil
}

func (s *ReceiptService) CalculatePoints(receipt *model.Receipt) int64 {
    var points int64

    // Rule 1: One point for every alphanumeric character in the retailer name
    for _, char := range receipt.Retailer {
        if unicode.IsLetter(char) || unicode.IsNumber(char) {
            points++
        }
    }

    // Rule 2: 50 points if the total is a round dollar amount
    total, _ := strconv.ParseFloat(receipt.Total, 64)
    if total == float64(int64(total)) {
        points += 50
    }

    // Rule 3: 25 points if the total is a multiple of 0.25
    if math.Mod(total*100, 25) == 0 {
        points += 25
    }

    // Rule 4: 5 points for every two items
    points += int64(len(receipt.Items) / 2 * 5)

    // Rule 5: Points for items with description length multiple of 3
    for _, item := range receipt.Items {
        trimmedLen := len(strings.TrimSpace(item.ShortDescription))
        if trimmedLen%3 == 0 {
            price, _ := strconv.ParseFloat(item.Price, 64)
            points += int64(math.Ceil(price * 0.2))
        }
    }

    // Rule 6: 6 points if the day in the purchase date is odd
    purchaseDate, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
    if purchaseDate.Day()%2 == 1 {
        points += 6
    }

    // Rule 7: 10 points if the time of purchase is between 2:00pm and 4:00pm
    purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
    hour := purchaseTime.Hour()
    minute := purchaseTime.Minute()
    if (hour == 14 && minute > 0) || (hour == 15) || (hour == 13 && minute < 60) {
        points += 10
    }

    return points
}

func generateUUID() string {
    return utils.GenerateID()
}