package redis

import (
	"fmt"
	"spenmo/payment-processing/payment-processing/config"
	"spenmo/payment-processing/payment-processing/pkg/cache"
	"strconv"
	"time"
)

type rCardStore struct {
	client cache.Cache
}

func NewCardStore(cache cache.Cache) RCardStore {
	return &rCardStore{client: cache}
}

type cardKey struct {
	walletID int64
}

func (k cardKey) toCreateKey() string {
	return fmt.Sprintf("%s:create_card:%d", config.AppConfig.Env, k.walletID)
}

const (
	redisTTLLockCard = 5 * time.Minute
)

func (s *rCardStore) LockCreateCard(walletID int64) error {
	return s.client.SetNX(cardKey{walletID: walletID}.toCreateKey(), strconv.FormatBool(true), redisTTLLockCard)
}

func (s *rCardStore) ReleaseLockCreateCard(walletID int64) error {
	return s.client.Del(cardKey{walletID: walletID}.toCreateKey())
}

func (s *rCardStore) CreateTransaction(limit int64, amount int64, cardID int64) error {
	panic("implement me")
}