package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/prawirdani/go-restapi-boilerplate/app"
	"github.com/prawirdani/go-restapi-boilerplate/config"
	_ "github.com/prawirdani/go-restapi-boilerplate/docs"
	"github.com/prawirdani/go-restapi-boilerplate/internal/index"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/utils"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title			Swagger Docs (API TITLE HERE)
// @version		1.0
// @description	This is an api Swagger.
// @BasePath		/v1
func main() {
	conf, err := config.LoadConfig()
	utils.PanicIfErr(err)

	mainRouter := app.NewMainRouter(conf)
	mainRouter.Get("/swagger/*", httpSwagger.WrapHandler)

	appRouter := chi.NewRouter()

	indexHandler := index.NewIndexHandler()
	appRouter.Route("/v1", func(r chi.Router) {
		indexHandler.Routes(r)
	})

	mainRouter.Mount("/", appRouter)

	svr := app.NewServer(conf, mainRouter)
	svr.Start()
}
