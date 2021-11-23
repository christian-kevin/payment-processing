package middleware

import (
	"context"
	"net/http"
	"spenmo/payment-processing/payment-processing/internal/pkg/contextutil"
	"spenmo/payment-processing/payment-processing/internal/pkg/dto/response"
	"spenmo/payment-processing/payment-processing/internal/pkg/errutil"
	log "spenmo/payment-processing/payment-processing/pkg/logger"
	"strings"
)

type tenant struct{}

func NewTenant() *tenant {
	return &tenant{}
}

func (t *tenant) Enforce(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenant := r.Header.Get("X-Tenant")
		if len(tenant) != 2 {
			response.WriteResponse(w, nil, errutil.ErrInvalidParam)
			return
		}

		ctx := context.WithValue(r.Context(), contextutil.XTenantKey, strings.ToLower(tenant))
		ctx = log.BuildContextLoggerWithCountry(ctx, strings.ToLower(tenant))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
