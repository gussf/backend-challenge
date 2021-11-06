package main

type CheckoutRequest struct {
	Products []ProductRequest
}

type ProductRequest struct {
	id       int64
	quantity int64
}
