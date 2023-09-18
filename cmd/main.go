package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/prawirdani/go-restapi-boilerplate/app"
	"github.com/prawirdani/go-restapi-boilerplate/config"
	"github.com/prawirdani/go-restapi-boilerplate/db"
	_ "github.com/prawirdani/go-restapi-boilerplate/docs"
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
		slog.Error("env variable not provided", "cause", err)
		os.Exit(1)
	}

	pgDB := db.NewPostgreSQL()
	defer pgDB.Close(context.Background())

	mainRouter := app.NewMainRouter()
	appRouter := app.NewSubRouter()

	userRepository := user.NewUserRepository(pgDB)
	userService := user.NewUserService(userRepository)
	userHandler := user.NewUserHandler(userService)

	appRouter.Route("/v1", func(r chi.Router) {
		userHandler.Routes(r)
	})

	mainRouter.Mount("/", appRouter)

	svr := app.NewServer(mainRouter)
	svr.Start()
}
