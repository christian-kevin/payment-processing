package store

type CardTransactionLog struct {
	ID        int64 `json:"id" db:"id"`
	CardID    int64 `json:"card_id" db:"card_id"`
	Amount    int64 `json:"amount" db:"amount"`
	CreatedAt int64 `json:"created_at" db:"created_at"`
}
