package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Eiji224/go_ozon/hw2/cart/internal/models"
)

type ProductClient interface {
	GetProduct(ctx context.Context, sku int64) (name string, price uint32, err error)
}

type CartRepository struct {
	carts         map[models.UserId]models.Cart
	productClient ProductClient
	mutex         *sync.RWMutex
}

var ErrNotFound = errors.New("not found")
var ErrZeroOrNegativeId = errors.New("id must be greater than 0")

func InitCartRepository(productClient ProductClient) *CartRepository {
	carts := make(map[models.UserId]models.Cart)

	return &CartRepository{
		carts:         carts,
		mutex:         &sync.RWMutex{},
		productClient: productClient,
	}
}

func (r *CartRepository) GetCart(ctx context.Context, userId int64) ([]*models.Product, error) {
	r.mutex.RLock()
	id := models.UserId(userId)
	cart, ok := r.carts[id]
	r.mutex.RUnlock()

	if !ok {
		return nil, ErrNotFound
	}

	result := make([]*models.Product, len(cart))
	skus := cart.GetSkuSlice()
	errs := make([]error, len(cart))

	wgroup := sync.WaitGroup{}
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	for i, sku := range skus {
		wgroup.Add(1)
		go func() {
			defer wgroup.Done()
			name, price, err := r.productClient.GetProduct(ctxTimeout, int64(sku))
			if err != nil {
				errs[sku] = err
				return
			}

			result[i] = &models.Product{
				Name:  name,
				Price: price,
				Count: cart[sku],
				Sku:   sku,
			}
		}()
	}
	wgroup.Wait()

	for _, e := range errs {
		if e != nil {
			return nil, e
		}
	}

	models.SortProducts(result)
	return result, nil
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
	_, found := cart[sku]
	if !found {
		cart[sku] = 0
	}
	cart[sku]++

	return nil
}
