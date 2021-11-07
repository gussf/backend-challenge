package main

import (
	"fmt"
	"testing"
)

func TestCheckoutProcessRequest(t *testing.T) {

	tests := []struct {
		name                string
		testProducts        []ProductDAO
		testProductRequest  []ProductRequest
		expectedLength      int
		expectedTotalAmount int
	}{
		{
			name: "Should not checkout gift products",
			testProducts: []ProductDAO{
				{Id: 1, Title: "a", Description: "a", Amount: 100, Is_gift: true},
			},
			testProductRequest: []ProductRequest{
				{Id: 1, Quantity: 1},
			},
			expectedLength:      0,
			expectedTotalAmount: 0,
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
			expectedLength:      2,
			expectedTotalAmount: (100 * 2) + (200 * 2),
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
			expectedLength:      2,
			expectedTotalAmount: 100 + 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inmemoryRepo := InMemoryRepository{Products: tt.testProducts}
			checkoutSvc := NewCheckoutService(inmemoryRepo)
			request := CheckoutRequest{
				Products: tt.testProductRequest,
			}

			response, _ := checkoutSvc.ProcessRequest(request)
			got := len(response.Products)
			if tt.expectedLength != got {
				t.Errorf("'%s' failed on ExpectedLength: want=%d, got=%d", tt.name, tt.expectedLength, got)
			}

			got = response.TotalAmount
			if tt.expectedTotalAmount != got {
				t.Errorf("'%s' failed on TotalAmount: want=%d, got=%d", tt.name, tt.expectedTotalAmount, got)
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

	fmt.Println(wantAmount, wantTotalDiscount, resp.TotalAmount, resp.TotalDiscount)
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
	product := ProductResponse{Id: 1, TotalAmount: 200, DiscountGiven: 100}

	response.AddProduct(product)
	want := 1
	got := len(response.Products)

	if want != got {
		t.Errorf("Incorrect Products length: want=%d, got=%d", want, got)
	}

	// AddProduct should also update checkout totals (amount and discount)

	want = product.TotalAmount
	got = response.TotalAmount
	if want != got {
		t.Errorf("Incorrect Response TotalAmount: want=%d, got=%d", want, got)
	}

	want = product.DiscountGiven
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
	discount := 0.05

	pResp := ConvertProductDAOToProductResponse(p, quantity, discount)

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

	want = int(float64(p.Amount*quantity) * discount)
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
