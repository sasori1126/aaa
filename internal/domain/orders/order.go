package orders

import "axis/ecommerce-backend/internal/models"

type OrderRepo interface {
	CreateOrder(order models.Order) (*models.Order, error)
	GetUserOrders(user uint) ([]models.Order, error)
	GetOrders(limit, offset int, userId uint) ([]models.Order, error)
	GetOrderField(fv models.FindByField) (*models.Order, error)
	UpdateOrderField(orderId uint, fv models.FindByField) error
}
