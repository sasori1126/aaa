package entities

type UpdateFields struct {
	Field string
	Value interface{}
}

type QueryPathParam struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type RequestPathData struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

var (
	Video = "video"
	Image = "image"
	Pdf   = "pdf"
)
