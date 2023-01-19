package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler, log Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Info(fmt.Sprintf("%s [%s] %s %s %s %v %s %s", r.RemoteAddr, time.Now().Format(time.RFC822Z),
			r.Method, r.RequestURI, r.URL.Scheme, http.StatusOK, time.Since(start), r.Header.Get("User-Agent")))
	})
}
