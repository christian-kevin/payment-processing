package redis

type Store interface {
	Ping() error
}

type RWalletStore interface {
	LockCreateWallet(userID int64) error
	ReleaseLockCreateWallet(userID int64) error

	LockUpdateWallet(userID, walletID int64) error
	ReleaseLockUpdateWallet(userID, walletID int64) error
}

type RCardStore interface {
	LockCreateCard(walletID int64) error
	ReleaseLockCreateCard(walletID int64) error

	CreateTransaction(limit int64, amount int64, cardID int64) error
}