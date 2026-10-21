package repository

import (
	"testing"

	"github.com/Eiji224/go_ozon/hw2/cart/internal/models"
	"github.com/stretchr/testify/require"
)

func initTestCartRepo(carts map[models.UserId]models.Cart) *CartRepository {
	return &CartRepository{carts}
}

func TestGetCart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		userId   int64
		initCart map[models.UserId]models.Cart
		wantErr  error
		wantCart models.Cart
	}{
		{"successful get", 1, map[models.UserId]models.Cart{1: {1: &models.Product{Name: "Book", Price: 100, Count: 1}}}, nil, models.Cart{1: &models.Product{Name: "Book", Price: 100, Count: 1}}},
		{"failed get", 1, map[models.UserId]models.Cart{}, ErrNotFound, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := initTestCartRepo(tt.initCart)
			got, err := repo.GetCart(tt.userId)

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
		expectedCart models.Cart
		expectedErr  error
	}{
		{"success update", 1, 1, "Book", 1, 100, models.Cart{1: &models.Product{Sku: 1, Name: "Book", Price: 100, Count: 1}}, nil},
		{"zero userId", 0, 1, "Book", 1, 100, nil, ErrZeroOrNegativeId},
		{"negative userId", -1, 1, "Book", 1, 100, nil, ErrZeroOrNegativeId},
		{"zero sku", 1, 0, "Book", 1, 100, nil, ErrZeroOrNegativeId},
		{"negative sku", 1, -1, "Book", 1, 100, nil, ErrZeroOrNegativeId},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := initTestCartRepo(make(map[models.UserId]models.Cart))
			err := repo.UpdateCart(tt.userId, tt.intSku, tt.productName, tt.count, tt.price)

			require.Equal(t, tt.expectedErr, err)

			actual, _ := repo.GetCart(tt.userId)
			require.Equal(t, tt.expectedCart, actual)
		})
	}
}
