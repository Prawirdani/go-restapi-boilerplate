package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prawirdani/go-restapi-boilerplate/config"
)

type Server struct {
	*http.Server
	Env string
}

func NewServer(c *config.Config, handler http.Handler) *Server {
	svr := http.Server{
		Addr:         c.Server.Port,
		Handler:      handler,
		ReadTimeout:  c.Server.ReadTimeout * time.Second,
		WriteTimeout: c.Server.WriteTimeout * time.Second,
	}
	return &Server{&svr, c.Server.Env}
}

func (s *Server) Start() {
	go func() {
		slog.Info("Server started", slog.String("environment", s.Env), slog.String("port", s.Addr))
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
