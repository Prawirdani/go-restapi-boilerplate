package logger

import (
	"io"
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
	Method     string
	Uri        string
	Address    string
	StatusCode int
	StatusText string
	TimeTaken  time.Duration
	Body       []byte
}

type resRecorder struct {
	http.ResponseWriter
	Status int
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
		slog.String("from", rl.Address),
		slog.Int("status_code", rl.StatusCode),
		slog.String("status_text", rl.StatusText),
		slog.String("body", string(rl.Body)),
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
		duration := time.Since(start)

		logAttributes := &requestLogAttributes{
			Method:     r.Method,
			Uri:        r.RequestURI,
			Address:    r.RemoteAddr,
			StatusCode: rec.Status,
			StatusText: http.StatusText(rec.Status),
			TimeTaken:  duration,
		}

		if r.Method == http.MethodPost && rec.Status > 201 {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				slog.Error("requestLogger.bodyCatcher", "cause", err)
			}
			defer r.Body.Close()

			logAttributes.Body = body
		}

		HttpRequestLogger(*logAttributes)
	})
}
