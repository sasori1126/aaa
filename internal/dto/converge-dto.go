package dto

// GetConvergePayTokenRequest represents the request to get payment token from Converge pay API.
// See: https://developer.elavon.com/na/docs/converge/1.0.0/integration-guide/integration_methods/checkoutjs
type GetConvergePayTokenRequest struct {
	Amount float64
}
