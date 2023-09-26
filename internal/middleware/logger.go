package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/logger"
)

/* Request Logger Middleware*/
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetReqID(r.Context())
		start := time.Now()
		rec := &logger.ResRecorder{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}
		next.ServeHTTP(rec, r)

		duration := time.Since(start)

		logAttributes := &logger.RequestLogAttributes{
			Method:     r.Method,
			Uri:        r.RequestURI,
			ClientIP:   r.RemoteAddr,
			RequestID:  requestID,
			StatusCode: rec.Status,
			StatusText: http.StatusText(rec.Status),
			TimeTaken:  duration,
		}

		logger.HttpRequestLogger(*logAttributes)
	})
}
