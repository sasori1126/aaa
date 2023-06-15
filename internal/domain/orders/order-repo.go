package orders

import (
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/internal/storage"
)

type OrderRepoDb struct {
	client storage.Storage
}

func (o OrderRepoDb) UpdateOrderField(orderId uint, fv models.FindByField) error {
	return o.client.UpdateOrderField(orderId, fv)
}

func (o OrderRepoDb) GetOrderField(fv models.FindByField) (*models.Order, error) {
	return o.client.GetOrderByField(fv)
}

func (o OrderRepoDb) GetOrders(limit, offset int, userId uint) ([]models.Order, error) {
	return o.client.GetOrders(limit, offset, userId)
}

func (o OrderRepoDb) GetUserOrders(user uint) ([]models.Order, error) {
	return o.client.GetUserOrders(user)
}

func (o OrderRepoDb) CreateOrder(order models.Order) (*models.Order, error) {
	return o.client.CreateOrder(order)
}

func NewOrderRepoDb(db storage.Storage) OrderRepo {
	return &OrderRepoDb{client: db}
}
