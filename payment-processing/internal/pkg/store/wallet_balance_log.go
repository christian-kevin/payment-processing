package store

type WalletBalanceLog struct {
	ID        int64 `json:"id" db:"id"`
	WalletID  int64 `json:"wallet_id" db:"wallet_id"`
	Amount    int64 `json:"amount" db:"amount"`
	Type      int   `json:"type" db:"type"`
	CreatedAt int64 `json:"created_at" db:"created_at"`
}
