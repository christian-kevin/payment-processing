package ratelimiter

import "context"

type RateLimiter interface {
	Allow(ctx context.Context, page string) error
}
