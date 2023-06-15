package dto

type ModelResponse struct {
	Id           string               `json:"id"`
	Name         string               `json:"name"`
	Image        string               `json:"image"`
	Manufacturer EmbedManufacturer    `json:"manufacturer"`
	Serials      []SerialResponse     `json:"serials"`
	Controllers  []ControllerResponse `json:"controllers"`
}

type ModelRequest struct {
	ManufacturerID string          `json:"manufacturer_id" binding:"required"`
	ControllerIds  []string        `json:"controller_ids" binding:"required"`
	Name           string          `json:"name" binding:"required"`
	Description    string          `json:"description" binding:"required"`
	Serials        []SerialRequest `json:"serials" binding:"required"`
}

type EmbedModel struct {
	Id           string            `json:"id"`
	Name         string            `json:"name"`
	Image        string            `json:"image"`
	Manufacturer EmbedManufacturer `json:"manufacturer"`
	Series       []string          `json:"series"`
}
