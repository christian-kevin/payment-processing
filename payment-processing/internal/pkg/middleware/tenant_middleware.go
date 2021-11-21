package middleware

import (
	"context"
	"net/http"
	"spenmo/payment-processing/payment-processing/internal/pkg/contextutil"
	log "spenmo/payment-processing/payment-processing/pkg/logger"
	"strings"
)

type tenant struct{}

func NewTenant() *tenant {
	return &tenant{}
}

func (t *tenant) Enforce(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenant := "ID"
		// TODO Get Tenant Here

		ctx := context.WithValue(r.Context(), contextutil.XTenantKey, strings.ToLower(tenant))
		ctx = log.BuildContextLoggerWithCountry(ctx, strings.ToLower(tenant))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
