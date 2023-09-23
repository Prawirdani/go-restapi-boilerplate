package main

import (
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/prawirdani/go-restapi-boilerplate/app"
	"github.com/prawirdani/go-restapi-boilerplate/config"
	"github.com/prawirdani/go-restapi-boilerplate/db"
	_ "github.com/prawirdani/go-restapi-boilerplate/docs"
	"github.com/prawirdani/go-restapi-boilerplate/internal/auth"
	"github.com/prawirdani/go-restapi-boilerplate/internal/user"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/logger"
)

func init() {
	logger.InitLogger()
}

// @title			Swagger Docs (RESTAPI BoilerPlate)
// @version		1.0
// @description	This is an api Swagger.
// @BasePath		/v1
func main() {
	if err := config.LoadEnv(); err != nil {
		slog.Error("env load error", "cause", err)
		os.Exit(1)
	}

	psqlDB := db.NewPostgreSQL()
	defer psqlDB.Close()

	mainRouter := app.NewMainRouter()
	v1 := app.NewSubRouter()

	userRepository := user.NewUserRepository(psqlDB)

	authService := auth.NewAuthService(userRepository)
	userService := user.NewUserService(userRepository)

	authHandler := auth.NewAuthHandler(authService)
	userHandler := user.NewUserHandler(userService)

	v1.Route("/v1", func(rt chi.Router) {
		rt.Route("/users", func(r chi.Router) {
			userHandler.Routes(r)
		})
		rt.Route("/auth", func(r chi.Router) {
			authHandler.Routes(r)
		})
	})

	mainRouter.Mount("/", v1)

	svr := app.NewServer(mainRouter)
	svr.Start()
}
