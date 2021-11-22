package middleware

import (
	"context"
	"net/http"
	"spenmo/payment-processing/payment-processing/internal/pkg/contextutil"
	"strconv"
)

type auth struct{}

func NewAuth() *auth {
	return &auth{}
}

const (
	userIDHeaderKey = "Op-User-ID"
	AuthCookieKey = "OP-AU"
)

func (a *auth) Enforce(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		userIDStr := r.Header.Get(userIDHeaderKey)
		if userIDStr == "" {
			next.ServeHTTP(rw, r)
			return
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			next.ServeHTTP(rw, r)
			return
		}

		authCookie := r.Header.Get(AuthCookieKey)
		if authCookie == "" {
			next.ServeHTTP(rw, r)
			return
		}

		ctx := context.WithValue(r.Context(), contextutil.XUserIDKey, userID)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}