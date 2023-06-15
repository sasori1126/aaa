package models

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/dto"
)

type Controller struct {
	configs.GormModel
	Name        string
	Description string
}

func (m Controller) ToResponse() dto.ControllerResponse {
	hid, _ := EncodeHashId(m.ID)
	return dto.ControllerResponse{
		Id:          hid,
		Name:        m.Name,
		Description: m.Description,
	}
}
