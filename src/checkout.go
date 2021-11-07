package main

type CheckoutRequest struct {
	Products []ProductRequest
}

type ProductRequest struct {
	id       int
	quantity int
}

type CheckoutService struct {
	repo Repository
}

func NewCheckoutService(r Repository) CheckoutService {
	return CheckoutService{
		repo: r,
	}
}
