package taxes

import "axis/ecommerce-backend/internal/models"

type TaxRepo interface {
	AddTax(data models.Tax) error
	DeleteTaxExemption(query models.TaxExemption) error
	GetTaxesByAddress(userId uint, query models.Tax) ([]models.Tax, []models.Tax, error)
	GetTaxesForTaxExemptions(userId uint) ([]models.Tax, []models.Tax, error)
	SaveTaxExemption(query models.TaxExemption) error
}
