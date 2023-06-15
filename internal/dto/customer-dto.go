package dto

type TempCustomerRequest struct {
	CustomerName    string `json:"customer_name"`
	CustomerEmail   string `json:"customer_email"`
	CustomerContact string `json:"customer_contact"`
	CustomerAddress string `json:"customer_address"`
	Items           []Item `json:"items"`
}

type Item struct {
	Title string `json:"title"`
}