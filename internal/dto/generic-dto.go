package dto

type ShippoParcel struct {
	Length       string
	Width        string
	Height       string
	DistanceUnit string
	Weight       string
	MassUnit     string
}

type SupportRequest struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Reason      string `json:"reason" binding:"required"`
	Message     string `json:"message" binding:"required"`
}
