package checkout

import "github.com/gussf/backend-challenge/src/repository"

type CheckoutRequest struct {
	Products []ProductRequest
}

type ProductRequest struct {
	Id       int
	Quantity int
}

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

// AddProduct converts DAO and updates totals for safety, the caller might forget to call one or the other
func (r *CheckoutResponse) AddProduct(pDAO repository.ProductDAO, quantity int, discount float32) {

	p := ConvertProductDAOToProductResponse(pDAO, quantity, discount)
	r.Products = append(r.Products, p)
	r.UpdateCheckoutTotals(p)
}

func ConvertProductDAOToProductResponse(p repository.ProductDAO, quantity int, discount float32) ProductResponse {
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

// Gifts shouldn't cost anything
func (r *CheckoutResponse) AddGiftProduct(pDAO repository.ProductDAO, quantity int) {

	p := ConvertProductDAOToProductResponse(pDAO, quantity, 0.00)
	p.TotalAmount, p.UnitAmount, p.DiscountGiven = 0, 0, 0
	r.Products = append(r.Products, p)
}
