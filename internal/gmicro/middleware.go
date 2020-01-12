package gmicro

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// LoggingMiddleware that logs requests received.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"uri":    r.URL,
			"method": r.Method,
		}).Info("Request received")

		next.ServeHTTP(rw, r)
	})
}

// ContentTypeMiddleware that sets application/json in response header.
func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
