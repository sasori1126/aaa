package models

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/dto"
)

type Distributor struct {
	configs.GormModel
	Name      string
	Type      string
	Address   Address
	Contact   Contact
	Site      string
	AddressID uint
	ContactID uint
}

func (d Distributor) ToResponse() dto.DistributorResponse {
	hid, _ := EncodeHashId(d.ID)
	address := d.Address.ToResponse()
	contact := d.Contact.ToResponse()
	return dto.DistributorResponse{
		ID:      hid,
		Name:    d.Name,
		Site:    d.Site,
		Address: address,
		Contact: contact,
	}
}
