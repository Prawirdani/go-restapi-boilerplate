package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prawirdani/go-restapi-boilerplate/config"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/httpErr"
	"github.com/rs/cors"
	"github.com/unrolled/secure"
)

func NewMainRouter(c *config.Config) *chi.Mux {
	r := chi.NewRouter()

	r.Use(panicRecover)

	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Add Allowed Origins, eg: frontend
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: c.Cors.Credentials,
	}).Handler)

	r.Use(secure.New(secure.Options{
		ContentTypeNosniff: true,
	}).Handler)

	if c.Server.Env == "development" {
		r.Use(middleware.DefaultLogger)
	}

	r.Use(middleware.Compress(6))

	return r
}

func panicRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				/* Add Logger here to detect untracked malfunction/error */
				log.Println("Recover:", rvr)
				httpErr.ExceptionIfErr(w, fmt.Errorf("%v", rvr))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
