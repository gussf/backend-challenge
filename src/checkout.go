package main

type CheckoutService struct {
	repo Repository
}

func NewCheckoutService(r Repository) CheckoutService {
	return CheckoutService{
		repo: r,
	}
}
