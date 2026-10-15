package repository

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/Eiji224/go_ozon/hw2/product_service/internal/models"
)

type ProductRepository struct {
	products map[models.Sku]*models.Product
}

func InitProductRepository() (*ProductRepository, error) {
	file, err := os.Open("products.json")
	if err != nil {
		return nil, errors.New("failed to open products.json")
	}
	defer file.Close()

	products := make(map[models.Sku]*models.Product)
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&products); err != nil {
		return nil, errors.New("failed to parse products from products.json")
	}

	return &ProductRepository{products: products}, nil
}

func (r *ProductRepository) GetProduct(sku models.Sku) (*models.Product, error) {
	product, ok := r.products[sku]
	if !ok {
		return nil, errors.New("Not found")
	}

	return product, nil
}
