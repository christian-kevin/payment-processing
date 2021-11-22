package store

type Wallet struct {
	ID         int64  `json:"id" db:"id"`
	Balance    int64  `json:"balance" db:"balance"`
	ParentID   int64  `json:"parent_id" db:"parent_id"`
	ParentType int    `json:"parent_type" db:"parent_type"`
	Country    string `json:"country" db:"country"`
}
