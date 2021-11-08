package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {

	ecommerceAddress := os.Getenv("ECOMMERCE_LISTEN_ADDRESS")
	discountGRPCAddress := os.Getenv("DISCOUNT_GRPC_ADDRESS")
	grpcDeadlineEnvvar, _ := strconv.Atoi(os.Getenv("GRPC_DEADLINE_MS"))
	log.Println(grpcDeadlineEnvvar)
	gRPC_Deadline := time.Duration(grpcDeadlineEnvvar * int(time.Millisecond))

	imr, err := NewInMemoryRepository("data/products.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	dSvc := NewDiscountService_gRPC(discountGRPCAddress, gRPC_Deadline)
	cSvc := NewCheckoutService(imr, dSvc)
	r := NewECommerceRouter(cSvc)

	http.HandleFunc("/checkout", r.Checkout)

	log.Println("Starting ecommerce server on", ecommerceAddress)
	log.Fatal(http.ListenAndServe(ecommerceAddress, nil))
}
