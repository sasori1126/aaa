package payments

import "axis/ecommerce-backend/internal/models"

type PaymentRepo interface {
	AddOrderPayment(payment *models.Payment) error
}

type CardDetail struct {
	Holder     string
	CardNumber string
	Expiry     string
	CVV        string
}

type PaymentIntegrationResponse struct {
	IsSuccessful bool    `json:"is_successful"`
	Code         string  `json:"code"`
	Message      string  `json:"message"`
	Reference    string  `json:"reference"`
	Amount       float32 `json:"amount"`
}

type PaymentImplementation interface {
	SetConfig() error
	GetCurrency() (*string, error)
	TakePayment(cardDetail *CardDetail, paymentAmount float32) (*PaymentIntegrationResponse, error)
	ValidatePaymentDetail(data *CardDetail) error
}

type Payment struct {
	Bambora PaymentImplementation
}

func NewPayment(cur string) (*Payment, error) {
	b, err := NewClient(cur)
	if err != nil {
		return nil, err
	}

	return &Payment{Bambora: b}, nil
}
