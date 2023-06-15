package models

import "axis/ecommerce-backend/configs"

type Oem struct {
	configs.GormModel
	PartNumber     string
	OemNote        string
	ManufacturerId uint
	Manufacturer
}
