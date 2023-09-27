package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/httputil"
	httpSwagger "github.com/swaggo/http-swagger"

	middleware "github.com/prawirdani/go-restapi-boilerplate/internal/middleware"
	"github.com/unrolled/secure"
)

func NewMainRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(panicRecoverer)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.RequestID)
	r.Use(middleware.RequestLogger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:4173", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Request-Id"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(secure.New(secure.Options{
		ContentTypeNosniff: true,
		FrameDeny: true,
	}).Handler)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		httputil.SendError(w, httputil.ErrNotFound("ops! you must be lost!"))
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		httputil.SendError(w, httputil.ErrMethodNotAllowed("ops! method not allowed"))
	})

	r.Use(chiMiddleware.Compress(6))

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

