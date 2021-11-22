package redis

type Store interface {
	Ping() error
}

type RWalletStore interface {
	LockCreateWallet(userID int64)
	ReleaseLockCreateWallet(userID int64)

	LockUpdateWallet(walletID int64)
	ReleaseLockUpdateWallet(walletID int64)
}

type RCardStore interface {
	LockCreateCard(walletID int64)
	ReleaseLockCreateCard(walletID int64)

	CreateTransaction(limit int64, amount int64, cardID int64)
}