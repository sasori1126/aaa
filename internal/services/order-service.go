package services

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal"
	"axis/ecommerce-backend/internal/domain/account"
	"axis/ecommerce-backend/internal/domain/carts"
	"axis/ecommerce-backend/internal/domain/orders"
	"axis/ecommerce-backend/internal/domain/payments"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/internal/notification/mail"
	"axis/ecommerce-backend/internal/notification/slack"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/trackerGenerator"
	"axis/ecommerce-backend/pkg/utils"

	"github.com/bugsnag/bugsnag-go/v2"
	"gorm.io/gorm"
)

type OrderService interface {
	CreateUserOrder(userHid string, request *dto.PlaceOrderRequest) *entities.ApiError
	AddOrderPayment(userID string, request *dto.AddOrderPayment) *entities.ApiError
	GetUserOrders(userHid string) ([]dto.OrderResponse, *entities.ApiError)
	GetOrders(limit, offset int, userHid string) ([]dto.OrderResponse, *entities.ApiError)
	GetOrderByID(orderHid string) (*dto.OrderResponse, *entities.ApiError)
	OrderUpdateStatus(orderHid string, req dto.OrderUpdateStatusRequest) *entities.ApiError
}

type DefaultOrderService struct {
	cartRepo    carts.CartRepo
	paymentRepo payments.PaymentRepo
	repo        orders.OrderRepo
	taxService  TaxService
	userRepo    account.UserRepo
}

func (d DefaultOrderService) AddOrderPayment(userHid string, request *dto.AddOrderPayment) *entities.ApiError {
	userId, err := models.DecodeHashId(userHid)
	if err != nil || userId == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"invalid user id",
			configs.BadRequest,
			entities.E{"order": "invalid user id"},
		)
	}
	getOrderId, err := models.DecodeHashId(request.OrderId)
	if err != nil || getOrderId == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error" + err.Error()},
		)
	}
	order, err := d.repo.GetOrderField(models.FindByField{
		Field: "id",
		Value: getOrderId,
	})
	if err != nil {
		bugsnag.Notify(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.FormatApiError(
				"order not found",
				configs.NotFound,
				entities.E{"order_id": "order not found"},
			)
		}

		return utils.FormatApiError(
			"failed to retrieve order",
			configs.ServerError,
			entities.E{"order_id": "failed to retrieve order"},
		)
	}

	if request.Amount < order.TotalAmount {
		return utils.FormatApiError(
			"payment amount cannot be less than order amount",
			configs.BadRequest,
			entities.E{"order": "payment amount less than order amount"},
		)
	}

	if order.Status != configs.OrderShipped && order.Status != configs.OrderPaid && order.Status != configs.OrderCancelled && order.Status != configs.OrderDelivered && order.Status != configs.OrderFailed {
		payment := models.Payment{
			PaymentMethod: request.PaymentMethod,
			PaymentAmount: request.Amount,
			Reference:     request.Reference,
			OrderId:       order.ID,
			Status:        configs.PaymentSuccessful,
			UserId:        userId,
		}

		err := d.paymentRepo.AddOrderPayment(&payment)
		if err != nil {
			bugsnag.Notify(err)
			return utils.FormatApiError(
				"failed to update cart",
				configs.ServerError,
				entities.E{"order": "server error failed " + err.Error() + ", try later"},
			)
		}

		err = d.repo.UpdateOrderField(order.ID, models.FindByField{
			Field: "status",
			Value: configs.OrderPaid,
		})
		if err != nil {
			return utils.FormatApiError(
				"failed to update order",
				configs.ServerError,
				entities.E{"order_id": "failed to update order"},
			)
		}
		data := make(map[string]string)
		data["NAME"] = order.User.Name
		data["USER_ORDERS_URL"] = "https://axisforestry.com/#/profile"
		data["PAYMENT_REFERENCE"] = payment.Reference

		emailTo := dto.MailData{
			Name:  order.User.Name,
			Email: order.User.Email,
		}
		replyTo := dto.MailData{
			Name:  "Axis",
			Email: configs.AxisPartSaleEmail,
		}

		err = mail.SendEmail(emailTo, replyTo, data, 32)
		if err != nil {
			slack.NewEvent(configs.GetSlackConfig(
				"failed to send email to user",
				"Error",
				"Create User Order: "+err.Error(),
				true,
			), nil, nil, nil)
		}
	} else {
		return utils.FormatApiError(
			"order status is "+order.Status,
			http.StatusBadRequest,
			entities.E{"order": "payment failed"},
		)
	}

	return nil
}

func (d DefaultOrderService) OrderUpdateStatus(orderHid string, req dto.OrderUpdateStatusRequest) *entities.ApiError {
	id, err := models.DecodeHashId(orderHid)
	if err != nil || id == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	err = d.repo.UpdateOrderField(id, models.FindByField{
		Field: "status",
		Value: req.Status,
	})
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to update order",
			configs.ServerError,
			entities.E{"order_id": "failed to update order"},
		)
	}

	return nil
}

func (d DefaultOrderService) GetOrderByID(orderHid string) (*dto.OrderResponse, *entities.ApiError) {
	id, err := models.DecodeHashId(orderHid)
	if err != nil || id == 0 {
		bugsnag.Notify(err)
		msg := "internal server error"
		return nil, utils.FormatApiError(msg, configs.ServerError, entities.E{"server": msg})
	}

	response, err := d.repo.GetOrderField(models.FindByField{Field: "id", Value: id})
	if err != nil {
		bugsnag.Notify(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg := "order not found"
			return nil, utils.FormatApiError(msg, configs.NotFound, entities.E{"order_id": msg})
		}

		msg := "failed to retrieve order"
		return nil, utils.FormatApiError(msg, configs.ServerError, entities.E{"order_id": msg})
	}

	order := response.ToResponse()
	return &order, nil
}

func (d DefaultOrderService) GetOrders(limit, offset int, userHid string) ([]dto.OrderResponse, *entities.ApiError) {
	userId, err := models.DecodeHashId(userHid)
	if err != nil && userHid != "" {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"invalid user id",
			configs.BadRequest,
			entities.E{"order": "invalid user id"},
		)
	}

	getOrders, err := d.repo.GetOrders(limit, offset, userId)
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"failed to retrieve orders",
			configs.ServerError,
			entities.E{"order": err.Error()},
		)
	}

	var orderRes []dto.OrderResponse
	orderRes = []dto.OrderResponse{}
	for _, order := range getOrders {
		or := order.ToResponse()
		orderRes = append(orderRes, or)
	}

	return orderRes, nil
}

func (d DefaultOrderService) GetUserOrders(userHid string) ([]dto.OrderResponse, *entities.ApiError) {
	userId, err := models.DecodeHashId(userHid)
	if err != nil || userId == 0 {
		bugsnag.Notify(err)
		msg := "invalid user id"
		return nil, utils.FormatApiError(msg, configs.BadRequest, entities.E{"order": msg})
	}

	getUserOrders, err := d.repo.GetUserOrders(userId)
	if err != nil {
		bugsnag.Notify(err)
		msg := "failed to getUserOrders"
		return nil, utils.FormatApiError(msg, configs.ServerError, entities.E{"order": msg})
	}

	orderRes := []dto.OrderResponse{}
	for _, order := range getUserOrders {
		orderRes = append(orderRes, order.ToResponse())
	}
	return orderRes, nil
}

func (d DefaultOrderService) CreateUserOrder(userHid string, request *dto.PlaceOrderRequest) *entities.ApiError {
	userId, err := models.DecodeHashId(userHid)
	if err != nil || userId == 0 {
		bugsnag.Notify(err)
		msg := "invalid user id"
		return utils.FormatApiError(msg, configs.BadRequest, entities.E{"order": msg})
	}

	user, err := d.userRepo.FindUserById(userId)
	if err != nil {
		bugsnag.Notify(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg := "user not found"
			return utils.FormatApiError(msg, configs.NotFound, entities.E{"order": msg})
		}
		return utils.FormatApiError(
			"Failed to process error ",
			configs.ServerError,
			entities.E{"order": "failed to retrieve user account"},
		)
	}

	cartId, err := models.DecodeHashId(request.CartId)
	if err != nil || cartId == 0 {
		bugsnag.Notify(err)
		msg := "invalid cart id"
		return utils.FormatApiError(msg, configs.BadRequest, entities.E{"order": msg})
	}

	getCart, err := d.cartRepo.GetCartIById(cartId)
	if err != nil {
		bugsnag.Notify(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg := "user's cart not found"
			return utils.FormatApiError(msg, configs.NotFound, entities.E{"order": msg})
		}
		msg := fmt.Sprintf("failed to get user cart by id: %d", cartId)
		return utils.FormatApiError(msg, configs.ServerError, entities.E{"order": "internal server error, try a later"})
	}

	if getCart.Status == internal.CartStatusOrdered {
		msg := "cart has already been submitted"
		return utils.FormatApiError(msg, configs.BadRequest, entities.E{"order": msg})
	}

	addressId, err := models.DecodeHashId(request.DeliveryAddressId)
	if err != nil || addressId == 0 {
		bugsnag.Notify(err)
		msg := "invalid address id"
		return utils.FormatApiError(msg, configs.BadRequest, entities.E{"order": msg})
	}

	userAddress, err := d.userRepo.GetUserAddressByID(addressId)
	if err != nil {
		bugsnag.Notify(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg := "user's address not found"
			return utils.FormatApiError(msg, configs.NotFound, entities.E{"order": msg})
		}
		return utils.FormatApiError(
			"failed to get user address",
			configs.ServerError,
			entities.E{"order": "internal server error, try a later"},
		)
	}

	if !payments.IsValidDeliveryMethod(request.DeliveryMethod) {
		msg := "invalid delivery method"
		return utils.FormatApiError(msg, configs.BadRequest, entities.E{"order": msg})
	}

	var shippingAmount float64
	deliveryMethod := configs.DeliveryMethodPickUp
	rateObjectId := request.DeliveryRateId
	if rateObjectId != "" && request.DeliveryMethod == configs.DeliveryMethodDeliver {
		shippoClient, err := internal.NewShippoClient()
		if err != nil {
			bugsnag.Notify(err)
			return utils.FormatApiError(
				"failed to retrieve rate",
				configs.ServerError,
				entities.E{"order": "internal server error, try a later"},
			)
		}

		getShippingRate, err := shippoClient.GetRate(rateObjectId)
		if err != nil {
			bugsnag.Notify(err)
			return utils.FormatApiError(
				"failed to retrieve rate",
				configs.ServerError,
				entities.E{"order": "internal server error, try a later"},
			)
		}
		deliveryMethod = fmt.Sprintf("%s %s", getShippingRate.Provider, getShippingRate.ServiceLevel.Name)

		shippingAmount, err = strconv.ParseFloat(getShippingRate.Amount, 32)
		if err != nil {
			bugsnag.Notify(err)
			return utils.FormatApiError(
				"failed to get shipping rates",
				configs.ServerError,
				entities.E{"order": "internal server error, try a later"},
			)
		}
		// Add extra 20% based on the Shippo returned shipping amount. This value need to be adjustable in the admin portal.
		// Please sync with function: CreateShipment in shippo.go.
		shippingAmount = utils.RoundOff(shippingAmount * 1.2)
	}

	var productAmount float64
	var orderItems []models.OrderItem
	for _, item := range getCart.Items {
		o := models.OrderItem{
			Description:  item.Description,
			Name:         item.Title,
			Note:         item.Note,
			PartId:       item.PartId,
			PricePerUnit: item.PricePerUnit,
			Quantity:     item.Quantity,
			Total:        item.Amount,
			Unit:         item.Unit,
			Weight:       item.Weight,
		}
		orderItems = append(orderItems, o)
		productAmount += item.Amount
	}
	if len(orderItems) == 0 {
		msg := "cart cannot be empty"
		return utils.FormatApiError(msg, configs.BadRequest, entities.E{"order": msg})
	}
	productAmount = utils.RoundOff(productAmount)

	taxLines, _, err := d.taxService.GetTaxesByAddress(userId, &userAddress.Address)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to retrieve cart taxes",
			configs.ServerError,
			entities.E{"order": "internal server error, try a later"},
		)
	}

	var taxTotal float64
	var productTax float64
	var shippingTax float64
	var orderTaxes []models.OrderTax
	for _, tax := range taxLines {
		currentProductTax := utils.RoundOff(productAmount * tax.Rate)
		currentShippingTax := utils.RoundOff(shippingAmount * tax.Rate)
		currentTaxAmount := currentProductTax + currentShippingTax

		productTax += currentProductTax
		shippingTax += currentProductTax
		taxTotal += currentTaxAmount

		orderTaxes = append(orderTaxes, models.OrderTax{
			Amount:            currentTaxAmount,
			Description:       tax.Description,
			Name:              tax.Name,
			ProductTaxAmount:  currentProductTax,
			Rate:              tax.Rate,
			ShippingTaxAmount: currentShippingTax,
		})
	}

	generator, err := trackerGenerator.New(trackerGenerator.Config{
		Charset: "alphanumericUpperCase",
		Pattern: "####-####-####",
	})
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to create tracking number",
			configs.ServerError,
			entities.E{"order": "internal server error, try a later"},
		)
	}

	trackingCode, err := generator.GenerateOne()
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to create tracking number",
			configs.ServerError,
			entities.E{"order": "internal server error, try a later"},
		)
	}

	if !payments.IsValidPaymentMethod(request.PaymentMethod) {
		return utils.FormatApiError(
			"invalid payment method",
			configs.BadRequest,
			entities.E{"order": "invalid payment method"},
		)
	}

	totalAmount := productAmount + shippingAmount + taxTotal
	var totalPaidAmount float64
	orderStatus := configs.OrderPending
	currency := user.DefaultCurrency
	var paymentReference string
	switch request.PaymentMethod {
	case configs.PaymentMethodCard:
		totalPaidAmount = totalAmount
		if isConvergePay(currency) {
			orderStatus = configs.OrderPaid
			paymentReference = request.ConvergePayTransactionId
		}

		if !isConvergePay(currency) {
			card := request.Card
			if card.Holder == "" || card.CVV == "" || card.Expiry == "" || card.CardNumber == "" {
				msg := "invalid card information"
				return utils.FormatApiError(msg, configs.BadRequest, entities.E{"order": msg})
			}

			payment, err := payments.NewPayment(currency)
			if err != nil {
				bugsnag.Notify(err)
				return utils.FormatApiError(
					"failed to initiate payment",
					configs.ServerError,
					entities.E{"order": "internal server error, try a later"},
				)
			}
			cr := payments.CardDetail{
				Holder:     request.Card.Holder,
				CardNumber: strings.TrimSpace(request.Card.CardNumber),
				Expiry:     strings.TrimSpace(request.Card.Expiry),
				CVV:        strings.TrimSpace(request.Card.CVV),
			}
			err = payment.Bambora.ValidatePaymentDetail(&cr)
			if err != nil {
				bugsnag.Notify(err)
				return utils.FormatApiError(
					"failed to validate payment",
					configs.ServerError,
					entities.E{"order": "internal server error, try a later" + err.Error()},
				)
			}

			paymentRes, err := payment.Bambora.TakePayment(&cr, float32(totalAmount))
			if err != nil {
				bugsnag.Notify(err)
				return utils.FormatApiError(
					"failed to make payment",
					configs.ServerError,
					entities.E{"order": "internal server error, try a later"},
				)
			}
			if !paymentRes.IsSuccessful {
				return utils.FormatApiError("failed to make payment", configs.BadRequest, entities.E{"order": paymentRes})
			}

			paymentReference = paymentRes.Reference
			orderStatus = configs.OrderPaid
		}
	case configs.PaymentMethodPhoneOrder:
		orderStatus = configs.OrderPendingOnPhone
		totalPaidAmount = 0
	case configs.PaymentMethodOnAccount:
		orderStatus = configs.OrderPendingOnAccount
		totalPaidAmount = 0
	default:
		return utils.FormatApiError(
			"invalid payment method",
			configs.BadRequest,
			entities.E{"order": "invalid payment method"},
		)
	}

	order := models.Order{
		UserId:                   userId,
		CartId:                   getCart.ID,
		PaidAmount:               totalPaidAmount,
		TrackingNumber:           trackingCode,
		Currency:                 currency,
		ShippingAddressId:        userAddress.ID,
		BillingAddressId:         userAddress.ID,
		ShippingAmount:           shippingAmount,
		PurchaseOrderNumber:      request.PurchaseOrderNumber,
		TotalTaxAmount:           taxTotal,
		TotalAmount:              totalAmount,
		SubTotalAmount:           productAmount,
		Items:                    orderItems,
		Taxes:                    orderTaxes,
		Status:                   orderStatus,
		PaymentMethod:            request.PaymentMethod,
		ConvergePayTransactionId: request.ConvergePayTransactionId,
		DeliveryRateId:           request.DeliveryRateId,
		DeliveryInstruction:      request.DeliveryInstruction,
		DeliveryMethod:           deliveryMethod,
		PickupPoint:              request.PickupPoint,
	}
	ordered, err := d.repo.CreateOrder(order)
	if err != nil {
		bugsnag.Notify(err)
		if errors.Is(err, configs.ErrOrderAlreadyExist) {
			getCart.Status = internal.CartStatusOrdered
			err := d.cartRepo.UpdateCart(getCart)
			if err != nil {
				bugsnag.Notify(err)
				return utils.FormatApiError(
					"failed to update cart",
					configs.ServerError,
					entities.E{"order": "server error failed " + err.Error() + ", try later"},
				)
			}
		}

		return utils.FormatApiError(
			"failed to place order, try again later",
			configs.BadRequest,
			entities.E{"order": "failed to place order, " + err.Error() + ", try later"},
		)
	}

	if ordered.Status == configs.OrderPaid {
		payment := models.Payment{
			PaymentMethod: request.PaymentMethod,
			PaymentAmount: ordered.TotalAmount,
			Reference:     paymentReference,
			OrderId:       ordered.ID,
			Status:        configs.PaymentSuccessful,
			UserId:        userId,
		}

		err := d.paymentRepo.AddOrderPayment(&payment)
		if err != nil {
			bugsnag.Notify(err)
			return utils.FormatApiError(
				"failed to update cart",
				configs.ServerError,
				entities.E{"order": "server error failed " + err.Error() + ", try later"},
			)
		}

		data := make(map[string]string)
		data["NAME"] = user.Name
		data["USER_ORDERS_URL"] = "https://axisforestry.com/#/profile"
		data["PAYMENT_REFERENCE"] = payment.Reference

		emailTo := dto.MailData{
			Name:  user.Name,
			Email: user.Email,
		}
		replyTo := dto.MailData{
			Name:  "Axis",
			Email: configs.AxisPartSaleEmail,
		}

		_ = mail.SendEmail(emailTo, replyTo, data, 32)
	}

	getCart.Status = internal.CartStatusOrdered
	err = d.cartRepo.UpdateCart(getCart)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to update cart",
			configs.ServerError,
			entities.E{"order": "server error failed " + err.Error() + ", try later"},
		)
	}

	orderedId, _ := models.EncodeHashId(ordered.ID)
	// Email Axis
	data := make(map[string]string)
	data["NAME"] = user.Name
	data["EMAIL"] = user.Email
	data["ORDER_URL"] = "https://the-a-team.axisforestry.com/#/orders/" + orderedId

	emailTo := dto.MailData{
		Name:  "Axis",
		Email: configs.AxisPartSaleEmail,
	}
	replyTo := dto.MailData{
		Name:  user.Name,
		Email: user.Email,
	}
	_ = mail.SendEmail(emailTo, replyTo, data, 25)

	// Email user
	userData := make(map[string]string)
	userData["NAME"] = user.Name
	userData["USER_ORDERS_URL"] = "https://axisforestry.com/#/profile"
	userData["ORDER_STATUS"] = order.Status

	emailToUser := dto.MailData{
		Name:  user.Name,
		Email: user.Email,
	}
	replyToUser := dto.MailData{
		Name:  "Axis",
		Email: configs.AxisPartSaleEmail,
	}
	_ = mail.SendEmail(emailToUser, replyToUser, data, 31)

	color := "#00FF00"
	link := "https://the-a-team.axisforestry.com/#/orders/" + orderedId
	slack.NewEvent(configs.GetSlackConfig(
		"New Axis Order",
		"Success",
		"Order was successfully placed. View: "+link,
		false,
	), nil, &link, &color)

	return nil
}

func NewDefaultOrderService(
	repo orders.OrderRepo,
	cartRepo carts.CartRepo,
	userRepo account.UserRepo,
	taxService TaxService,
	paymentRepo payments.PaymentRepo,
) OrderService {
	return &DefaultOrderService{
		cartRepo:    cartRepo,
		paymentRepo: paymentRepo,
		repo:        repo,
		taxService:  taxService,
		userRepo:    userRepo,
	}
}

func isConvergePay(currency string) bool {
	// Use Converge pay for USD and Bambora pay for CAD, and treat consumers with phone country code +1 only uses CAD.
	return currency == "USD"
}
