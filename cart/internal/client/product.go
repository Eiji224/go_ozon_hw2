package client

import (
	"context"
	"errors"

	pb "github.com/Eiji224/go_ozon/hw2/product_service/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductClient struct {
	conn   *grpc.ClientConn
	client pb.ProductServiceClient
}

func InitProductClient(addr string) (*ProductClient, error) {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.New("failed to start product gRPC client")
	}

	return &ProductClient{conn: conn, client: pb.NewProductServiceClient(conn)}, nil
}

func (c *ProductClient) Close() error {
	return c.conn.Close()
}

func (c *ProductClient) GetProduct(ctx context.Context, sku int) (name string, price int, err error) {
	req := &pb.GetProductRequest{Sku: int64(sku)}

	resp, err := c.client.GetProduct(ctx, req)
	if err != nil {
		return "", -1, err
	}

	return resp.Name, int(resp.Price), nil
}
