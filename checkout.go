package main

type CheckoutRequest struct {
	Products []Product
}

type Product struct {
	id       int64
	quantity int64
}
