package sales

import (
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/internal/storage"
)

type SaleRepoDb struct {
	client storage.Storage
}

func (s SaleRepoDb) StoreKeslaOrder(data models.KeslaOrder) error {
	return s.client.StoreKeslaOrder(data)
}

func (s SaleRepoDb) StoreControllerOrder(data models.ControllerOrder) error {
	return s.client.StoreControllerOrder(data)
}

func (s SaleRepoDb) StoreAxisHeadOrder(data models.AxisHead) error {
	return s.client.StoreAxisHeadOrder(data)
}

func NewSaleRepoDb(db storage.Storage) SaleRepo {
	return &SaleRepoDb{client: db}
}
