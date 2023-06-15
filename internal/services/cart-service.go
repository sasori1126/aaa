package services

import (
	"errors"
	"fmt"
	"sort"
	"strconv"

	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal"
	"axis/ecommerce-backend/internal/domain/account"
	"axis/ecommerce-backend/internal/domain/carts"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/internal/notification/slack"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"

	"github.com/bugsnag/bugsnag-go/v2"
	shippoModels "github.com/coldbrewcloud/go-shippo/models"
	"gorm.io/gorm"
)

type CartService interface {
	AddCartItem(userId string, request *dto.CartItemRequest, part dto.PartResponse) (*dto.CartItemResponse, *entities.ApiError)
	DelCartItem(cartItemsId []string, userHid string) *entities.ApiError
	GetCartByID(cartHid string) (*dto.CartResponse, *entities.ApiError)
	GetCarts(limit, offset int, userHid, status string) ([]dto.CartResponse, *entities.ApiError)
	GetDeliveryRates(
		deliveryAddressId,
		userHid string,
	) ([]dto.DeliveryRateResponse, []dto.TaxLine, []dto.TaxLine, *entities.ApiError)
	GetUserActiveCart(userId string) (*dto.CartResponse, *entities.ApiError)
	UpdateCartItemQuantity(cartItemId, userHid string, quantity int, isAddition bool) *entities.ApiError
}

type DefaultCartService struct {
	repo       carts.CartRepo
	taxService TaxService
	userRepo   account.UserRepo
}

func (d DefaultCartService) GetCartByID(cartHid string) (*dto.CartResponse, *entities.ApiError) {
	id, err := models.DecodeHashId(cartHid)
	if err != nil || id == 0 {
		errMessage := "internal server error"
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(errMessage, configs.ServerError, entities.E{"order": errMessage})
	}

	getCart, err := d.repo.GetCartIById(id)
	if err != nil {
		bugsnag.Notify(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.FormatApiError(
				"user's cart not found",
				configs.NotFound,
				entities.E{"order": "user's cart not found"},
			)
		}
		return nil, utils.FormatApiError(
			"failed to get user cart",
			configs.ServerError,
			entities.E{"order": "internal server error, try a later"},
		)
	}

	cart := getCart.ToResponse()
	return &cart, nil
}

func (d DefaultCartService) GetCarts(limit, offset int, userHid, status string) ([]dto.CartResponse, *entities.ApiError) {
	userId, err := models.DecodeHashId(userHid)
	if err != nil && userHid != "" {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"invalid user id",
			configs.BadRequest,
			entities.E{"order": "invalid user id"},
		)
	}

	getCarts, err := d.repo.GetCarts(limit, offset, userId, status)
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"failed to retrieve carts",
			configs.ServerError,
			entities.E{"order": err.Error()},
		)
	}

	var orderRes []dto.CartResponse
	orderRes = []dto.CartResponse{}
	for _, cart := range getCarts {
		cid, _ := models.EncodeHashId(cart.ID)
		uid, _ := models.EncodeHashId(cart.UserId)
		var items []dto.CartItemResponse
		for _, item := range cart.Items {
			iid, _ := models.EncodeHashId(item.ID)
			p := item.Part.ToResponse()
			i := dto.CartItemResponse{
				Id:           iid,
				CartId:       cid,
				PartId:       p.Id,
				PricePerUnit: item.PricePerUnit,
				Quantity:     item.Quantity,
				Unit:         item.Unit,
				Description:  item.Description,
				Note:         item.Note,
				Title:        item.Title,
				Amount:       item.Amount,
				Part:         p,
			}

			items = append(items, i)
		}
		cr := dto.CartResponse{
			Id: cid,
			User: dto.EmbedUser{
				Id:          uid,
				Name:        cart.User.Name,
				Email:       cart.User.Email,
				PhoneNumber: cart.User.PhoneNumber,
				IsActive:    cart.User.IsActive,
				Verified:    cart.User.IsActive,
			},
			Items:       items,
			SubTotal:    cart.SubTotal,
			TotalAmount: cart.TotalAmount,
			Status:      cart.Status,
		}
		orderRes = append(orderRes, cr)
	}

	return orderRes, nil
}

func (d DefaultCartService) GetDeliveryRates(
	deliveryAddressId, userHid string,
) ([]dto.DeliveryRateResponse, []dto.TaxLine, []dto.TaxLine, *entities.ApiError) {
	id, err := models.DecodeHashId(deliveryAddressId)
	if err != nil || id == 0 {
		bugsnag.Notify(err)
		errMsg := "Invalid delivery address id"
		go slack.NewEvent(configs.GetSlackConfig(errMsg, "Error", "Create User Order: "+err.Error(), true), nil, nil, nil)
		return nil, nil, nil, utils.FormatApiError(errMsg, configs.BadRequest, entities.E{"deliveryAddressId": errMsg})
	}

	userId, err := models.DecodeHashId(userHid)
	if err != nil || userId == 0 {
		bugsnag.Notify(err)
		errMsg := "Invalid user id"
		go slack.NewEvent(configs.GetSlackConfig(errMsg, "Error", "Create User Order: "+err.Error(), true), nil, nil, nil)
		return nil, nil, nil, utils.FormatApiError(errMsg, configs.BadRequest, entities.E{"user id": errMsg})
	}

	getCart, err := d.repo.GetUserActiveCart(userId)
	if err != nil {
		bugsnag.Notify(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errMsg := "User's cart is empty"
			go slack.NewEvent(configs.GetSlackConfig(errMsg, "Error", "Create User Order: "+err.Error(), true), nil, nil, nil)
			return nil, nil, nil, utils.FormatApiError(errMsg, configs.NotFound, entities.E{"cart": errMsg})
		}

		errMsg := "Internal server error, failed to get user cart"
		return nil, nil, nil, utils.FormatApiError(errMsg, configs.ServerError, entities.E{"server": errMsg})
	}

	userAddress, err := d.userRepo.GetUserAddressByField(models.FindByField{Field: "id", Value: id})
	if err != nil {
		bugsnag.Notify(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errMsg := "Address not found"
			return nil, nil, nil, utils.FormatApiError(errMsg, configs.NotFound, entities.E{"address_id": errMsg})
		}
		errMsg := "Internal server error, failed to get user address by id"
		return nil, nil, nil, utils.FormatApiError(errMsg, configs.ServerError, entities.E{"address_id": errMsg})
	}

	shippoClient, err := internal.NewShippoClient()
	if err != nil {
		bugsnag.Notify(err)
		errMsg := "Failed to get shipping rates"
		return nil, nil, nil, utils.FormatApiError(errMsg, configs.ServerError, entities.E{"error": err.Error()})
	}

	userDefaultCurrency := getCart.User.DefaultCurrency
	getAxisAddress := configs.GetAxisDeliveryAddress(userAddress.Address.Country, &userDefaultCurrency)
	from := shippoModels.AddressInput{
		Name:     getAxisAddress.Name,
		Street1:  getAxisAddress.Street1,
		Street2:  getAxisAddress.Street2,
		City:     getAxisAddress.City,
		State:    getAxisAddress.State,
		Zip:      getAxisAddress.Zip,
		Country:  getAxisAddress.Country,
		Phone:    getAxisAddress.Phone,
		Email:    getAxisAddress.Email,
		Company:  getAxisAddress.Company,
		StreetNo: getAxisAddress.StreetNo,
	}

	to := shippoModels.AddressInput{
		Name:     userAddress.Contact.Name,
		Street1:  userAddress.Address.StreetName,
		City:     userAddress.Address.City,
		State:    userAddress.Address.State,
		Zip:      userAddress.Address.ZipCode,
		Country:  userAddress.Address.Country,
		Phone:    userAddress.Contact.Phone,
		Email:    userAddress.Contact.Email,
		Company:  userAddress.Contact.Organisation,
		StreetNo: "",
	}

	cartWeight := getWeight(*getCart)
	parcel := dto.ShippoParcel{
		Length:       "10",
		Width:        "10",
		Height:       "10",
		DistanceUnit: shippoModels.DistanceUnitInch,
		Weight:       fmt.Sprintf("%.4f", 300.00),
		MassUnit:     shippoModels.MassUnitKiloGram,
	}

	shipment, err := shippoClient.CreateShipment(from, to, cartWeight, parcel)
	if err != nil {
		bugsnag.Notify(err)
		errMsg := "Failed to get shipping rates"
		return nil, nil, nil, utils.FormatApiError(errMsg, configs.ServerError, entities.E{"error": err.Error()})
	}

	taxLines, taxExemptionLines, err := d.taxService.GetTaxesByAddress(userId, &userAddress.Address)
	if err != nil {
		bugsnag.Notify(err)
		msg := "failed to retrieve taxes"
		return nil, nil, nil, utils.FormatApiError(msg, configs.ServerError, entities.E{"order": msg})
	}

	deliveryRates := []dto.DeliveryRateResponse{}
	for _, rate := range shipment.Rates {
		deliveryAmount, _ := strconv.ParseFloat(rate.Amount, 64)
		r := dto.DeliveryRateResponse{
			Amount:           deliveryAmount,
			Attributes:       rate.Attributes,
			CarrierAccount:   rate.CarrierAccount,
			Currency:         rate.Currency,
			CurrencyLocal:    rate.CurrencyLocal,
			Days:             rate.Days,
			DurationTerms:    rate.DurationTerms,
			Expected:         rate.ServiceLevel.Name,
			Id:               rate.ObjectID,
			Provider:         rate.Provider,
			ProviderImage200: rate.ProviderImage200,
			ProviderImage75:  rate.ProviderImage75,
			Zone:             rate.Zone,
		}
		if r.Id != "" {
			deliveryRates = append(deliveryRates, r)
		}
	}
	sort.Slice(deliveryRates, func(i, j int) bool { return deliveryRates[i].Amount > deliveryRates[j].Amount })

	taxes := []dto.TaxLine{}
	for _, tax := range taxLines {
		deliveryTaxAmounts := []dto.DeliveryTaxAmount{}
		for _, deliveryRate := range deliveryRates {
			deliveryTaxAmounts = append(deliveryTaxAmounts, dto.DeliveryTaxAmount{
				DeliveryRateId:    deliveryRate.Id,
				DeliveryTaxAmount: utils.RoundOff(deliveryRate.Amount * tax.Rate),
			})
		}
		productTaxAmount := getCart.TotalAmount * tax.Rate
		taxes = append(taxes, dto.TaxLine{
			DeliveryTaxAmounts: deliveryTaxAmounts,
			Country:            tax.Country,
			Name:               tax.Name,
			ProductTaxAmount:   utils.RoundOff(productTaxAmount),
			Rate:               tax.Rate,
			State:              tax.State,
		})
	}

	taxExemptions := []dto.TaxLine{}
	for _, tax := range taxExemptionLines {
		taxExemptions = append(taxExemptions, dto.TaxLine{
			DeliveryTaxAmounts: make([]dto.DeliveryTaxAmount, 0),
			Country:            tax.Country,
			Name:               tax.Name,
			ProductTaxAmount:   0,
			Rate:               tax.Rate,
			State:              tax.State,
		})
	}
	return deliveryRates, taxes, taxExemptions, nil
}

func getWeight(cart models.Cart) float32 {
	var totalWeight float64 = 0
	for _, item := range cart.Items {
		w := item.Weight
		totalWeight += w
	}
	if totalWeight <= 0 {
		totalWeight = 5
	}
	return float32(totalWeight)
}

func (d DefaultCartService) UpdateCartItemQuantity(cartItemId, userHid string, quantity int, isAddition bool) *entities.ApiError {
	id, err := models.DecodeHashId(cartItemId)
	if err != nil || id == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"invalid id",
			configs.BadRequest,
			entities.E{"id": "invalid id"},
		)
	}

	item, err := d.repo.GetCartItem(id)
	if err != nil {
		bugsnag.Notify(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.FormatApiError(
				"Item not found",
				configs.NotFound,
				entities.E{"cart_item_id": "Item not found"},
			)
		}
		return utils.FormatApiError(
			"failed to get item, try later",
			configs.ServerError,
			entities.E{"id": "internal server, try later"},
		)
	}

	userId, err := models.DecodeHashId(userHid)
	if err != nil || id == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"invalid user id",
			configs.BadRequest,
			entities.E{"user id": "invalid user id"},
		)
	}

	getCart, err := d.repo.GetUserActiveCart(userId)
	if err != nil {
		bugsnag.Notify(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.FormatApiError(
				"user's cart is empty",
				configs.NotFound,
				entities.E{"cart": "user's cart is empty"},
			)
		}

		return utils.FormatApiError(
			"failed to get user cart",
			configs.ServerError,
			entities.E{"cart": "internal server error, try alater"},
		)
	}

	getCart.SubTotal -= item.Amount
	getCart.TotalAmount -= item.Amount

	var newQuantity float64
	if isAddition {
		newQuantity = item.Quantity + float64(quantity)
	} else {
		newQuantity = item.Quantity - float64(quantity)
	}

	if newQuantity <= 0 {
		return utils.FormatApiError(
			"quantity cannot be less than 1",
			configs.ServerError,
			entities.E{"quantity": "quantity cannot be less than 1"},
		)
	}

	item.Amount = utils.RoundOff(newQuantity * item.PricePerUnit)
	item.Quantity = newQuantity

	err = d.repo.UpdateCartItemQuantity(item)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to get item, try later",
			configs.ServerError,
			entities.E{"id": "internal server, try later"},
		)
	}

	getCart.SubTotal += item.Amount
	getCart.TotalAmount += item.Amount
	err = d.repo.UpdateCart(getCart)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to update cart, try later",
			configs.ServerError,
			entities.E{"id": "internal server, try later"},
		)
	}

	return nil
}

func (d DefaultCartService) DelCartItem(cartItemsId []string, userHid string) *entities.ApiError {
	itemsId, err := models.DecodeMultipleHash(cartItemsId...)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"invalid id",
			configs.BadRequest,
			entities.E{"id": "invalid id"},
		)
	}

	userId, err := models.DecodeHashId(userHid)
	if err != nil || userId == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"invalid user id",
			configs.BadRequest,
			entities.E{"user id": "invalid user id"},
		)
	}

	getCart, err := d.repo.GetUserActiveCart(userId)
	if err != nil {
		bugsnag.Notify(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.FormatApiError(
				"user's cart is empty",
				configs.NotFound,
				entities.E{"cart": "user's cart is empty"},
			)
		}

		return utils.FormatApiError(
			"failed to get user cart",
			configs.ServerError,
			entities.E{"cart": "internal server error, try alater"},
		)
	}

	items, err := d.repo.GetCartItemsByQuery(models.QueryByField{Query: "id IN (?)", Value: itemsId})
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to remove items, try later",
			configs.ServerError,
			entities.E{"id": "internal server, try later"},
		)
	}

	var amountRemoved float64 = 0
	for _, item := range items {
		amountRemoved += item.Amount
	}

	err = d.repo.DelCartItem(itemsId)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to remove items, try later",
			configs.ServerError,
			entities.E{"id": "internal server, try later"},
		)
	}

	getCart.SubTotal -= amountRemoved
	getCart.TotalAmount -= amountRemoved
	err = d.repo.UpdateCart(getCart)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError(
			"failed to update cart, try later",
			configs.ServerError,
			entities.E{"id": "internal server, try later"},
		)
	}

	return nil
}

func (d DefaultCartService) AddCartItem(
	userHid string,
	request *dto.CartItemRequest,
	part dto.PartResponse,
) (*dto.CartItemResponse, *entities.ApiError) {
	userId, err := models.DecodeHashId(userHid)
	if err != nil || userId == 0 {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	getActiveCart, err := d.repo.GetUserActiveCart(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			getActiveCart, err = d.repo.CreateCart(models.Cart{
				UserId:      userId,
				Status:      internal.CartStatusActive,
				SubTotal:    0,
				TotalAmount: 0,
			})

			if err != nil {
				bugsnag.Notify(err)
				return nil, utils.FormatApiError(
					"failed to create cart",
					configs.ServerError,
					entities.E{"server": "internal server error"},
				)
			}
		} else {
			bugsnag.Notify(err)
			return nil, utils.FormatApiError(
				"failed to create cart",
				configs.ServerError,
				entities.E{"server": "internal server error"},
			)
		}
	}

	pid, err := models.DecodeHashId(request.PartId)
	if err != nil || pid == 0 {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	amount := utils.RoundOff(request.Quantity * part.Price)
	weight := part.Weight * request.Quantity
	item := models.CartItem{
		PartId:       pid,
		PricePerUnit: part.Price,
		Quantity:     request.Quantity,
		Unit:         "piece",
		Description:  part.Description,
		Note:         request.Note,
		Title:        part.Name,
		Weight:       weight,
		Amount:       amount,
		Cart:         *getActiveCart,
	}

	getItem, err := d.repo.AddCartItem(item)
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	getActiveCart.SubTotal += amount
	getActiveCart.TotalAmount += amount

	err = d.repo.UpdateCart(getActiveCart)
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"internal server error",
			configs.ServerError,
			entities.E{"server": "internal server error"},
		)
	}

	itemResponse := getItem.ToResponse()
	return &itemResponse, nil
}

func (d DefaultCartService) GetUserActiveCart(userId string) (*dto.CartResponse, *entities.ApiError) {
	id, err := models.DecodeHashId(userId)
	if err != nil || id == 0 {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError(
			"invalid id",
			configs.BadRequest,
			entities.E{"id": "invalid id"},
		)
	}

	cart, err := d.repo.GetUserActiveCart(id)
	if err != nil {
		bugsnag.Notify(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart, err = d.repo.CreateCart(models.Cart{
				UserId:      id,
				Status:      internal.CartStatusActive,
				SubTotal:    0,
				TotalAmount: 0,
			})
			if err != nil {
				bugsnag.Notify(err)
				return nil, utils.FormatApiError(
					"failed to get cart, try later",
					configs.ServerError,
					entities.E{"id": "internal server, try later", "error": err.Error()},
				)
			}
		} else {
			return nil, utils.FormatApiError(
				"failed to get cart, try later",
				configs.ServerError,
				entities.E{"id": "internal server, try later", "error": err.Error()},
			)
		}
	}

	response := cart.ToResponse()
	response.SubTotal = utils.RoundOff(response.SubTotal)
	response.TotalAmount = utils.RoundOff(response.TotalAmount)
	return &response, nil
}

func NewDefaultCartService(repo carts.CartRepo, userRepo account.UserRepo, taxService TaxService) CartService {
	return &DefaultCartService{repo: repo, userRepo: userRepo, taxService: taxService}
}
