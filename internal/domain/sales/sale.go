package sales

import "axis/ecommerce-backend/internal/models"

type SaleRepo interface {
	StoreKeslaOrder(data models.KeslaOrder) error
	StoreControllerOrder(data models.ControllerOrder) error
	StoreAxisHeadOrder(data models.AxisHead) error
}
