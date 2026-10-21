package server

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Eiji224/go_ozon/hw2/product_service/internal/repository"
	pb "github.com/Eiji224/go_ozon/hw2/product_service/pkg/api"
	"github.com/stretchr/testify/require"
)

const testProducts = `{"1":{"name":"Book","price":100},"2":{"name":"Pen","price":5}}`

func writeFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "products.json")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func TestGetProduct(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		req         *pb.GetProductRequest
		resp        *pb.GetProductResponse
		expectedErr error
	}{
		{"successful response", &pb.GetProductRequest{Sku: 1}, &pb.GetProductResponse{Sku: 1, Name: "Book", Price: 100}, nil},
		{"zero sku", &pb.GetProductRequest{Sku: 0}, nil, ErrInvalidSku},
		{"negative sku", &pb.GetProductRequest{Sku: -1}, nil, ErrInvalidSku},
		{"not found", &pb.GetProductRequest{Sku: 11}, nil, ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, _ := repository.InitProductRepository(writeFile(t, testProducts))
			server := InitProductServer(repo)

			got, err := server.GetProduct(t.Context(), tt.req)

			require.Equal(t, tt.resp, got)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}
