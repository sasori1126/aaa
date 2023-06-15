package models

import "axis/ecommerce-backend/configs"

type UserEquipment struct {
	configs.GormModel
	User    User
	UserId  uint
	Model   Model
	ModelId uint
	Name    string
	Serial  string
	Unit    string
	Year    string
}
