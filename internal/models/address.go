package models

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/dto"
)

type Address struct {
	configs.GormModel
	City       string
	Country    string
	Lat        float64
	Lon        float64
	Province   string // Deprecated.
	State      string // State for US, and Province for CA.
	StreetName string
	ZipCode    string
}

func (a Address) ToResponse() dto.AddressResponse {
	hid, _ := EncodeHashId(a.ID)
	return dto.AddressResponse{
		ID:         hid,
		City:       a.City,
		StreetName: a.StreetName,
		Country:    a.Country,
		Province:   a.Province,
		State:      a.State,
		ZipCode:    a.ZipCode,
		Lat:        a.Lat,
		Lon:        a.Lon,
	}
}
