package repository

import (
	"errors"
	"sync"

	"github.com/Eiji224/go_ozon/hw2/cart/internal/models"
)

type CartRepository struct {
	carts map[models.UserId]models.Cart
	mutex *sync.RWMutex
}

var ErrNotFound = errors.New("not found")
var ErrZeroOrNegativeId = errors.New("id must be greater than 0")

func InitCartRepository() *CartRepository {
	carts := make(map[models.UserId]models.Cart)

	return &CartRepository{
		carts: carts,
		mutex: &sync.RWMutex{},
	}
}

func (r *CartRepository) GetCart(userId int64) ([]*models.Product, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	id := models.UserId(userId)
	cart, ok := r.carts[id]
	if !ok {
		return nil, ErrNotFound
	}

	return cart.GetCartSlice(), nil
}

func (r *CartRepository) UpdateCart(userId int64, intSku int64, name string, count uint32, price uint32) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if userId <= 0 {
		return ErrZeroOrNegativeId
	}
	if intSku <= 0 {
		return ErrZeroOrNegativeId
	}

	cart, ok := r.carts[models.UserId(userId)]
	if !ok {
		cart = make(models.Cart)
		r.carts[models.UserId(userId)] = cart
	}

	sku := models.Sku(intSku)
	product, found := cart[sku]
	if !found {
		product = &models.Product{
			Sku:   sku,
			Name:  name,
			Count: 0,
			Price: price,
		}
		cart[sku] = product
	}
	product.Count += count

	return nil
}
