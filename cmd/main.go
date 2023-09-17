package main

import (
	"context"
	"flag"
	"log/slog"

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

//	@title			Swagger Docs (RESTAPI BoilerPlate)
//	@version		1.0
//	@description	This is an api Swagger.
//	@BasePath		/v1
func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		slog.Error("Load config error", "error", err)
	}

	pgDB := db.NewPostgreSQL(*conf)
	defer pgDB.Close(context.Background())

	var initSchema bool
	flag.BoolVar(&initSchema, "initschema", false, "Initialize the schema")
	flag.Parse()
	if initSchema {
		db.InitSchema(pgDB)
	}

	mainRouter := app.NewMainRouter(conf)
	appRouter := app.NewSubRouter()

	userRepository := user.NewUserRepository(pgDB)
	userService := user.NewUserService(userRepository)
	userHandler := user.NewUserHandler(userService)

	appRouter.Route("/v1", func(r chi.Router) {
		userHandler.Routes(r)
	})

	mainRouter.Mount("/", appRouter)

	svr := app.NewServer(conf, mainRouter)
	svr.Start()
}

func someHandler(c chi.Context) {
}
