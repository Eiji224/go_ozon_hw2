package router

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Eiji224/go_ozon/hw2/cart/internal/dto"
	"github.com/Eiji224/go_ozon/hw2/cart/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrUserIdLessThanZero = "user_id must be digit and greater than 0"
var ErrSkuLessThanZero = "sku_id must be digit and greater than 0"
var ErrCountDigitGreaterThanZero = "count must be digit and greater than 0"
var ErrProductNotFound = "product not found"
var ErrCartNotFound = "cart not found"

type CartRepo interface {
	GetCart(ctx context.Context, userId int64) ([]*models.Product, error)
	UpdateCart(userId int64, intSku int64, name string, count uint32, price uint32) error
}

type ProductClient interface {
	GetProduct(ctx context.Context, sku int64) (name string, price uint32, err error)
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
	userId, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil || userId <= 0 {
		writeError(w, ErrUserIdLessThanZero, http.StatusBadRequest)
		return
	}
	cart, err := router.cartRepo.GetCart(r.Context(), userId)
	if err != nil {
		writeError(w, ErrCartNotFound, http.StatusNotFound)
		return
	}

	var totalPrice uint32
	for _, item := range cart {
		totalPrice += item.Price * item.Count
	}

	response := dto.GetCartResponse{
		Items:      cart,
		TotalPrice: totalPrice,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (router *Router) addProductToCart(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	userId, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil || userId <= 0 {
		writeError(w, ErrUserIdLessThanZero, http.StatusBadRequest)
		return
	}
	skuId, err := strconv.ParseInt(r.PathValue("sku_id"), 10, 64)
	if err != nil || skuId <= 0 {
		writeError(w, ErrSkuLessThanZero, http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		writeError(w, ErrCountDigitGreaterThanZero, http.StatusBadRequest)
		return
	}

	name, price, err := router.productClient.GetProduct(context.Background(), skuId)
	if err != nil {
		st, errParsed := status.FromError(err)
		if errParsed {
			switch st.Code() {
			case codes.InvalidArgument:
				writeError(w, st.Message(), http.StatusBadRequest)
				return
			case codes.NotFound:
				writeError(w, ErrProductNotFound, http.StatusPreconditionFailed)
				return
			default:
				writeError(w, "internal server error", http.StatusInternalServerError)
				return
			}
		}

		writeError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := router.cartRepo.UpdateCart(userId, skuId, name, product.Count, price); err != nil {
		writeError(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func writeError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
