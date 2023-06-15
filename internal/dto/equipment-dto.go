package dto

type EquipmentRequest struct {
	ModelId  string `json:"model_id"`
	SerialId string `json:"serial_id"`
}

type EquipmentResponse struct {
	Id    string
	Model struct {
		Id           string            `json:"id"`
		Name         string            `json:"name"`
		Image        string            `json:"image"`
		Manufacturer EmbedManufacturer `json:"manufacturer"`
		Serial       SerialResponse    `json:"serial"`
	} `json:"model"`
}
