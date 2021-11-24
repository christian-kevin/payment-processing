package redis

import (
	"fmt"
	"spenmo/payment-processing/payment-processing/internal/pkg/constant"
	"spenmo/payment-processing/payment-processing/internal/pkg/errutil"
	"spenmo/payment-processing/payment-processing/pkg/cache"
	"time"
)

type rRateLimitStore struct {
	client cache.Cache
}

func NewRateLimitStore(cache cache.Cache) RateLimitStore {
	return &rRateLimitStore{client: cache}
}

func createRLimitKey(page string, r *RateLimit) string {
	return fmt.Sprintf("rate_limit:%s:%d", page, getCurrentUnitValue(r))
}

func getCurrentUnitValue(r *RateLimit) int64 {
	if r.Unit == constant.Hour {
		return time.Now().Truncate(time.Hour).Unix()
	}
	if r.Unit == constant.Minute {
		return time.Now().Truncate(time.Minute).Unix()
	}
	return time.Now().Truncate(time.Second).Unix()
}

func (s *rRateLimitStore) Allow(page string, r *RateLimit) error {
	_, err := s.client.DecrX(createRLimitKey(page, r), r.Limit, 1, r.Ttl)
	if err != nil {
		if err == cache.ErrLimitExceeded {
			return errutil.ErrRateLimitExceeded
		}
		return err
	}
	return nil
}
