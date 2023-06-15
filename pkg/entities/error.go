package entities

type E map[string]interface{}

type ApiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Errors  E      `json:"errors"`
}

//func (a *ApiError) Error() string {
//	return a.Message
//}
