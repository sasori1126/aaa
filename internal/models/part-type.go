package models

import (
	"axis/ecommerce-backend/configs"
)

type PartType struct {
	configs.GormModel
	Name        string
	Description string
}
