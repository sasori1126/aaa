package taxes

import (
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/internal/storage"
)

type TaxRepoDb struct {
	client storage.Storage
}

func (t TaxRepoDb) AddTax(data models.Tax) error {
	return t.client.CreateTax(data)
}

func (t TaxRepoDb) DeleteTaxExemption(query models.TaxExemption) error {
	return t.client.DeleteTaxExemption(query)
}

func (t TaxRepoDb) GetTaxesByAddress(userId uint, query models.Tax) ([]models.Tax, []models.Tax, error) {
	return t.client.GetTaxesByAddress(userId, query)
}

func (t TaxRepoDb) GetTaxesForTaxExemptions(userId uint) ([]models.Tax, []models.Tax, error) {
	return t.client.GetTaxesForTaxExemptions(userId)
}

func (t TaxRepoDb) SaveTaxExemption(query models.TaxExemption) error {
	return t.client.SaveTaxExemption(query)
}

func NewTaxRepoDb(db storage.Storage) TaxRepo {
	return &TaxRepoDb{client: db}
}
