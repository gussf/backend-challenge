package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type CheckoutRequest struct {
	Products []ProductRequest
}

type ProductRequest struct {
	id       int
	quantity int
}

type ECommerceRouter struct {
	cs CheckoutService
}

func NewECommerceRouter(cSvc CheckoutService) ECommerceRouter {
	return ECommerceRouter{
		cs: cSvc,
	}
}

func (router ECommerceRouter) Checkout(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only POST method is allowed"))
		return
	}

	checkoutReq, err := ParseProductsFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse request: " + err.Error()))
		log.Println("Failed to parse request: " + err.Error())
	}

	log.Println(checkoutReq)
}

func ParseProductsFromRequest(r *http.Request) (CheckoutRequest, error) {

	var checkoutReq CheckoutRequest

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&checkoutReq)
	if err != nil {
		return checkoutReq, err
	}

	return checkoutReq, nil
}
