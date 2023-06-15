package entities

type ApiResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Links   E           `json:"links"`
	Meta    E           `json:"meta"`
}
