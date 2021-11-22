package store

type Card struct {
	ID         int64  `json:"id" db:"id"`
	WalletID   int64  `json:"wallet_id" db:"wallet_id"`
	CardNumber string `json:"card_number" db:"card_number"`
	ExpiryDate string `json:"expiry_date" db:"expiry_date"`
	Name       string `json:"name" db:"name"`
	CreatedAt  int64  `json:"created_at" db:"created_at"`
	IsDeleted  int    `json:"is_deleted" db:"is_deleted"`
}
