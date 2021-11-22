package handler

import "spenmo/payment-processing/payment-processing/internal/pkg/store/redis"

type Module struct {
	rCardStore   redis.RCardStore
	rWalletStore redis.RWalletStore
}

func NewModule(rCardStore redis.RCardStore,
	rWalletStore redis.RWalletStore,
) *Module {
	return &Module{
		rCardStore:   rCardStore,
		rWalletStore: rWalletStore,
	}
}
