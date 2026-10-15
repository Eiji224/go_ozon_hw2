package server

import (
	"context"

	"github.com/Eiji224/go_ozon/hw2/product_service/internal/models"
	pb "github.com/Eiji224/go_ozon/hw2/product_service/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductRepo interface {
	GetProduct(sku models.Sku) (*models.Product, error)
}

type ProductServer struct {
	pb.ProductServiceServer
	productRepo ProductRepo
}

func InitProductServer(productRepo ProductRepo) *ProductServer {
	return &ProductServer{productRepo: productRepo}
}

func (s *ProductServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	if req.Sku <= 0 {
		return nil, status.Error(codes.InvalidArgument, "sku must be greater than 0")
	}

	product, err := s.productRepo.GetProduct(models.Sku(req.Sku))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "product %d not found", req.Sku)
	}

	response := &pb.GetProductResponse{
		Sku:   req.Sku,
		Name:  product.Name,
		Price: product.Price,
	}

	return response, nil
}
