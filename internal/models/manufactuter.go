package models

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/dto"
)

type Manufacturer struct {
	configs.GormModel
	Name        string
	Description string
}

func (m Manufacturer) ToResponse() dto.ManufacturerResponse {
	hid, _ := EncodeHashId(m.ID)
	return dto.ManufacturerResponse{
		Id:          hid,
		Name:        m.Name,
		Description: m.Description,
	}
}
