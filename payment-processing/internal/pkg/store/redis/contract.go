package redis

import "time"

type Store interface {
	Ping() error
}

type RWalletStore interface {
	LockCreateWallet(userID int64) error
	ReleaseLockCreateWallet(userID int64) error

	LockUpdateWallet(walletID int64) error
	ReleaseLockUpdateWallet(walletID int64) error
}

type RCardStore interface {
	LockCreateCard(walletID int64) error
	ReleaseLockCreateCard(walletID int64) error

	CreateTransaction(limit int64, amount int64, cardID int64, retType int) error
}

type RateLimitStore interface {
	Allow(page string, r *RateLimit) error
}



type RateLimit struct {
	Limit int64
	Ttl   time.Duration
	Unit  string
}
