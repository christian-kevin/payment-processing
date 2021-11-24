package ratelimiter

import (
	"context"
	"spenmo/payment-processing/payment-processing/config"
	"spenmo/payment-processing/payment-processing/internal/pkg/constant"
	"spenmo/payment-processing/payment-processing/internal/pkg/errutil"
	"spenmo/payment-processing/payment-processing/internal/pkg/store/redis"
	log "spenmo/payment-processing/payment-processing/pkg/logger"
	"time"
)

type rateLimiter struct {
	rRateLimitStore redis.RateLimitStore
}

func NewRateLimiter (rRateLimitStore redis.RateLimitStore) RateLimiter {
	return &rateLimiter{rRateLimitStore: rRateLimitStore}
}

func (r *rateLimiter) Allow(ctx context.Context, page string) error {
	var ttl time.Duration
	if config.AppConfig.RateLimitUnit == constant.Minute {
		ttl = time.Duration(config.AppConfig.RateLimitValue) * time.Minute
	} else if config.AppConfig.RateLimitUnit == constant.Hour {
		ttl = time.Duration(config.AppConfig.RateLimitValue) * time.Hour
	} else {
		log.Get().Error(ctx, "wrong rate limit unit in config")
		return errutil.ErrInvalidParam
	}

	rRate := redis.RateLimit{
		Limit: config.AppConfig.RateLimitValue,
		Ttl:   ttl,
		Unit:  config.AppConfig.RateLimitUnit,
	}
	return r.rRateLimitStore.Allow(page, &rRate)
}
