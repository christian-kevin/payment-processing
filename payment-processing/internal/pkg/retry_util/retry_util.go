package retry_util

import (
	"context"
	log "spenmo/payment-processing/payment-processing/pkg/logger"
	"time"
)

func ExponentialRetry(ctx context.Context, numRetry int, cmd func() error) (err error) {
	sleepTime := 500 * time.Millisecond
	for i := 0; i < numRetry; i++ {
		if i > 0 {
			log.Get().Errorf(ctx, "retry attempt: %d", i)
		}

		err = cmd()
		if err == nil {
			break
		}
		time.Sleep(sleepTime)
		sleepTime *= 2
	}
	return err
}
