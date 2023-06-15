package dto

type ResourceResponse struct {
	Id           string `json:"id"`
	Title        string `json:"title"`
	Type         string `json:"type"`
	ResourcePath string `json:"resource_path"`
}
