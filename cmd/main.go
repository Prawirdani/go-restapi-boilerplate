package main

import (
	"github.com/prawirdani/go-restapi-boilerplate/app"
	"github.com/prawirdani/go-restapi-boilerplate/config"
	"github.com/prawirdani/go-restapi-boilerplate/internal/index"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/utils"
)

func main() {
	conf, err := config.LoadConfig()
	utils.PanicIfErr(err)

	r := app.NewMainRouter(conf)
	svr := app.NewServer(conf, r)

	indexHandler := index.NewIndexHandler()
	indexHandler.Routes(r)

	svr.Start()
}
