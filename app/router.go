package app

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prawirdani/go-restapi-boilerplate/config"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/httpError"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/json"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/logger"
	"github.com/rs/cors"
	"github.com/unrolled/secure"
)

func NewMainRouter(c *config.Config) *chi.Mux {
	r := chi.NewRouter()

	if c.Server.Env == "development" {
		r.Use(logger.RequestLogger) // Json formatted Request Log for Prod Env
	} else {
		r.Use(middleware.DefaultLogger) // Human Readable Request Log for Dev Env
	}

	r.Use(panicRecoverer)

	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Add Allowed Origins, eg: frontend
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: c.Cors.Credentials,
	}).Handler)

	r.Use(secure.New(secure.Options{
		ContentTypeNosniff: true,
	}).Handler)

	r.NotFound(notFoundHandler)
	r.MethodNotAllowed(methodNotAllowed)

	r.Use(middleware.Compress(6))

	return r
}

func NewSubRouter() *chi.Mux {
	return chi.NewRouter()
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	slog.Warn("Route Not Found", "Route", r.RequestURI)
	json.SendError(w, httpError.NotFound("ops! you must be lost!"))
}

func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	slog.Warn("Method Not Allowed", slog.String("Route", r.RequestURI), slog.String("method", r.Method))
	json.SendError(w, httpError.MethodNotAllowed("ops! method not allowed"))
}

/* Panic recoverer middleware, it keep the service alive when crashes */
func panicRecoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				slog.Error("Recoverer Log", "cause", rvr)
				json.SendError(w, fmt.Errorf("%v", rvr))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
