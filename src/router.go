package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gussf/backend-challenge/src/checkout"
)

type CheckoutJSONResponse struct {
	Total_amount               int                   `json:"total_amount"`
	Total_amount_with_discount int                   `json:"total_amount_with_discount"`
	Total_discount             int                   `json:"total_discount"`
	Products                   []ProductJSONResponse `json:"products"`
}

type ProductJSONResponse struct {
	Id           int  `json:"id"`
	Quantity     int  `json:"quantity"`
	Unit_amount  int  `json:"unit_amount"`
	Total_amount int  `json:"total_amount"`
	Discount     int  `json:"discount"`
	Is_gift      bool `json:"is_gift"`
}

type ECommerceRouter struct {
	checkoutSvc checkout.CheckoutService
}

func NewECommerceRouter(cs checkout.CheckoutService) ECommerceRouter {
	return ECommerceRouter{
		checkoutSvc: cs,
	}
}

func (router ECommerceRouter) Checkout(w http.ResponseWriter, r *http.Request) {

	enc := json.NewEncoder(w)

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only POST method is allowed"))
		return
	}

	checkoutReq, err := ParseCheckoutRequestFromBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse request: " + err.Error()))
		log.Println("Failed to parse request: " + err.Error())
		return
	}

	if checkoutReq.HasNoProducts() {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Request must have at least one product"))
		return
	}

	resp := router.checkoutSvc.ProcessRequest(checkoutReq)
	w.Header().Add("Content-Type", "application/json")
	jsonResp := ConvertCheckoutResponseToCheckoutJSONResponse(resp)
	enc.Encode(jsonResp)
}

func ParseCheckoutRequestFromBody(r *http.Request) (checkout.CheckoutRequest, error) {

	var checkoutReq checkout.CheckoutRequest

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&checkoutReq)
	if err != nil {
		return checkoutReq, err
	}

	return checkoutReq, nil
}

func ConvertCheckoutResponseToCheckoutJSONResponse(r *checkout.CheckoutResponse) CheckoutJSONResponse {

	resp := CheckoutJSONResponse{Products: make([]ProductJSONResponse, 0)}

	for _, p := range r.Products {
		resp.Products = append(resp.Products, ConvertProductResponseToProductJSONResponse(p))
	}

	resp.Total_amount = r.TotalAmount
	resp.Total_amount_with_discount = r.TotalAmount - r.TotalDiscount
	resp.Total_discount = r.TotalDiscount

	return resp
}

func ConvertProductResponseToProductJSONResponse(p checkout.ProductResponse) ProductJSONResponse {
	return ProductJSONResponse{
		Id:           p.Id,
		Quantity:     p.Quantity,
		Unit_amount:  p.UnitAmount,
		Total_amount: p.TotalAmount,
		Discount:     p.DiscountGiven,
		Is_gift:      p.IsGift,
	}
}
