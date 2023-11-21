package main

import (
	"context"
	"log"

	pbAuth "github.com/kirvader/BodyController/domains/users/services/base/auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pbAuth.NewAuthClient(conn)

	resp, err := client.CreateUser(context.Background(), &pbAuth.CreateUserRequest{
		UserCredentials: &pbAuth.UserCredentials{
			Username: "kk",
			Password: "lol",
		},
	})
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	log.Printf("response got: %v", resp)
}