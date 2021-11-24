package request

type CreateTransaction struct {
	CardNumber string `json:"card_number"`
	ExpiryDate string `json:"expiry_date"`
	Amount     string `json:"amount"`
	Country    string `json:"country"`
}
