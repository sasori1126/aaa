package models

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/dto"
)

type Model struct {
	configs.GormModel
	Name           string `gorm:"unique"`
	Description    string
	Manufacturer   Manufacturer
	ManufacturerID uint
	Controllers    []*Controller `gorm:"many2many:model_controllers;"`
	Serials        []*Serial     `gorm:"many2many:model_serials;"`
	ImagePath      string
	Status         string
}

func (m Model) ToResponse() dto.ModelResponse {
	hid, _ := EncodeHashId(m.ID)
	mid, _ := EncodeHashId(m.Manufacturer.ID)
	mf := dto.EmbedManufacturer{
		Id:          mid,
		Name:        m.Manufacturer.Name,
		Description: m.Manufacturer.Description,
	}

	var serials []dto.SerialResponse
	serials = []dto.SerialResponse{}
	for _, serial := range m.Serials {
		id, _ := EncodeHashId(serial.ID)
		sr := dto.SerialResponse{
			Id:          id,
			SerialStart: serial.SerialStart,
			SerialEnd:   serial.SerialEnd,
			YearStart:   serial.YearStart,
			YearEnd:     serial.YearEnd,
		}

		serials = append(serials, sr)
	}

	var controllers []dto.ControllerResponse
	controllers = []dto.ControllerResponse{}
	for _, coll := range m.Controllers {
		id, _ := EncodeHashId(coll.ID)
		c := dto.ControllerResponse{
			Id:          id,
			Name:        coll.Name,
			Description: coll.Description,
		}

		controllers = append(controllers, c)
	}

	return dto.ModelResponse{
		Id:           hid,
		Name:         m.Name,
		Image:        m.ImagePath,
		Serials:      serials,
		Controllers:  controllers,
		Manufacturer: mf,
	}
}

type ModelController struct {
	ModelId      uint
	ControllerId uint
}
