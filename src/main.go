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
	blackFridayDateEnvvar := os.Getenv("BLACK_FRIDAY_DATE_MMDD")

	gRPC_Deadline := time.Duration(grpcDeadlineEnvvar * int(time.Millisecond))

	blackFridayDate := Parse_MMDD_DateFromString(blackFridayDateEnvvar)

	imr, err := NewInMemoryRepository("data/products.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	dSvc := NewDiscountService_gRPC(discountGRPCAddress, gRPC_Deadline)
	cSvc := NewCheckoutService(imr, dSvc, blackFridayDate)
	r := NewECommerceRouter(cSvc)

	http.HandleFunc("/checkout", r.Checkout)

	log.Println("Starting ecommerce server on", ecommerceAddress)
	log.Println("Black friday:", blackFridayDate.Month(), blackFridayDate.Day())
	log.Fatal(http.ListenAndServe(ecommerceAddress, nil))
}

func Parse_MMDD_DateFromString(date string) time.Time {
	layout := "0102"
	blackFridayDate, err := time.Parse(layout, date)
	if err != nil {
		log.Fatalf("Failed to parse BlackFriday date (%s): %v", date, err)
	}
	return blackFridayDate
}
