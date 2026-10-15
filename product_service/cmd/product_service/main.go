package main

import (
	"log"
	"net"

	"github.com/Eiji224/go_ozon/hw2/product_service/internal/repository"
	sv "github.com/Eiji224/go_ozon/hw2/product_service/internal/server"
	pb "github.com/Eiji224/go_ozon/hw2/product_service/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	reflection.Register(server)

	rep, err := repository.InitProductRepository()
	if err != nil {
		log.Fatalf("failed to start: %v", err)
	}
	ps := sv.InitProductServer(rep)
	pb.RegisterProductServiceServer(server, ps)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
