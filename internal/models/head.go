package models

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/dto"
)

type Head struct {
	configs.GormModel
	Name             string
	ShortDescription string
	LongDescription  string
	Overview         string
	HeadType         HeadType
	HeadTypeID       uint
	Resources        []Resource `gorm:"many2many:head_resources;"`
	Images           []Image    `gorm:"many2many:head_images;"`
	Specification    string
}

func (h Head) ToResponse() dto.HeadResponse {
	hid, _ := EncodeHashId(h.ID)
	ht := h.HeadType.ToResponse()

	var images []dto.ImageResponse
	images = []dto.ImageResponse{}
	for _, image := range h.Images {
		img := image.ToResponse()
		images = append(images, img)
	}

	var rrs []dto.ResourceResponse
	rrs = []dto.ResourceResponse{}
	for _, resource := range h.Resources {
		rr := resource.ToResponse()
		rrs = append(rrs, rr)
	}

	return dto.HeadResponse{
		Id:               hid,
		Name:             h.Name,
		ShortDescription: h.ShortDescription,
		Overview:         h.Overview,
		LongDescription:  h.LongDescription,
		HeadType:         ht,
		Resources:        rrs,
		Images:           images,
		Specification:    dto.Specification{},
	}
}
