package main

import (
	"testing"
	"time"
)

type StubDiscountService struct{}

func (s StubDiscountService) GetDiscountForProduct(id int32) float32 {
	return 0.1
}

func TestCheckoutProcessRequest(t *testing.T) {

	tests := []struct {
		name                  string
		testProducts          []ProductDAO
		testProductRequest    []ProductRequest
		expectedLength        int
		expectedTotalAmount   int
		expectedTotalDiscount int
	}{
		{
			name: "Should not checkout gift products",
			testProducts: []ProductDAO{
				{Id: 1, Title: "a", Description: "a", Amount: 100, Is_gift: true},
			},
			testProductRequest: []ProductRequest{
				{Id: 1, Quantity: 1},
			},
			expectedLength:        0,
			expectedTotalAmount:   0,
			expectedTotalDiscount: 0,
		},
		{
			name: "Checkout valid products",
			testProducts: []ProductDAO{
				{Id: 1, Title: "a", Description: "a", Amount: 100, Is_gift: false},
				{Id: 2, Title: "b", Description: "b", Amount: 200, Is_gift: false},
			},
			testProductRequest: []ProductRequest{
				{Id: 1, Quantity: 2},
				{Id: 2, Quantity: 2},
			},
			expectedLength:        2,
			expectedTotalAmount:   600,
			expectedTotalDiscount: 60,
		},
		{
			name: "Checkout valid products with one gift product(gift shouldnt be checked out)",
			testProducts: []ProductDAO{
				{Id: 1, Title: "a", Description: "a", Amount: 100, Is_gift: false},
				{Id: 2, Title: "b", Description: "b", Amount: 200, Is_gift: false},
				{Id: 3, Title: "c", Description: "c", Amount: 50, Is_gift: true},
			},
			testProductRequest: []ProductRequest{
				{Id: 1, Quantity: 1},
				{Id: 2, Quantity: 1},
				{Id: 3, Quantity: 1},
			},
			expectedLength:        2,
			expectedTotalAmount:   300,
			expectedTotalDiscount: 30,
		},
		{
			name: "Checkout product that doesnt exist in repository",
			testProducts: []ProductDAO{
				{Id: 1, Title: "a", Description: "a", Amount: 100, Is_gift: false},
			},
			testProductRequest: []ProductRequest{
				{Id: 4, Quantity: 1},
			},
			expectedLength:        0,
			expectedTotalAmount:   0,
			expectedTotalDiscount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inmemoryRepo := InMemoryRepository{Products: tt.testProducts}
			checkoutSvc := NewCheckoutService(inmemoryRepo, StubDiscountService{}, time.Now().Add(24*time.Hour)) // Add 1 day to avoid Black Friday
			request := CheckoutRequest{
				Products: tt.testProductRequest,
			}

			response := checkoutSvc.ProcessRequest(request)
			got := len(response.Products)
			if tt.expectedLength != got {
				t.Errorf("'%s' Incorrect ExpectedLength: want=%d, got=%d", tt.name, tt.expectedLength, got)
			}

			got = response.TotalAmount
			if tt.expectedTotalAmount != got {
				t.Errorf("'%s' Incorrect TotalAmount: want=%d, got=%d", tt.name, tt.expectedTotalAmount, got)
			}

			got = response.TotalDiscount
			if tt.expectedTotalDiscount != got {
				t.Errorf("'%s' Incorrect TotalDiscount: want=%d, got=%d", tt.name, tt.expectedTotalDiscount, got)
			}
		})
	}
}

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

func TestCheckedOutProductIsAGift(t *testing.T) {

	product := ProductDAO{Is_gift: true}

	ret := CheckedOutProductIsAGift(product)
	want := true

	if want != ret {
		t.Errorf("Incorrect output: want=%t got=%t", want, ret)
	}
}

func TestAddProduct(t *testing.T) {

	response := CheckoutResponse{}
	pDAO := ProductDAO{Id: 1, Amount: 200}
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
	p := ProductDAO{
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

func TestItsBlackFriday(t *testing.T) {
	inmemoryRepo := InMemoryRepository{}

	tests := []struct {
		name            string
		checkoutService CheckoutService
		want            bool
	}{
		{
			name:            "Should be black friday",
			checkoutService: NewCheckoutService(inmemoryRepo, StubDiscountService{}, time.Now()),
			want:            true,
		},
		{
			name:            "Should Not be black friday",
			checkoutService: NewCheckoutService(inmemoryRepo, StubDiscountService{}, time.Now().Add(24*time.Hour)),
			want:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.checkoutService.ItsBlackFriday()
			if tt.want != got {
				t.Errorf("%s: Incorrect black friday assertion: want=%t, got=%t", tt.name, tt.want, got)
			}
		})
	}
}

func TestAddBlackFridayGift(t *testing.T) {

	tests := []struct {
		name       string
		products   []ProductDAO
		wantLength int
		wantCosts  int
	}{
		{
			name: "Should add one gift to checkout with no cost",
			products: []ProductDAO{
				{Id: 1, Title: "a", Description: "a", Amount: 100, Is_gift: true},
			},
			wantLength: 1,
			wantCosts:  0,
		},
		{
			name: "Should not find a gift to add to checkout",
			products: []ProductDAO{
				{Id: 1, Title: "a", Description: "a", Amount: 100, Is_gift: false},
			},
			wantLength: 0,
			wantCosts:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			inmemoryRepo := InMemoryRepository{Products: tt.products}
			checkoutSvc := NewCheckoutService(inmemoryRepo, StubDiscountService{}, time.Now())
			checkoutResp := &CheckoutResponse{}

			checkoutSvc.AddBlackFridayGift(checkoutResp)
			got := len(checkoutResp.Products)
			if tt.wantLength != got {
				t.Errorf("%s: Incorrect Products Length: want=%d, got=%d", tt.name, tt.wantLength, got)
			}

			got = checkoutResp.TotalAmount
			if tt.wantCosts != got {
				t.Errorf("%s: Incorrect TotalAmount: want=%d, got=%d", tt.name, tt.wantLength, got)
			}
		})
	}

}

func TestAddGiftProduct(t *testing.T) {

	product := ProductDAO{Id: 1, Title: "a", Description: "a", Amount: 100, Is_gift: true}
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
