package router

import "net/http"

func InitRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /user/{user_id}/cart", getCartContent)

	return mux
}
