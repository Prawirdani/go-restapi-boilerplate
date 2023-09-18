package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/httputil"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/logger"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/unrolled/secure"
)

func NewMainRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(panicRecoverer)
	r.Use(logger.RequestLogger)

	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Add Allowed Origins, eg: frontend
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler)

	r.Use(secure.New(secure.Options{
		ContentTypeNosniff: true,
	}).Handler)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		httputil.SendError(w, httputil.ErrNotFound("ops! you must be lost!"))
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		httputil.SendError(w, httputil.ErrMethodNotAllowed("ops! method not allowed"))
	})

	r.Use(middleware.Compress(6))

	// Enable Swagger on development environment
	if os.Getenv("ENV") == "development" {
		slog.Info(fmt.Sprintf("Swagger available at http://localhost:%s/swagger/index.html", os.Getenv("APP_PORT")))
		r.Get("/swagger/*", httpSwagger.WrapHandler)
	}

	return r
}

func NewSubRouter() *chi.Mux {
	return chi.NewRouter()
}

/* Panic recoverer middleware, it keep the service alive when crashes */
func panicRecoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				slog.Error("Recoverer Log", "cause", rvr)
				httputil.SendError(w, fmt.Errorf("%v", rvr))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
