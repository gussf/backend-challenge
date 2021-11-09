package checkout

import (
	"log"
	"time"

	"github.com/gussf/backend-challenge/src/discount"
	"github.com/gussf/backend-challenge/src/repository"
)

type CheckoutService struct {
	repo            repository.Repository
	discountSvc     discount.DiscountService
	blackFridayDate time.Time
}

func NewCheckoutService(r repository.Repository, d discount.DiscountService, bf time.Time) CheckoutService {
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
			case repository.ErrProductNotFound:
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

func CheckedOutProductIsAGift(p repository.ProductDAO) bool {
	return p.Is_gift
}

func (c CheckoutService) ItsBlackFriday() bool {
	today := time.Now()
	return c.blackFridayDate.Month() == today.Month() && c.blackFridayDate.Day() == today.Day()
}

func (c CheckoutService) AddBlackFridayGift(r *CheckoutResponse) {
	gift, err := c.repo.FindGift()
	if err != nil {
		switch err {
		case repository.ErrNoGiftFound:
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
