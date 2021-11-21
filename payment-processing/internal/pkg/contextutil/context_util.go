package contextutil

import (
	"context"
	"spenmo/payment-processing/payment-processing/internal/pkg/errutil"
	"strings"
)

type key int64

const (
	XTenantKey key = 1
)

func ExtractTenantID(ctx context.Context) (string, error) {
	val, ok := ctx.Value(XTenantKey).(string)
	if !ok {
		return "", errutil.ErrContextValueNotFound
	}

	return strings.ToLower(val), nil
}
