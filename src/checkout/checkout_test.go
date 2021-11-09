package checkout

import (
	"testing"

	"github.com/gussf/backend-challenge/src/repository"
)

func TestUpdateCheckoutTotals(t *testing.T) {

	resp := &CheckoutResponse{}
	product := ProductResponse{TotalAmount: 300, DiscountGiven: 50}

	wantAmount := product.TotalAmount
	wantTotalDiscount := product.DiscountGiven

	resp.UpdateCheckoutTotals(product)
	if wantAmount != resp.TotalAmount {
		t.Errorf("Incorrect TotalAmount: want=%d got=%d", wantAmount, resp.TotalAmount)
	}

	if wantTotalDiscount != resp.TotalDiscount {
		t.Errorf("Incorrect TotalDiscount: want=%d got=%d", wantTotalDiscount, resp.TotalDiscount)
	}
}

func TestAddProduct(t *testing.T) {

	response := CheckoutResponse{}
	pDAO := repository.ProductDAO{Id: 1, Amount: 200}
	quantity := 1
	var discount float32 = 0.10

	response.AddProduct(pDAO, quantity, float32(discount))
	want := 1
	got := len(response.Products)

	if want != got {
		t.Errorf("Incorrect Products length: want=%d, got=%d", want, got)
	}

	// AddProduct should also update checkout totals (amount and discount)

	want = pDAO.Amount * quantity
	got = response.TotalAmount
	if want != got {
		t.Errorf("Incorrect Response TotalAmount: want=%d, got=%d", want, got)
	}

	want = int(float32(pDAO.Amount*quantity) * discount)
	got = response.TotalDiscount
	if want != got {
		t.Errorf("Incorrect Response TotalDiscount: want=%d, got=%d", want, got)
	}
}

func TestConvertProductDAOToProductResponse(t *testing.T) {
	p := repository.ProductDAO{
		Id: 1, Title: "a", Description: "a", Amount: 200, Is_gift: false,
	}
	quantity := 2
	var discount float32 = 0.05

	pResp := ConvertProductDAOToProductResponse(p, quantity, float32(discount))

	want := p.Id
	got := pResp.Id
	if want != got {
		t.Errorf("Incorrect Product ID: want=%d, got=%d", want, got)
	}

	want = p.Amount
	got = pResp.UnitAmount
	if want != got {
		t.Errorf("Incorrect Product UnitAmount: want=%d, got=%d", want, got)
	}

	want = quantity
	got = pResp.Quantity
	if want != got {
		t.Errorf("Incorrect Product Quantity: want=%d, got=%d", want, got)
	}

	want = p.Amount * quantity
	got = pResp.TotalAmount
	if want != got {
		t.Errorf("Incorrect Product TotalAmount: want=%d, got=%d", want, got)
	}

	want = int(float32(p.Amount*quantity) * discount)
	got = pResp.DiscountGiven
	if want != got {
		t.Errorf("Incorrect Product DiscountGiven: want=%d, got=%d", want, got)
	}

	wantB := p.Is_gift
	gotB := pResp.IsGift
	if wantB != gotB {
		t.Errorf("Incorrect Product IsGift: want=%t, got=%t", wantB, gotB)
	}
}

func TestAddGiftProduct(t *testing.T) {

	product := repository.ProductDAO{Id: 1, Title: "a", Description: "a", Amount: 100, Is_gift: true}
	quantity := 1
	wantLength := 1
	wantCosts := 0

	checkoutResp := &CheckoutResponse{}
	checkoutResp.AddGiftProduct(product, quantity)

	got := len(checkoutResp.Products)
	if wantLength != got {
		t.Errorf("Incorrect Response Products Length: want=%d, got=%d", wantLength, got)
	}

	addedGift := checkoutResp.Products[0]

	got = addedGift.TotalAmount
	if wantCosts != got {
		t.Errorf("Incorrect Gift Product TotalAmount: want=%d, got=%d", wantCosts, got)
	}

	got = addedGift.DiscountGiven
	if wantCosts != got {
		t.Errorf("Incorrect Gift Product DiscountGiven: want=%d, got=%d", wantCosts, got)
	}

	got = addedGift.UnitAmount
	if wantCosts != got {
		t.Errorf("Incorrect Gift Product UnitAmount: want=%d, got=%d", wantCosts, got)
	}

	got = addedGift.Quantity
	if quantity != got {
		t.Errorf("Incorrect Gift Product Quantity: want=%d, got=%d", quantity, got)
	}

	gotB := addedGift.IsGift
	want := true
	if want != gotB {
		t.Errorf("Incorrect Gift Product IsGift: want=%t, got=%t", want, gotB)
	}
}
