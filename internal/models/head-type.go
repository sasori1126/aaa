package models

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/dto"
)

type HeadType struct {
	configs.GormModel
	Name        string
	Description string
}

func (t HeadType) ToResponse() dto.HeadTypeResponse {
	hid, _ := EncodeHashId(t.ID)
	return dto.HeadTypeResponse{
		Id:          hid,
		Name:        t.Name,
		Description: t.Description,
	}
}
