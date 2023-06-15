package dto

type ImageRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Ext         string `json:"ext"`
	Path        string `json:"path"`
}

type ImageResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}
