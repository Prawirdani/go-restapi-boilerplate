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

type RequestLogAttributes struct {
	Method     string
	Uri        string
	ClientIP   string
	RequestID  string
	StatusCode int
	StatusText string
	TimeTaken  time.Duration
}

type ResRecorder struct {
	http.ResponseWriter
	Status int
	Body   []byte
}

func (rr *ResRecorder) WriteHeader(code int) {
	rr.Status = code
	rr.ResponseWriter.WriteHeader(code)
}

func HttpRequestLogger(rl RequestLogAttributes) {
	slog.Info(
		"HTTP Request Log",
		slog.String("method", rl.Method),
		slog.String("url", rl.Uri),
		slog.String("from", rl.ClientIP),
		slog.String("request-id", rl.RequestID),
		slog.Int("status_code", rl.StatusCode),
		slog.String("status_text", rl.StatusText),
		slog.Float64("time_taken(ms)", float64(rl.TimeTaken.Microseconds())/float64(1000)),
	)
}
