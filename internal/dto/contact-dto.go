package dto

type ContactRequest struct {
	Name   string `json:"name,omitempty" binding:"required"`
	Email  string `json:"email,omitempty" binding:"required"`
	Phone  string `json:"phone,omitempty" binding:"required"`
	Fax    string `json:"fax,omitempty"`
	Postal string `json:"postal,omitempty"`
}

type ContactResponse struct {
	Id     string `json:"id"`
	Name   string `json:"name,omitempty"`
	Email  string `json:"email,omitempty"`
	Phone  string `json:"phone,omitempty"`
	Fax    string `json:"fax,omitempty"`
	Postal string `json:"postal,omitempty"`
}
