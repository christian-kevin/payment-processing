package store

type Limit struct {
	ID         int64 `json:"id" db:"id"`
	ParentID   int64 `json:"parent_id" db:"parent_id"`
	ParentType int   `json:"parent_type" db:"parent_type"`
	Type       int   `json:"type" db:"type"`
	Amount     int64 `json:"amount" db:"amount"`
}
