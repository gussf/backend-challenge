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
		response.AddProduct(productDAO, p.Quantity, discount)
	}

	if c.ItsBlackFriday() {
		log.Println("Its black friday! Attempt to add gift product to checkout")
		c.AddBlackFridayGift(&response)
	}

	return &response, nil
}

func CheckedOutProductIsAGift(p ProductDAO) bool {
	return p.Is_gift
}

// AddProduct converts DAO and updates totals for safety, the caller might forget to call one or the other
func (r *CheckoutResponse) AddProduct(pDAO ProductDAO, quantity int, discount float32) {

	p := ConvertProductDAOToProductResponse(pDAO, quantity, discount)
	r.Products = append(r.Products, p)
	r.UpdateCheckoutTotals(p)
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

func (r *CheckoutResponse) UpdateCheckoutTotals(p ProductResponse) {
	r.TotalAmount += p.TotalAmount
	r.TotalDiscount += p.DiscountGiven
}

func (c CheckoutService) ItsBlackFriday() bool {
	today := time.Now()
	return c.blackFridayDate.Month() == today.Month() && c.blackFridayDate.Day() == today.Day()
}

func (c CheckoutService) AddBlackFridayGift(r *CheckoutResponse) {
	gift, err := c.repo.FindGift()
	if err != nil {
		log.Printf("Could not add a gift product to checkout: %v", err)
		return
	}
	r.AddGiftProduct(gift, 1)
	log.Printf("Gift product=%d added to checkout", gift.Id)
}

// Gifts shouldn't cost anything
func (r *CheckoutResponse) AddGiftProduct(pDAO ProductDAO, quantity int) {

	p := ConvertProductDAOToProductResponse(pDAO, quantity, 0.00)
	p.TotalAmount, p.UnitAmount, p.DiscountGiven = 0, 0, 0
	r.Products = append(r.Products, p)
}
