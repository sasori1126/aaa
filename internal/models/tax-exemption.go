package models

import "axis/ecommerce-backend/configs"

type TaxExemption struct {
	configs.GormModel

	TaxId  uint `gorm:"uniqueIndex:idx_tax_exemption;"`
	UserId uint `gorm:"uniqueIndex:idx_tax_exemption;"`
}
