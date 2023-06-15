package payments

import (
	"axis/ecommerce-backend/configs"
	"encoding/json"
	"errors"
	"github.com/marvinhosea/bambora-go"
	"github.com/marvinhosea/bambora-go/client"
	"strconv"
	"strings"
)

type PaymentResponse struct {
	Code      int64    `json:"code"`
	Category  int64    `json:"category"`
	Message   string   `json:"message"`
	Reference string   `json:"reference"`
	Details   []Detail `json:"details"`
}

type Detail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Client struct {
	bambora  *client.Api
	currency *string
	config   *configs.BamboraConfig
}

func (c *Client) GetCurrency() (*string, error) {
	if c.currency == nil {
		return nil, errors.New("currency not set")
	}

	return c.currency, nil
}

func (c *Client) TakePayment(cardDetail *CardDetail, paymentAmount float32) (*PaymentIntegrationResponse, error) {
	err := c.ValidatePaymentDetail(cardDetail)
	if err != nil {
		return nil, err
	}

	formatExpiry := strings.Split(cardDetail.Expiry, "/")
	expiryMonth := formatExpiry[0]
	expiryYear := formatExpiry[1]

	c.bambora.Init(c.config.MerchantId, c.config.ProfilePasscode, c.config.PaymentPasscode)
	card, err := c.bambora.Card.Tokenize(&bambora.CardParams{
		Number:      String(cardDetail.CardNumber),
		ExpiryMonth: String(expiryMonth),
		ExpiryYear:  String(expiryYear),
		CVD:         String(cardDetail.CVV),
	})
	if err != nil {
		return nil, err
	}

	profile, err := c.bambora.Profile.New(&bambora.ProfileParams{
		CardName:  String(cardDetail.Holder),
		CardToken: String(card.Token),
	})
	if err != nil {
		return nil, err
	}

	payment, err := c.bambora.Payment.TakePayment(&bambora.PaymentParams{
		Amount:        paymentAmount,
		PaymentMethod: "payment_profile",
		Profile: bambora.PaymentProfile{
			CustomerCode: profile.CustomerCode,
			CardId:       "1",
			Complete:     true,
		},
	})
	if err != nil {
		return nil, err
	}

	pr := &PaymentResponse{}
	err = json.Unmarshal(payment.LastResponse.RawJson, pr)
	if err != nil {
		return nil, err
	}
	var isSuccessful = false
	if pr.Message == "Approved" {
		isSuccessful = true
	}

	return &PaymentIntegrationResponse{
		IsSuccessful: isSuccessful,
		Message:      pr.Message,
		Reference:    pr.Reference,
		Code:         strconv.FormatInt(pr.Code, 10),
		Amount:       paymentAmount,
	}, nil
}

func (c *Client) ValidatePaymentDetail(data *CardDetail) error {
	//getCard, ok := data.(CardDetail)
	//if !ok {
	//	return errors.New("invalid payment detail")
	//}

	if err := ValidateCard(data); err != nil {
		return err
	}
	return nil
}

func (c *Client) SetConfig() error {
	con, err := configs.GetBamboraConfig(*c.currency)
	if err != nil {
		return err
	}

	c.config = con
	return nil
}

func NewClient(cur string) (PaymentImplementation, error) {
	currency := cur
	c := &Client{
		currency: &currency,
	}

	err := c.SetConfig()
	if err != nil {
		return nil, err
	}

	c.bambora = &client.Api{}
	return c, nil
}
