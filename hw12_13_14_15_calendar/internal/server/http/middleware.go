package internalhttp

import (
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler, log Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.InfoHttp(r, http.StatusOK, time.Since(start))
	})
}
