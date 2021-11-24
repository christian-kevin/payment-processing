package redis

import (
	"fmt"
	"spenmo/payment-processing/payment-processing/config"
	"spenmo/payment-processing/payment-processing/internal/pkg/constant"
	"spenmo/payment-processing/payment-processing/internal/pkg/errutil"
	"spenmo/payment-processing/payment-processing/internal/pkg/timeutil"
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

func buildCardLimitKey(cardID int64, retType int) string {
	return fmt.Sprintf("%s:transaction_limit:%d:%d", config.AppConfig.Env, cardID, retType)
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

/**
	TODO add local time
 */
func (s *rCardStore) CreateTransaction(limit int64, amount int64, cardID int64, retType int) error {
	var ttl time.Duration
	if retType == constant.LimitTypeDaily {
		now := time.Now()
		bodToday := timeutil.Bod(now)
		bodTomorrow := bodToday.AddDate(0, 0, 1)
		ttl = bodTomorrow.Sub(now)
	} else if retType == constant.LimitTypeMonthly {
		now := time.Now()
		currentYear, currentMonth, _ := now.Date()
		currentLocation := now.Location()
		firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
		ttl = firstOfMonth.Sub(now)
	} else {
		return errutil.ErrInvalidParam
	}

	_, err := s.client.DecrX(buildCardLimitKey(cardID, retType), limit, amount, ttl)
	return err
}