package dto

type AddressRequest struct {
	City       string  `json:"city,omitempty" binding:"required"`
	StreetName string  `json:"street_name,omitempty" binding:"required"`
	Country    string  `json:"country,omitempty" binding:"required"`
	Province   string  `json:"province,omitempty" binding:"required"`
	State      string  `json:"state,omitempty" binding:"required"`
	ZipCode    string  `json:"zip_code,omitempty" binding:"required"`
	Lat        float64 `json:"lat,omitempty" binding:"numeric"`
	Lon        float64 `json:"lon,omitempty" binding:"numeric"`
}

type AddressResponse struct {
	ID         string  `json:"id"`
	City       string  `json:"city,omitempty"`
	StreetName string  `json:"street_name,omitempty"`
	Country    string  `json:"country,omitempty"`
	Province   string  `json:"province,omitempty"`
	State      string  `json:"state,omitempty"`
	ZipCode    string  `json:"zip_code,omitempty"`
	Lat        float64 `json:"lat,omitempty"`
	Lon        float64 `json:"lon,omitempty"`
}
