package main

import "log"

type CheckoutResponse struct {
	TotalAmount   int
	TotalDiscount int
	Products      []ProductResponse
}

type ProductResponse struct {
	Id            int
	Quantity      int
	UnitAmount    int
	TotalAmount   int
	DiscountGiven int
	IsGift        bool
}

type CheckoutRequest struct {
	Products []ProductRequest
}

type ProductRequest struct {
	Id       int
	Quantity int
}

type CheckoutService struct {
	repo Repository
}

func NewCheckoutService(r Repository) CheckoutService {
	return CheckoutService{
		repo: r,
	}
}

func (c CheckoutService) ProcessRequest(req CheckoutRequest) (CheckoutResponse, error) {
	var response CheckoutResponse

	for _, p := range req.Products {
		productDAO, err := c.repo.Find(p.Id)
		if err != nil {
			log.Println("Product with id =", p.Id, "not found")
			continue
		}

		if productDAO.Is_gift {
			log.Println("Product with id =", p.Id, "is a gift and therefore cannot be checked out")
			continue
		}

		// get discount from grpc server
		// todo
		discount := 0

		response.TotalAmount += productDAO.Amount * p.Quantity
		response.TotalDiscount += discount
		response.Products = append(response.Products, ConvertProductDAOToProductResponse(productDAO, p.Quantity, discount))
	}
	return response, nil
}

func ConvertProductDAOToProductResponse(p ProductDAO, quantity int, discount int) ProductResponse {
	return ProductResponse{
		Id:            p.Id,
		Quantity:      quantity,
		UnitAmount:    p.Amount,
		TotalAmount:   p.Amount * quantity,
		DiscountGiven: discount,
		IsGift:        p.Is_gift,
	}
}
