package models

import (
	"axis/ecommerce-backend/configs"
)

type Serial struct {
	configs.GormModel
	ModelId     uint
	SerialStart string
	SerialEnd   string
	YearStart   string
	YearEnd     string
}
