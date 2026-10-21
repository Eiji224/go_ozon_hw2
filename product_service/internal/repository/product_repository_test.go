package repository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Eiji224/go_ozon/hw2/product_service/internal/models"
)

const testProducts = `{"1":{"name":"Book","price":100},"2":{"name":"Pen","price":5}}`

func writeFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "products.json")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func TestInitProductRepository(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		prepare func(t *testing.T) string
		wantErr bool
	}{
		{"ok", func(t *testing.T) string { return writeFile(t, testProducts) }, false},
		{"file not found", func(t *testing.T) string { return "no_such.json" }, true},
		{"invalid json", func(t *testing.T) string { return writeFile(t, "{asddsa") }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, err := InitProductRepository(tt.prepare(t))

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, repo)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, repo)
		})
	}
}

func TestGetProduct(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		sku     models.Sku
		want    *models.Product
		wantErr error
	}{
		{"found", 1, &models.Product{Name: "Book", Price: 100}, nil},
		{"not found", 123, nil, ErrNotFound},
		{"zero sku", 0, nil, ErrNotFound},
		{"negative sku", -1, nil, ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, err := InitProductRepository(writeFile(t, testProducts))
			require.NoError(t, err)

			got, err := repo.GetProduct(tt.sku)

			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
		})
	}
}
