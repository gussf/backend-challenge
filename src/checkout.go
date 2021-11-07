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

func (c CheckoutService) ProcessRequest(req CheckoutRequest) (*CheckoutResponse, error) {
	response := CheckoutResponse{}

	for _, p := range req.Products {
		productDAO, err := c.repo.Find(p.Id)
		if err != nil {
			log.Println("Product with id =", p.Id, "not found")
			continue
		}

		if CheckedOutProductIsAGift(productDAO) {
			log.Println("Product with id =", p.Id, "is a gift and therefore cannot be checked out")
			continue
		}

		// get discount from grpc server
		// todo
		discount := 0.05

		productResp := ConvertProductDAOToProductResponse(productDAO, p.Quantity, discount)
		response.AddProduct(productResp)
	}
	return &response, nil
}

func CheckedOutProductIsAGift(p ProductDAO) bool {
	return p.Is_gift
}

func ConvertProductDAOToProductResponse(p ProductDAO, quantity int, discount float64) ProductResponse {
	return ProductResponse{
		Id:            p.Id,
		Quantity:      quantity,
		UnitAmount:    p.Amount,
		TotalAmount:   p.Amount * quantity,
		DiscountGiven: int(float64(p.Amount*quantity) * discount),
		IsGift:        p.Is_gift,
	}
}

// AddProduct updates totals for safety, the caller might forget to call one or the other
func (r *CheckoutResponse) AddProduct(p ProductResponse) {
	r.Products = append(r.Products, p)
	r.UpdateCheckoutTotals(p)
}

func (r *CheckoutResponse) UpdateCheckoutTotals(p ProductResponse) {
	r.TotalAmount += p.TotalAmount
	r.TotalDiscount += p.DiscountGiven
}
