package app

import (
	"github.com/Eiji224/go_ozon/hw2/cart/internal/router"
)

type App struct {
	router *router.Router
}

func InitApp(router *router.Router) *App {
	return &App{
		router: router,
	}
}

func (app *App) StartServer() {
	app.router.Serve()
}
