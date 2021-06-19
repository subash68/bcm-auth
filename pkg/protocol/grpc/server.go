package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	v1 "github.com/subash68/authenticator/pkg/api/v1"
	"google.golang.org/grpc"
)


func RunServer(ctx context.Context, v1API v1.AuthServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)

	if err != nil {
		return err
	}

	server := grpc.NewServer()
	v1.RegisterAuthServiceServer(server, v1API)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			log.Println("shutting down gRPC server...")
			server.GracefulStop()
			<-ctx.Done()
		}
	}()

	log.Println("starting gRPC server...")
	return server.Serve(listen)
}