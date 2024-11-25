package model

type Receipt struct {
    ID            string  `json:"id,omitempty"`
    Retailer      string  `json:"retailer" validate:"required,retailer"`
    PurchaseDate  string  `json:"purchaseDate" validate:"required,date"`
    PurchaseTime  string  `json:"purchaseTime" validate:"required,time"`
    Items         []Item  `json:"items" validate:"required,min=1,dive"`
    Total         string  `json:"total" validate:"required,price"`
}

type Item struct {
    ShortDescription string `json:"shortDescription" validate:"required"`
    Price           string `json:"price" validate:"required,price"`
}