package grpc

import (
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"homework10/internal/app"
	"homework10/internal/ports/grpc/loggers"
	"homework10/internal/ports/grpc/proto"
	"log"
	"net"
)

func NewGRPCServer(port string, a app.App) (*grpc.Server, net.Listener) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			loggers.Logger, loggers.PanicInterceptor)))
	grpcClient := NewService(a)
	proto.RegisterAdServiceServer(grpcServer, grpcClient)
	return grpcServer, lis
}

func TestNewGRPCServer(sizeBuff int, a app.App) (*grpc.Server, *bufconn.Listener) {
	lis := bufconn.Listen(sizeBuff)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			loggers.Logger, loggers.PanicInterceptor)))
	grpcClient := NewService(a)
	proto.RegisterAdServiceServer(grpcServer, grpcClient)
	return grpcServer, lis
}
