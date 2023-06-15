package payments

import (
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/internal/storage"
)

type PaymentRepoDb struct {
	client storage.Storage
}

func (p PaymentRepoDb) AddOrderPayment(payment *models.Payment) error {
	return p.client.AddOrderPayment(payment)
}

func NewPaymentRepoDb(db storage.Storage) PaymentRepo {
	return &PaymentRepoDb{client: db}
}
