package redis

import (
	"fmt"
	"spenmo/payment-processing/payment-processing/config"
	"spenmo/payment-processing/payment-processing/pkg/cache"
	"strconv"
	"time"
)

type rWalletStore struct {
	client cache.Cache
}

func NewWalletStore(cache cache.Cache) RWalletStore {
	return &rWalletStore{client: cache}
}

type walletKey struct {
	userID  int64
	walletID int64
}

func (k walletKey) toCreateKey() string {
	return fmt.Sprintf("%s:create_wallet:%d", config.AppConfig.Env, k.userID)
}

func (k walletKey) toUpdateKey() string {
	return fmt.Sprintf("%s:update_wallet:%d:%d", config.AppConfig.Env, k.walletID, k.userID)
}

const (
	redisTTLLockWallet = 5 * time.Minute
)

func (s *rWalletStore) LockCreateWallet(userID int64) error {
	return s.client.SetNX(walletKey{userID: userID}.toCreateKey(), strconv.FormatBool(true), redisTTLLockWallet)
}

func (s *rWalletStore) ReleaseLockCreateWallet(userID int64) error {
	return s.client.Del(walletKey{userID: userID}.toCreateKey())
}

func (s *rWalletStore) LockUpdateWallet(userID, walletID int64) error {
	return s.client.SetNX(walletKey{userID: userID, walletID: walletID}.toUpdateKey(), strconv.FormatBool(true),
	redisTTLLockWallet)
}

func (s *rWalletStore) ReleaseLockUpdateWallet(userID, walletID int64) error {
	return s.client.Del(walletKey{userID: userID, walletID: walletID}.toUpdateKey())
}
