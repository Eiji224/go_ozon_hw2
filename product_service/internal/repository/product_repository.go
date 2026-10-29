package repository

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/Eiji224/go_ozon/hw2/product_service/internal/models"
)

type ProductRepository struct {
	products map[models.Sku]*models.Product
	mutex    *sync.RWMutex
}

var ErrNotFound = errors.New("not found")

func InitProductRepository(fileName string) (*ProductRepository, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, errors.New("failed to open products.json")
	}
	defer file.Close()

	products := make(map[models.Sku]*models.Product)
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&products); err != nil {
		return nil, errors.New("failed to parse products from products.json")
	}

	return &ProductRepository{products: products, mutex: &sync.RWMutex{}}, nil
}

func (r *ProductRepository) GetProduct(sku models.Sku) (*models.Product, error) {
	r.mutex.RLock()
	product, ok := r.products[sku]
	r.mutex.RUnlock()
	if !ok {
		return nil, ErrNotFound
	}

	return product, nil
}
