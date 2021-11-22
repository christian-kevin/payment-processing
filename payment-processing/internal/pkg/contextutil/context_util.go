package contextutil

import (
	"context"
	"spenmo/payment-processing/payment-processing/internal/pkg/errutil"
	log "spenmo/payment-processing/payment-processing/pkg/logger"
	"strings"
)

type key int

const (
	XTenantKey = key(iota)
	XUserIDKey
)

func ExtractTenantID(ctx context.Context) (string, error) {
	val, ok := ctx.Value(XTenantKey).(string)
	if !ok {
		return "", errutil.ErrContextValueNotFound
	}

	return strings.ToLower(val), nil
}

func ExtractUserID(ctx context.Context) (int64, error) {
	val, ok := ctx.Value(XUserIDKey).(int64)
	if !ok {
		log.Get().Error(log.GetEmptyContext(), "value context", ctx.Value(XUserIDKey))
		return 0, errutil.ErrContextValueNotFound
	}

	return val, nil
}
