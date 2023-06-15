package models

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/dto"
)

type Contact struct {
	configs.GormModel
	Name         string
	Email        string
	Phone        string
	Organisation string
	Fax          string
	Postal       string
}

func (c Contact) ToResponse() dto.ContactResponse {
	hid, _ := EncodeHashId(c.ID)
	return dto.ContactResponse{
		Id:     hid,
		Name:   c.Name,
		Email:  c.Email,
		Phone:  c.Phone,
		Fax:    c.Fax,
		Postal: c.Postal,
	}
}
