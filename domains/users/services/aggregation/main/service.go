package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	pbAuth "github.com/kirvader/BodyController/domains/users/services/aggregation/proto"
)

type UsersService struct {
	authClient *pbAuth.AuthClient

	pbAuth.UnimplementedAuthServer
}

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	authClient := pbAuth.NewAuthClient(conn)

	svc := &UsersService{
		authClient: authClient,
	}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()

	pbAuth.RegisterUsersServer(grpcServer, svc)
	reflection.Register(grpcServer)
	log.Printf("server listening at %v", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
