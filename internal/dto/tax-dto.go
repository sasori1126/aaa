package dto

// Add Tax request.
type TaxRequest struct {
	Country     string  `json:"country" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Rate        float64 `json:"rate" binding:"required"`
	State       string  `json:"state" binding:"required"`
}

type GetTaxRequest struct {
	City    string
	Country string
	Region  string
	State   string
	ZipCode string
}

type TaxLine struct {
	City               string              `json:"city"`
	Country            string              `json:"country"`
	DeliveryTaxAmounts []DeliveryTaxAmount `json:"delivery_tax_amounts"`
	Id                 string              `json:"id"`
	Name               string              `json:"name"`
	ProductTaxAmount   float64             `json:"product_tax_amount"`
	Rate               float64             `json:"rate"`
	State              string              `json:"state"`
	ZipCode            string              `json:"zipcode"`
}

type DeliveryTaxAmount struct {
	DeliveryRateId    string  `json:"delivery_rate_id"`
	DeliveryTaxAmount float64 `json:"delivery_tax_amount"`
}
