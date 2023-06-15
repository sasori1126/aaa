package payments

import (
	"axis/ecommerce-backend/configs"
	"errors"
	"strconv"
	"strings"
)

func ValidateCard(card *CardDetail) error {
	//re := regexp.MustCompile(`^(?:4[0-9]{12}(?:[0-9]{3})?|[25][1-7][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\d{3})\d{11})$`)
	//if re.MatchString(card.CardNumber) {
	//	return errors.New("invalid card number" + card.CardNumber)
	//}

	formatExpiry := strings.Split(card.Expiry, "/")
	if len(formatExpiry) != 2 {
		return errors.New("invalid expiry date")
	}

	expiryMonth, err := strconv.Atoi(formatExpiry[0])
	if err != nil {
		return errors.New("invalid expiry date" + err.Error())
	}
	if expiryMonth < 1 && expiryMonth > 12 {
		return errors.New("invalid expiry date")
	}
	expiryYear := formatExpiry[1]
	if len(expiryYear) != 2 {
		return errors.New("invalid expiry date")
	}

	if len(card.CVV) == 0 {
		return errors.New("invalid credit card")
	}

	return nil
}

func String(v string) *string {
	return &v
}

func IsValidPaymentMethod(method string) bool {
	methods := map[string]bool{
		configs.PaymentMethodCard:       true,
		configs.PaymentMethodPhoneOrder: true,
		configs.PaymentMethodOnAccount:  true,
	}

	return methods[method]
}

func IsValidDeliveryMethod(method string) bool {
	methods := map[string]bool{
		configs.DeliveryMethodDeliver: true,
		configs.DeliveryMethodPickUp:  true,
	}

	return methods[method]
}
