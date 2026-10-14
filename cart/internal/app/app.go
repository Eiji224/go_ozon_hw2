package app

import (
	"log"
	"net/http"
)

type App struct {
	router *http.ServeMux
}

func InitApp(router *http.ServeMux) *App {
	return &App{
		router: router,
	}
}

func (a *App) StartServer() {
	log.Fatal(http.ListenAndServe(":8080", a.router))
}
