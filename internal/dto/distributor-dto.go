package dto

type DistributorRequest struct {
	Name    string         `json:"name,omitempty" binding:"required"`
	Address AddressRequest `json:"address,omitempty" binding:"required"`
	Site    string         `json:"site,omitempty" binding:"required"`
	Contact ContactRequest `json:"contact,omitempty" binding:"required"`
}

type DistributorResponse struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Address AddressResponse `json:"address"`
	Site    string          `json:"website"`
	Contact ContactResponse `json:"contact"`
}
