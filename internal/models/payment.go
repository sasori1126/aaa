package models

import "axis/ecommerce-backend/configs"

type Payment struct {
	configs.GormModel
	PaymentMethod string
	PaymentAmount float64
	Reference     string
	OrderId       uint
	Status        string
	UserId        uint
}
