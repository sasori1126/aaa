package customers

import "axis/ecommerce-backend/internal/models"

type CustomerRepo interface {
	CreateCustomerTempReq(cr models.CustomerTempRequest) error
	FindAllOrders() ([]models.CustomerTempRequest, error)
}