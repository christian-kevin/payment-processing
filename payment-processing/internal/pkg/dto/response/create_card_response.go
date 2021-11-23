package response

type CreateCard struct {
	CardID     int64  `json:"card_id"`
	CardNumber string `json:"card_number"`
	ExpiryDate string `json:"expiry_date"`
	Name       string `json:"name"`
}
