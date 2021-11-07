package main

import (
	"log"
	"net/http"
)

func main() {
	imr, err := NewInMemoryRepository("data/products.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	cSvc := NewCheckoutService(imr)
	r := NewECommerceRouter(cSvc)

	http.HandleFunc("/checkout", r.Checkout)

	log.Println("Starting ecommerce on :3000...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
