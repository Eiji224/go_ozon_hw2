package repository

import (
	"errors"

	"github.com/Eiji224/go_ozon/hw2/cart/internal/models"
)

type CartRepository struct {
	carts map[models.UserId]models.Cart
}

func InitCartRepository() *CartRepository {
	carts := make(map[models.UserId]models.Cart)

	return &CartRepository{
		carts: carts,
	}
}

func (r *CartRepository) GetCart(userId int64) (models.Cart, error) {
	id := models.UserId(userId)
	cart, ok := r.carts[id]
	if !ok {
		return nil, errors.New("Not found")
	}

	return cart, nil
}

func (r *CartRepository) UpdateCart(userId int64, intSku int64, name string, count uint32, price uint32) error {
	cart, err := r.GetCart(userId)
	if err != nil {
		cart = make(models.Cart)
		r.carts[models.UserId(userId)] = cart
	}

	sku := models.Sku(intSku)
	product, ok := cart[sku]
	if !ok {
		cart[sku] = &models.Product{
			Sku:   sku,
			Name:  name,
			Count: 0,
			Price: price,
		}
	}
	product.Count += count

	return nil
}
