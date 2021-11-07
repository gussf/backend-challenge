package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	imr, err := NewInMemoryRepository("data/products.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	cSvc := NewCheckoutService(imr)
	r := NewECommerceRouter(cSvc)

	http.HandleFunc("/checkout", r.Checkout)

	address := os.Getenv("LISTEN_ADDRESS")

	log.Println("Starting ecommerce server on", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
