package dto

import (
	"time"
)

type OrderResponse struct {
	BillingAddress      UserAddressResponse `json:"billing_address"`
	CancelledAt         *time.Time          `json:"cancelled_at"`
	CancelReason        string              `json:"cancel_reason"`
	CartId              string              `json:"cart_id"`
	Currency            string              `json:"currency"`
	DeliveryInstruction string              `json:"delivery_instruction"`
	DeliveryMethod      string              `json:"delivery_method"`
	DeliveryRateId      string              `json:"delivery_rate_id"`
	Id                  string              `json:"id"`
	Items               []OrderItemResponse `json:"items"`
	OrderDate           time.Time           `json:"order_date"`
	PaidAmount          float64             `json:"paid_amount"`
	PaymentMethod       string              `json:"payment_method"`
	PickupPoint         string              `json:"pickup_point"`
	PurchaseOrderNumber string              `json:"purchase_order_number"`
	ReturnedAt          *time.Time          `json:"returned_at"`
	ShippedAt           *time.Time          `json:"shipped_at"`
	ShippingAddress     UserAddressResponse `json:"shipping_address"`
	ShippingAmount      float64             `json:"shipping_amount"`
	Status              string              `json:"status"`
	SubTotalAmount      float64             `json:"sub_total_amount"`
	Taxes               []OrderTaxResponse  `json:"taxes"`
	TotalAmount         float64             `json:"total_amount"`
	TotalTaxAmount      float64             `json:"total_tax_amount"`
	TrackingNumber      string              `json:"tracking_number"`
	User                EmbedUser           `json:"user"`
	UserId              string              `json:"user_id"`
}

type OrderItemResponse struct {
	Description  string    `json:"description"`
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	Note         string    `json:"note"`
	OrderId      string    `json:"order_id"`
	Part         EmbedPart `json:"part"`
	PricePerUnit float64   `json:"price_per_unit"`
	Quantity     float64   `json:"quantity"`
	Total        float64   `json:"total"`
	Unit         string    `json:"unit"`
	Weight       float64   `json:"weight"`
}

type OrderTaxResponse struct {
	Amount            float64 `json:"amount"`
	Description       string  `json:"description"`
	Id                string  `json:"id"`
	Name              string  `json:"name"`
	OrderId           string  `json:"order_id"`
	ProductTaxAmount  float64 `json:"product_tax_amount"`
	Rate              float64 `json:"rate"`
	ShippingTaxAmount float64 `json:"shipping_tax_amount"`
}

type AddOrderPayment struct {
	Amount        float64 `json:"amount"`
	OrderId       string  `json:"order_id"`
	PaymentDate   string  `json:"payment_date"`
	PaymentMethod string  `json:"payment_method"  binding:"required,oneof='bank' 'cheque' 'onAccount'"`
	Reference     string  `json:"reference"`
}

type OrderUpdateStatusRequest struct {
	Status string `json:"status"  binding:"required,oneof='pending' 'paid' 'declined' 'delivered' 'shipped'"`
}
