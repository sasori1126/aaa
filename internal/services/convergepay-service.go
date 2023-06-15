package services

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"

	"github.com/bugsnag/bugsnag-go/v2"
)

// GetPaymentToken generates the required payment token by Converge.
// See: https://developer.elavon.com/na/docs/converge/1.0.0/integration-guide/integration_methods/checkoutjs
// Note: In order to call Converge pay API, it requires whitelist IP to get permission, even for local development.
// Need to call Converge pay support team to whitelist IP.
func GetPaymentToken(request *dto.GetConvergePayTokenRequest) (string, *entities.ApiError) {
	config, err := configs.GetConvergePayConfig()
	if err != nil {
		msg := "failed to get converge pay API config"
		return "", utils.FormatApiError(msg, configs.ServerError, entities.E{"convergePaymentToken": "server error"})
	}

	payload := newPaymentTokenPayload(*config, request)
	httpRequest, err := http.NewRequest("POST", config.APIBaseURL+"/transaction_token", bytes.NewBufferString(payload))
	if err != nil {
		msg := "failed to create converge getPaymentToken request"
		return "", utils.FormatApiError(msg, configs.ServerError, entities.E{"convergePaymentToken": "server error"})
	}
	httpRequest.Header.Set("content-type", "application/x-www-form-urlencoded")

	httpClient := &http.Client{Timeout: time.Second * 5}
	response, err := httpClient.Do(httpRequest)
	if err != nil {
		bugsnag.Notify(err)
		msg := "failed to make converge getPaymentToken request"
		return "", utils.FormatApiError(msg, configs.ServerError, entities.E{"convergePaymentToken": "server error"})
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		bugsnag.Notify(err)
		msg := "failed to read converge getPaymentToken response"
		return "", utils.FormatApiError(msg, configs.ServerError, entities.E{"convergePaymentToken": "server error"})
	}
	return string(body), nil
}

func newPaymentTokenPayload(config configs.ConvergePayConfig, request *dto.GetConvergePayTokenRequest) string {
	payload := []string{
		fmt.Sprintf("ssl_amount=%f", request.Amount),
		fmt.Sprintf("ssl_merchant_id=%s", config.MerchantID),
		fmt.Sprintf("ssl_pin=%s", config.MerchantPIN),
		"ssl_transaction_type=ccsale",
		fmt.Sprintf("ssl_user_id=%s", config.UserID),
	}
	return strings.Join(payload, "&")
}
