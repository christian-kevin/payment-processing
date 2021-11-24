package middleware

import (
	"net/http"
	"spenmo/payment-processing/payment-processing/internal/pkg/dto/response"
	"spenmo/payment-processing/payment-processing/internal/pkg/errutil"
	"spenmo/payment-processing/payment-processing/internal/pkg/ratelimiter"
	log "spenmo/payment-processing/payment-processing/pkg/logger"
)

type MustRateLimit struct {
	limiter  ratelimiter.RateLimiter
	isActive bool
}

func NewMustRateLimit(limiter ratelimiter.RateLimiter, isActive bool) *MustRateLimit {
	return &MustRateLimit{limiter: limiter, isActive: isActive}
}

func (l *MustRateLimit) Enforce(next http.Handler, page string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if l.isActive {
			if page != "" {
				err := l.limiter.Allow(r.Context(), page)
				if err != nil && err != errutil.ErrRateLimitExceeded {
					log.Get().Errorf(r.Context(), "failed to get rate limit, err :%s", err.Error())
					response.WriteResponse(w, nil, errutil.ErrServerError)
					return
				}
				if err == errutil.ErrRateLimitExceeded {
					response.WriteResponse(w, nil, errutil.ErrRateLimitExceeded)
					return
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
