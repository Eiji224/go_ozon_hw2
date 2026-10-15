package router

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Eiji224/go_ozon/hw2/cart/internal/dto"
	"github.com/Eiji224/go_ozon/hw2/cart/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CartRepo interface {
	GetCart(userId int) (models.Cart, error)
	UpdateCart(userId int, intSku int, name string, count int, price int) error
}

type ProductClient interface {
	GetProduct(ctx context.Context, sku int) (name string, price int, err error)
}

type Router struct {
	r             *http.ServeMux
	cartRepo      CartRepo
	productClient ProductClient
}

func InitRouter(cartRepo CartRepo, productClient ProductClient) *Router {
	router := &Router{
		r:             http.NewServeMux(),
		cartRepo:      cartRepo,
		productClient: productClient,
	}

	router.r.HandleFunc("GET /user/{user_id}/cart", router.getCartContent)
	router.r.HandleFunc("POST /user/{user_id}/cart/{sku_id}", router.addProductToCart)

	return router
}

func (router *Router) Serve() {
	log.Fatal(http.ListenAndServe(":8080", router.r))
}

func (router *Router) getCartContent(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.PathValue("user_id"))
	if err != nil {
		http.Error(w, `{"error": "invalid response"}`, http.StatusBadRequest)
		return
	}
	cart, err := router.cartRepo.GetCart(userId)
	if err != nil {
		http.Error(w, `{"error": "cart not found"}`, http.StatusNotFound)
		return
	}
	items := cart.GetCartSlice()

	var totalPrice int
	for _, item := range items {
		totalPrice += item.Price * item.Count
	}

	response := dto.GetCartResponse{
		Items:      items,
		TotalPrice: totalPrice,
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (router *Router) addProductToCart(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	userId, err := strconv.Atoi(r.PathValue("user_id"))
	if err != nil {
		http.Error(w, `{"error": "invalid response"}`, http.StatusBadRequest)
		return
	}
	skuId, err := strconv.Atoi(r.PathValue("sku_id"))
	if err != nil {
		http.Error(w, `{"error": "invalid response"}`, http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, `{"error": "invalid response"}`, http.StatusBadRequest)
		return
	}

	name, price, err := router.productClient.GetProduct(context.Background(), skuId)
	status, errParsed := status.FromError(err)
	if errParsed {
		switch status.Code() {
		case codes.InvalidArgument:
			http.Error(w, `{"error": "`+status.Message()+`"}`, http.StatusBadRequest)
			return
		case codes.NotFound:
			http.Error(w, `{"error": "`+status.Message()+`"}`, http.StatusNotFound)
			return
		default:
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
			fmt.Printf("internal server error: %v", err)
			return
		}
	}

	if err := router.cartRepo.UpdateCart(userId, skuId, name, product.Count, price); err != nil {
		fmt.Printf("internal server error: %v", err)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}
}
