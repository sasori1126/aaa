package models

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/dto"
)

type Resource struct {
	configs.GormModel
	Title        string
	Type         string
	Tag          []Tag
	ResourcePath string
}

type Tag struct {
	Name       string
	ResourceId uint
}

func (r Resource) ToResponse() dto.ResourceResponse {
	hid, _ := EncodeHashId(r.ID)
	return dto.ResourceResponse{
		Id:           hid,
		Title:        r.Title,
		Type:         r.Type,
		ResourcePath: r.ResourcePath,
	}
}
