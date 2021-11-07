package main

import (
	"encoding/json"
	"log"
	"net/http"
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
	cs CheckoutService
}

func NewECommerceRouter(cSvc CheckoutService) ECommerceRouter {
	return ECommerceRouter{
		cs: cSvc,
	}
}

func (router ECommerceRouter) Checkout(w http.ResponseWriter, r *http.Request) {

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
	}

	resp, err := router.cs.ProcessRequest(checkoutReq)

	jsonResp := ConvertCheckoutResponseToCheckoutJSONResponse(resp)
	enc := json.NewEncoder(w)
	enc.Encode(jsonResp)
}

func ParseCheckoutRequestFromBody(r *http.Request) (CheckoutRequest, error) {

	var checkoutReq CheckoutRequest

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&checkoutReq)
	if err != nil {
		return checkoutReq, err
	}

	return checkoutReq, nil
}

func ConvertCheckoutResponseToCheckoutJSONResponse(r CheckoutResponse) CheckoutJSONResponse {

	var resp CheckoutJSONResponse

	for _, p := range r.Products {
		resp.Products = append(resp.Products, ConvertProductResponseToProductJSONResponse(p))
	}

	resp.Total_amount = r.TotalAmount
	resp.Total_amount_with_discount = r.TotalAmount - r.TotalDiscount
	resp.Total_discount = r.TotalDiscount

	return resp
}

func ConvertProductResponseToProductJSONResponse(p ProductResponse) ProductJSONResponse {
	return ProductJSONResponse{
		Id:           p.Id,
		Quantity:     p.Quantity,
		Unit_amount:  p.UnitAmount,
		Total_amount: p.TotalAmount,
		Discount:     p.DiscountGiven,
		Is_gift:      p.IsGift,
	}
}
