package main

import (
	"github.com/Eiji224/go_ozon/hw2/cart/internal/app"
	"github.com/Eiji224/go_ozon/hw2/cart/internal/repository"
	"github.com/Eiji224/go_ozon/hw2/cart/internal/router"
)

func main() {
	cartRepo := repository.InitCartRepository()

	router := router.InitRouter(cartRepo)
	app := app.InitApp(router)

	app.StartServer()
}
