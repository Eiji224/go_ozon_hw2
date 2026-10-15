package main

import (
	"log"

	"github.com/Eiji224/go_ozon/hw2/cart/internal/app"
	"github.com/Eiji224/go_ozon/hw2/cart/internal/client"
	"github.com/Eiji224/go_ozon/hw2/cart/internal/repository"
	"github.com/Eiji224/go_ozon/hw2/cart/internal/router"
)

func main() {
	productClient, err := client.InitProductClient("localhost:50051")
	if err != nil {
		log.Fatal(err)
	}
	defer productClient.Close()

	cartRepo := repository.InitCartRepository()

	router := router.InitRouter(cartRepo, productClient)
	app := app.InitApp(router)

	app.StartServer()
}
