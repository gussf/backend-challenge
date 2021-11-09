package checkout

import (
	"testing"
	"time"

	"github.com/gussf/backend-challenge/src/repository"
)

type StubDiscountService struct{}

func (s StubDiscountService) GetDiscountForProduct(id int32) float32 {
	return 0.1
}

func TestCheckoutProcessRequest(t *testing.T) {

	tests := []struct {
		name                  string
		testProducts          []repository.ProductDAO
		testProductRequest    []ProductRequest
		expectedLength        int
		expectedTotalAmount   int
		expectedTotalDiscount int
	}{
		{
			name: "Should not checkout gift products",
			testProducts: []repository.ProductDAO{
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
			testProducts: []repository.ProductDAO{
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
			testProducts: []repository.ProductDAO{
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
			testProducts: []repository.ProductDAO{
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
			inmemoryRepo := repository.InMemoryRepository{Products: tt.testProducts}
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

func TestCheckedOutProductIsAGift(t *testing.T) {

	product := repository.ProductDAO{Is_gift: true}

	ret := CheckedOutProductIsAGift(product)
	want := true

	if want != ret {
		t.Errorf("Incorrect output: want=%t got=%t", want, ret)
	}
}

func TestItsBlackFriday(t *testing.T) {
	inmemoryRepo := repository.InMemoryRepository{}

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
		products   []repository.ProductDAO
		wantLength int
		wantCosts  int
	}{
		{
			name: "Should add one gift to checkout with no cost",
			products: []repository.ProductDAO{
				{Id: 1, Title: "a", Description: "a", Amount: 100, Is_gift: true},
			},
			wantLength: 1,
			wantCosts:  0,
		},
		{
			name: "Should not find a gift to add to checkout",
			products: []repository.ProductDAO{
				{Id: 1, Title: "a", Description: "a", Amount: 100, Is_gift: false},
			},
			wantLength: 0,
			wantCosts:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			inmemoryRepo := repository.InMemoryRepository{Products: tt.products}
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
