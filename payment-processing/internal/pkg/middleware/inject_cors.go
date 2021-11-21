package middleware

import (
	"net/http"
)

func InjectCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// wildcard doesn't work if client using credential mode, thus need conditional values
		origin := "*"
		if r.Header.Get("Origin") != "" {
			origin = r.Header.Get("Origin")
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		next.ServeHTTP(w, r)
	})
}
