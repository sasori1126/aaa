package services

import (
	"axis/ecommerce-backend/internal/domain/payments"
	"axis/ecommerce-backend/pkg/entities"
)

type PaymentService interface {
	AddOrderPayment(orderId string) *entities.ApiError
}

type DefaultPaymentService struct {
	repo payments.PaymentRepo
}

func (d DefaultPaymentService) AddOrderPayment(orderId string) *entities.ApiError {
	//TODO implement me
	panic("implement me")
}

func NewDefaultService(repo payments.PaymentRepo) PaymentService {
	return &DefaultPaymentService{repo: repo}
}
