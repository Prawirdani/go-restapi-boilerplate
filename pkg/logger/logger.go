package logger

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

func InitLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)
}

type requestLogAttributes struct {
	Method      string
	Uri         string
	ForwardedIP string
	StatusCode  int
	StatusText  string
	TimeTaken   time.Duration
}

type resRecorder struct {
	http.ResponseWriter
	Status int
	Body   []byte
}

func (rr *resRecorder) WriteHeader(code int) {
	rr.Status = code
	rr.ResponseWriter.WriteHeader(code)
}

func HttpRequestLogger(rl requestLogAttributes) {
	slog.Info(
		"HTTP Request Log",
		slog.String("method", rl.Method),
		slog.String("url", rl.Uri),
		slog.String("from", rl.ForwardedIP),
		slog.Int("status_code", rl.StatusCode),
		slog.String("status_text", rl.StatusText),
		slog.Float64("time_taken(ms)", float64(rl.TimeTaken.Microseconds())/float64(1000)),
	)
}

/* Request Logger Middleware*/
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := &resRecorder{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}
		next.ServeHTTP(rec, r)

		// Retrieve forwarded client ip from reverse-proxy
		forwardedIP := r.Header.Get("X-Real-IP")
		duration := time.Since(start)

		logAttributes := &requestLogAttributes{
			Method:      r.Method,
			Uri:         r.RequestURI,
			ForwardedIP: forwardedIP,
			StatusCode:  rec.Status,
			StatusText:  http.StatusText(rec.Status),
			TimeTaken:   duration,
		}

		HttpRequestLogger(*logAttributes)
	})
}
