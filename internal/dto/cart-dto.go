package dto

type CartResponse struct {
	Id          string             `json:"id"`
	User        EmbedUser          `json:"user"`
	Items       []CartItemResponse `json:"items"`
	Status      string             `json:"status"`
	SubTotal    float64            `json:"sub_total"`
	TotalAmount float64            `json:"total_amount"`
}

type CartItemRemoveRequest struct {
	CartItemsId []string `json:"cart_items_id"`
}

type GetCartRequest struct {
	CartId string `json:"cart_id"`
}

type UpdateCartItemQuantityRequest struct {
	CartItemId string `json:"cart_item_id"`
	Quantity   int    `json:"quantity"`
}

type CartItemRequest struct {
	PartId   string  `json:"part_id" binding:"required"`
	Quantity float64 `json:"quantity" binding:"required"`
	Note     string  `json:"note"`
}

type CartItemResponse struct {
	Id           string       `json:"id"`
	CartId       string       `json:"cart_id"`
	PartId       string       `json:"part_id"`
	PricePerUnit float64      `json:"price_per_unit"`
	Quantity     float64      `json:"quantity"`
	Unit         string       `json:"unit"`
	Description  string       `json:"description"`
	Note         string       `json:"note"`
	Title        string       `json:"title"`
	Amount       float64      `json:"amount"`
	Part         PartResponse `json:"part"`
}

type DeliveryRateResponse struct {
	Amount           float64  `json:"amount"`
	Attributes       []string `json:"attributes"`
	CarrierAccount   string   `json:"carrier_account"`
	Currency         string   `json:"currency"`
	CurrencyLocal    string   `json:"currency_local"`
	Days             int      `json:"days"`
	DurationTerms    string   `json:"duration_terms"`
	Expected         string   `json:"expected"`
	Id               string   `json:"id"`
	Provider         string   `json:"provider"`
	ProviderImage200 string   `json:"provider_image_200"`
	ProviderImage75  string   `json:"provider_image_75"`
	Zone             string   `json:"zone"`
}

type DeliveryTaxLine struct {
	TaxAmount float64 `json:"tax_amount"`
	TaxName   float64 `json:"tax_name"`
	TaxRate   float64 `json:"tax_rate"`
}

type PlaceOrderRequest struct {
	Card                     CardDetailRequest `json:"card"`
	CartId                   string            `json:"cart_id" binding:"required"`
	ConvergePayTransactionId string            `json:"convergepay_transaction_id"`
	DeliveryAddressId        string            `json:"delivery_address_id"`
	DeliveryRateId           string            `json:"delivery_rate_id"`
	DeliveryInstruction      string            `json:"delivery_instruction"`
	DeliveryMethod           string            `json:"delivery_method" binding:"required"`
	PaymentMethod            string            `json:"payment_method" binding:"required"`
	PickupPoint              string            `json:"pickup_point"`
	PurchaseOrderNumber      string            `json:"purchase_order_number"`
}

type CardDetailRequest struct {
	Holder     string `json:"holder"`
	CardNumber string `json:"card_number"`
	Expiry     string `json:"expiry"`
	CVV        string `json:"CVV"`
}
