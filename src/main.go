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

	port := os.Getenv("LISTEN_PORT")
	address := "0.0.0.0:" + port

	log.Println("Starting ecommerce server on", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
