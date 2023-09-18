package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	*http.Server
}

func NewServer(handler http.Handler) *Server {
	svr := http.Server{
		Addr:         ":" + os.Getenv("APP_PORT"),
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	return &Server{&svr}
}

func (s *Server) Start() {
	go func() {
		slog.Info(fmt.Sprintf("Server listening on localhost%s", s.Addr))
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Debug("Server startup failed", slog.Any("cause", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	slog.Info("Shutdown signal received")

	ctx, shutdown := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdown()

	if err := s.Shutdown(ctx); err != nil {
		slog.Warn("Server shutdown failed", slog.Any("cause", err))
		os.Exit(1)
	}
	slog.Info("Server gracefully stopped")
}
