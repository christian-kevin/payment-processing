package handler

import (
	"spenmo/payment-processing/payment-processing/internal/pkg/store/mysql"
	"spenmo/payment-processing/payment-processing/internal/pkg/store/redis"
)

type Module struct {
	rCardStore              redis.RCardStore
	rWalletStore            redis.RWalletStore
	walletStore             mysql.WalletStore
	cardStore               mysql.CardStore
	limitStore              mysql.LimitStore
	cardTransactionLogStore mysql.CardTransactionLogStore
	walletBalanceLogStore   mysql.WalletBalanceLogStore
}

func NewModule(rCardStore redis.RCardStore,
	rWalletStore redis.RWalletStore,
	walletStore mysql.WalletStore,
	cardStore mysql.CardStore,
	limitStore mysql.LimitStore,
	cardTransactionLogStore mysql.CardTransactionLogStore,
	walletBalanceLogStore mysql.WalletBalanceLogStore,
) *Module {
	return &Module{
		rCardStore:              rCardStore,
		rWalletStore:            rWalletStore,
		walletStore:             walletStore,
		limitStore:              limitStore,
		cardStore:               cardStore,
		cardTransactionLogStore: cardTransactionLogStore,
		walletBalanceLogStore:   walletBalanceLogStore,
	}
}
