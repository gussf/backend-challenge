package main

import (
	"fmt"
	"net/http"
)

type ECommerceRouter struct {
	cs CheckoutService
}

func NewECommerceRouter(cSvc CheckoutService) ECommerceRouter {
	return ECommerceRouter{
		cs: cSvc,
	}
}

func (e ECommerceRouter) Checkout(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only POST method is allowed"))
	}

	fmt.Println("Hello world")
}
