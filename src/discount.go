package main

import (
	"context"
	"log"
	"time"

	discount "github.com/gussf/backend-challenge/src/pb"
	"google.golang.org/grpc"
)

type DiscountService_gRPC struct {
	client   discount.DiscountClient
	deadline time.Duration
}

func NewDiscountService_gRPC(connAddress string, deadline time.Duration) DiscountService_gRPC {
	conn, err := grpc.DialContext(context.Background(), connAddress, grpc.WithInsecure())
	if err != nil {
		log.Printf("Could not connect to gRPC Discount Server(%s): %v", connAddress, err)
	}
	c := discount.NewDiscountClient(conn)

	return DiscountService_gRPC{
		client:   c,
		deadline: deadline,
	}
}

func (svc DiscountService_gRPC) GetDiscountForProduct(id int32) float32 {

	clientDeadline := time.Now().Add(svc.deadline)
	ctx, _ := context.WithDeadline(context.Background(), clientDeadline)

	r, err := svc.client.GetDiscount(ctx, &discount.GetDiscountRequest{ProductID: id})
	if err != nil {
		log.Printf("Failed to get discount for product=%d, returning discount=0.00: %v", id, err)
		return 0.00
	}

	log.Printf("Discount=%.2f received for product=%d", r.GetPercentage(), id)
	return r.GetPercentage()
}

type DiscountService interface {
	GetDiscountForProduct(id int32) float32
}
