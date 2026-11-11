package repository

import (
	"context"
	"sync"
	"testing"

	"github.com/Eiji224/go_ozon/hw2/cart/internal/models"
	"github.com/stretchr/testify/require"
)

type FakeProductService struct {
	expectedName   string
	expectedPrice  uint32
	exptectedError error
}

func (s *FakeProductService) GetProduct(ctx context.Context, sku int64) (string, uint32, error) {
	return s.expectedName, s.expectedPrice, s.exptectedError
}

func initTestCartRepo(carts map[models.UserId]models.Cart, expectedProduct *models.Product) *CartRepository {
	if len(carts) == 0 && expectedProduct == nil {
		return &CartRepository{carts, &FakeProductService{"", 0, ErrNotFound}, &sync.RWMutex{}}
	}

	return &CartRepository{carts, &FakeProductService{expectedProduct.Name, expectedProduct.Price, nil}, &sync.RWMutex{}}
}

func TestGetCart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		userId      int64
		initCart    map[models.UserId]models.Cart
		initProduct *models.Product
		wantErr     error
		wantCart    []*models.Product
	}{
		{"successful get", 1, map[models.UserId]models.Cart{1: {1: 1}}, &models.Product{Sku: 1, Name: "Book", Price: 100, Count: 1}, nil, []*models.Product{{Sku: 1, Name: "Book", Price: 100, Count: 1}}},
		{"failed get", 1, map[models.UserId]models.Cart{}, nil, ErrNotFound, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := initTestCartRepo(tt.initCart, tt.initProduct)
			got, err := repo.GetCart(context.Background(), tt.userId)

			require.Equal(t, tt.wantCart, got)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestUpdateCart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		userId       int64
		intSku       int64
		productName  string
		count        uint32
		price        uint32
		initProduct  *models.Product
		expectedCart []*models.Product
		expectedErr  error
	}{
		{"success update", 1, 1, "Book", 1, 100, &models.Product{Sku: 1, Name: "Book", Price: 100, Count: 1}, []*models.Product{{Sku: 1, Name: "Book", Price: 100, Count: 1}}, nil},
		{"zero userId", 0, 1, "Book", 1, 100, nil, nil, ErrZeroOrNegativeId},
		{"negative userId", -1, 1, "Book", 1, 100, nil, nil, ErrZeroOrNegativeId},
		{"zero sku", 1, 0, "Book", 1, 100, nil, nil, ErrZeroOrNegativeId},
		{"negative sku", 1, -1, "Book", 1, 100, nil, nil, ErrZeroOrNegativeId},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := initTestCartRepo(make(map[models.UserId]models.Cart), tt.initProduct)
			err := repo.UpdateCart(tt.userId, tt.intSku, tt.productName, tt.count, tt.price)

			require.Equal(t, tt.expectedErr, err)

			actual, _ := repo.GetCart(context.Background(), tt.userId)
			require.Equal(t, tt.expectedCart, actual)
		})
	}
}
