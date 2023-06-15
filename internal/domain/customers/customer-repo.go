package customers

import (
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/internal/storage"
)

type CustomerRepoDb struct {
	client storage.Storage
}

func (c CustomerRepoDb) CreateCustomerTempReq(cr models.CustomerTempRequest) error {
	err := c.client.CreateCustomerTempReq(cr)
	if err != nil {
		return err
	}
	return nil
}

func (c CustomerRepoDb) FindAllOrders() ([]models.CustomerTempRequest, error) {
	return c.client.FindAllOrders()
}

func NewCustomerRepoDb(db storage.Storage) CustomerRepo {
	return CustomerRepoDb{client: db}
}
