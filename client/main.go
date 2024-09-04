package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(":50051", opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	c := order.NewOrderClient(conn)
	c.Create(context.Background(), &order.CreateOrderRequest{})
}
