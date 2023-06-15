package models

import (
	"axis/ecommerce-backend/configs"
)

type Tax struct {
	configs.GormModel

	Country               string `gorm:"uniqueIndex:idx_tax;"`
	Description           string
	IsAllowedTaxExemption bool
	Name                  string `gorm:"uniqueIndex:idx_tax;"`
	State                 string `gorm:"uniqueIndex:idx_tax;"` // State for US, and Province for CA.
	Rate                  float64
}
