package app

import (
	"context"
	"log"
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
		log.Printf("ENVIROMENT: %s", s.Env)
		log.Printf("Server listening on localhost%s", s.Addr)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server startup error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received")

	ctx, shutdown := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdown()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server gracefully stopped")
}
