package main

import (
	"log"
	"time"
)

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
	repo            Repository
	discountSvc     DiscountService
	blackFridayDate time.Time
}

func NewCheckoutService(r Repository, d DiscountService, bf time.Time) CheckoutService {
	return CheckoutService{
		repo:            r,
		discountSvc:     d,
		blackFridayDate: bf,
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

		discount := c.discountSvc.GetDiscountForProduct(int32(p.Id))
		productResp := ConvertProductDAOToProductResponse(productDAO, p.Quantity, discount)
		response.AddProduct(productResp)
	}

	if c.ItsBlackFriday() {
		log.Println("Its black friday! Attempt to add gift product to checkout")
		c.AddGiftProductToCheckout(&response)
	}

	return &response, nil
}

func CheckedOutProductIsAGift(p ProductDAO) bool {
	return p.Is_gift
}

func ConvertProductDAOToProductResponse(p ProductDAO, quantity int, discount float32) ProductResponse {
	return ProductResponse{
		Id:            p.Id,
		Quantity:      quantity,
		UnitAmount:    p.Amount,
		TotalAmount:   p.Amount * quantity,
		DiscountGiven: int(float32(p.Amount*quantity) * discount),
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

func (c CheckoutService) ItsBlackFriday() bool {
	today := time.Now()
	return c.blackFridayDate.Month() == today.Month() && c.blackFridayDate.Day() == today.Day()
}

func (c CheckoutService) AddGiftProductToCheckout(r *CheckoutResponse) {
	gift, err := c.repo.FindGift()
	if err != nil {
		log.Printf("Could not add a gift product to checkout: %v", err)
		return
	}
	productResp := ConvertProductDAOToProductResponse(gift, 1, 0)
	r.AddGiftProduct(productResp)
	log.Printf("Gift product=%d added to checkout", productResp.Id)
	log.Println(r)
}

// Gifts shouldn't cost anything
func (r *CheckoutResponse) AddGiftProduct(p ProductResponse) {
	p.TotalAmount = 0
	p.UnitAmount = 0
	p.DiscountGiven = 0
	r.Products = append(r.Products, p)
}
