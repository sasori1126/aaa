package models

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/dto"
)

type Image struct {
	configs.GormModel
	Name        string
	Description string
	Ext         string
	Path        string
	Parts       []Part `gorm:"many2many:part_images;"`
	Heads       []Head `gorm:"many2many:head_images;"`
}

func (i Image) ToResponse() dto.ImageResponse {
	hid, _ := EncodeHashId(i.ID)
	return dto.ImageResponse{
		Id:   hid,
		Name: i.Name,
		Path: i.Path,
	}
}
