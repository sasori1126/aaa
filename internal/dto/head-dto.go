package dto

import "mime/multipart"

type HeadTypeResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type HeadTypeRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type HeadResponse struct {
	Id               string             `json:"id"`
	Name             string             `json:"name"`
	Images           []ImageResponse    `json:"images"`
	ShortDescription string             `json:"short_description"`
	LongDescription  string             `json:"long_description"`
	Overview         string             `json:"overview"`
	HeadType         HeadTypeResponse   `json:"head_type"`
	Resources        []ResourceResponse `json:"resources"`
	Specification    Specification      `json:"specification"`
}

type UploadFileRequest struct {
	Title string `form:"title"`
	ID    string `form:"id"`
}

type HeadRequest struct {
	Name             string                  `form:"name" binding:"required"`
	Images           []*multipart.FileHeader `form:"images"`
	ShortDescription string                  `form:"short_description" binding:"required"`
	LongDescription  string                  `form:"long_description" binding:"required"`
	Overview         string                  `form:"overview" binding:"required"`
	HeadTypeId       string                  `form:"head_type_id" binding:"required"`
	Resources        []*multipart.FileHeader `form:"resources"`
	Specification    string                  `form:"specification"`
}

type Specification map[string][]map[string]string
