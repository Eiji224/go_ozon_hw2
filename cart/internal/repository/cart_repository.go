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

func (r *CartRepository) GetCart(userId int) (models.Cart, error) {
	id := models.UserId(userId)
	cart, ok := r.carts[id]
	if !ok {
		return nil, errors.New("Not found")
	}

	return cart, nil
}

func (r *CartRepository) UpdateCart(userId int, intSku int, name string, count int, price int) error {
	cart, err := r.GetCart(userId)
	if err != nil {
		cart = make(models.Cart)
		r.carts[models.UserId(userId)] = cart
	}

	sku := models.Sku(intSku)
	cart[sku] = &models.Product{
		Sku:   sku,
		Name:  name,
		Count: count,
		Price: price,
	}

	return nil
}
