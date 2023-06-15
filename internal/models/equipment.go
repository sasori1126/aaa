package models

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/dto"
)

type Equipment struct {
	configs.GormModel
	UserId   uint
	User     User
	ModelId  uint
	Model    Model
	SerialId uint
	Serial   Serial
}

func (e Equipment) ToResponse() dto.EquipmentResponse {
	id, _ := EncodeHashId(e.ID)
	mid, _ := EncodeHashId(e.ModelId)
	sid, _ := EncodeHashId(e.SerialId)
	m := e.Model
	mfid, _ := EncodeHashId(m.ManufacturerID)
	mf := dto.EmbedManufacturer{
		Id:          mfid,
		Name:        m.Manufacturer.Name,
		Description: m.Manufacturer.Description,
	}
	mm := struct {
		Id           string
		Name         string
		Image        string
		Manufacturer dto.EmbedManufacturer
		Serial       dto.SerialResponse
	}{
		Id:           mid,
		Name:         m.Name,
		Image:        m.ImagePath,
		Manufacturer: mf,
		Serial: dto.SerialResponse{
			Id:          sid,
			SerialStart: e.Serial.SerialStart,
			SerialEnd:   e.Serial.SerialEnd,
			YearStart:   e.Serial.YearStart,
			YearEnd:     e.Serial.YearEnd,
		},
	}

	return dto.EquipmentResponse{
		Id: id,
		Model: struct {
			Id           string                `json:"id"`
			Name         string                `json:"name"`
			Image        string                `json:"image"`
			Manufacturer dto.EmbedManufacturer `json:"manufacturer"`
			Serial       dto.SerialResponse    `json:"serial"`
		}(mm),
	}
}
