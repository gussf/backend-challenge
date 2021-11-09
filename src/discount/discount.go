package discount

type DiscountService interface {
	GetDiscountForProduct(id int32) float32
}
