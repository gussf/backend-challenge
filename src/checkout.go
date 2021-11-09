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

func (c CheckoutService) ProcessRequest(req CheckoutRequest) *CheckoutResponse {
	response := &CheckoutResponse{}

	for _, p := range req.Products {
		productDAO, err := c.repo.Find(p.Id)
		if err != nil {
			switch err {
			case ErrProductNotFound:
				log.Printf("Product with id=%d not found in repository", p.Id)
				continue
			default:
				log.Printf("Something unexpected went wrong obtaining product=%d: %v", p.Id, err)
				continue
			}
		}

		if CheckedOutProductIsAGift(productDAO) {
			log.Printf("Product with id=%d is a gift and therefore cannot be checked out", p.Id)
			continue
		}

		discount := c.discountSvc.GetDiscountForProduct(int32(p.Id))
		response.AddProduct(productDAO, p.Quantity, discount)
	}

	// Only add gift if there are products in checkout
	if len(response.Products) > 0 && c.ItsBlackFriday() {
		log.Println("Its black friday! Attempt to add gift product to checkout")
		c.AddBlackFridayGift(response)
	}

	return response
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
		switch err {
		case ErrNoGiftFound:
			log.Printf("No gift was found in repository")
			return
		default:
			log.Printf("Something went wrong obtaining a gift: %v", err)
			return
		}
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
