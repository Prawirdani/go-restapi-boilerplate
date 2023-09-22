package middleware

import (
	"net/http"
	"time"

	"github.com/prawirdani/go-restapi-boilerplate/pkg/logger"
)

/* Request Logger Middleware*/
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := &logger.ResRecorder{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}
		next.ServeHTTP(rec, r)

		// Retrieve forwarded client ip from reverse-proxy
		forwardedIP := r.Header.Get("X-Real-IP") // !Todo should create dynamic version, if not using reverse-proxy then straighly collect the client ip address.
		duration := time.Since(start)

		logAttributes := &logger.RequestLogAttributes{
			Method:      r.Method,
			Uri:         r.RequestURI,
			ForwardedIP: forwardedIP,
			StatusCode:  rec.Status,
			StatusText:  http.StatusText(rec.Status),
			TimeTaken:   duration,
		}

		logger.HttpRequestLogger(*logAttributes)
	})
}
