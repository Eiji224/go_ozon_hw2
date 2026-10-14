package main

import (
	"github.com/Eiji224/go_ozon/hw2/cart/internal/app"
	"github.com/Eiji224/go_ozon/hw2/cart/internal/router"
)

func main() {
	router := router.InitRouter()
	app := app.InitApp(router)

	app.StartServer()
}
