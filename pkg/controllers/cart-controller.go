package controllers

import (
	"net/http"

	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/account"
	"axis/ecommerce-backend/internal/domain/carts"
	"axis/ecommerce-backend/internal/domain/parts"
	"axis/ecommerce-backend/internal/domain/taxes"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/services"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (s *Serve) GetDeliveryRates(c *gin.Context) {
	addressId := c.Param("addressId")
	if addressId == "" {
		msg := "bad request"
		c.JSON(configs.BadRequest, utils.FormatApiError(msg, configs.BadRequest, entities.E{"head_id": msg}))
		return
	}

	user, ok := c.Get("user")
	if !ok {
		msg := "user account not found"
		c.JSON(configs.Unauthorized, utils.FormatApiError(msg, configs.Unauthorized, entities.E{"account": msg}))
		return
	}

	service := services.NewDefaultCartService(
		carts.NewCartRepoDb(s.Db),
		account.NewUserRepoDb(s.Db),
		services.NewDefaultTaxService(taxes.NewTaxRepoDb(s.Db)),
	)
	rates, taxes, taxExemptions, apiErr := service.GetDeliveryRates(addressId, user.(dto.UserSession).ID)
	if apiErr != nil {
		c.JSON(apiErr.Code, utils.FormatApiError(apiErr.Message, apiErr.Code, apiErr.Errors))
		return
	}

	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "delivery rates",
		Data:    entities.E{"delivery_rates": rates, "taxes": taxes, "tax_exemptions": taxExemptions},
		Links:   entities.E{},
		Meta:    entities.E{},
	})
}

func (s *Serve) GetUserActiveCart(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		msg := "user account not found"
		c.JSON(configs.Unauthorized, utils.FormatApiError(msg, configs.Unauthorized, entities.E{"account": msg}))
		return
	}

	service := services.NewDefaultCartService(
		carts.NewCartRepoDb(s.Db),
		account.NewUserRepoDb(s.Db),
		services.NewDefaultTaxService(taxes.NewTaxRepoDb(s.Db)),
	)
	cart, apiErr := service.GetUserActiveCart(user.(dto.UserSession).ID)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "user cart",
		Data:    cart,
		Links:   entities.E{},
		Meta:    entities.E{},
	})
}

func (s *Serve) RemoveCartItem(c *gin.Context) {
	request := &dto.CartItemRemoveRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		vErrs, ok := err.(validator.ValidationErrors)
		msg := "bad request, failed to process request"
		if ok {
			fieldErrors := utils.FormatValidationError(vErrs)
			c.JSON(configs.BadRequest, utils.FormatApiError(msg, configs.BadRequest, fieldErrors))
			return
		}
		c.JSON(configs.BadRequest, utils.FormatApiError(msg, configs.BadRequest, entities.E{"message": err.Error()}))
		return
	}

	user, ok := c.Get("user")
	if !ok {
		msg := "user account not found"
		c.JSON(configs.Unauthorized, utils.FormatApiError(msg, configs.Unauthorized, entities.E{"account": msg}))
		return
	}

	cs := services.NewDefaultCartService(
		carts.NewCartRepoDb(s.Db),
		account.NewUserRepoDb(s.Db),
		services.NewDefaultTaxService(taxes.NewTaxRepoDb(s.Db)),
	)
	apiErr := cs.DelCartItem(request.CartItemsId, user.(dto.UserSession).ID)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{}, "Item removed successfully", utils.Empty(), utils.Empty()))
}

func (s *Serve) IncreaseCartItem(c *gin.Context) {
	request := &dto.UpdateCartItemQuantityRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		vErrs, ok := err.(validator.ValidationErrors)
		msg := "bad request, failed to process request"
		if ok {
			fieldErrors := utils.FormatValidationError(vErrs)
			c.JSON(configs.BadRequest, utils.FormatApiError(msg, configs.BadRequest, fieldErrors))
			return
		}
		c.JSON(configs.BadRequest, utils.FormatApiError(msg, configs.BadRequest, entities.E{"message": err.Error()}))
		return
	}

	user, ok := c.Get("user")
	if !ok {
		msg := "user account not found"
		c.JSON(configs.Unauthorized, utils.FormatApiError(msg, configs.Unauthorized, entities.E{"account": msg}))
		return
	}

	service := services.NewDefaultCartService(
		carts.NewCartRepoDb(s.Db),
		account.NewUserRepoDb(s.Db),
		services.NewDefaultTaxService(taxes.NewTaxRepoDb(s.Db)),
	)
	apiErr := service.UpdateCartItemQuantity(request.CartItemId, user.(dto.UserSession).ID, request.Quantity, true)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{}, "Updated cart item successfully", utils.Empty(), utils.Empty()))
}

func (s *Serve) DecreaseCartItem(c *gin.Context) {
	request := &dto.UpdateCartItemQuantityRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		vErrs, ok := err.(validator.ValidationErrors)
		msg := "bad request, failed to process request"
		if ok {
			fieldErrors := utils.FormatValidationError(vErrs)
			c.JSON(configs.BadRequest, utils.FormatApiError(msg, configs.BadRequest, fieldErrors))
			return
		}
		c.JSON(configs.BadRequest, utils.FormatApiError(msg, configs.BadRequest, entities.E{"message": err.Error()}))
		return
	}

	user, ok := c.Get("user")
	if !ok {
		msg := "user account not found"
		c.JSON(configs.Unauthorized, utils.FormatApiError(msg, configs.Unauthorized, entities.E{"account": msg}))
		return
	}

	cs := services.NewDefaultCartService(
		carts.NewCartRepoDb(s.Db),
		account.NewUserRepoDb(s.Db),
		services.NewDefaultTaxService(taxes.NewTaxRepoDb(s.Db)),
	)
	apiErr := cs.UpdateCartItemQuantity(request.CartItemId, user.(dto.UserSession).ID, request.Quantity, false)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{}, "Updated cart item successfully", utils.Empty(), utils.Empty()))
}

func (s *Serve) AddCartItem(c *gin.Context) {
	ctx, span := otel.Tracer("").Start(c.Request.Context(), "userAddItem")
	defer span.End()

	cr := &dto.CartItemRequest{}
	if err := c.ShouldBindJSON(cr); err != nil {
		vErrs, ok := err.(validator.ValidationErrors)
		msg := "bad request, failed to process request"
		if ok {
			fieldErrors := utils.FormatValidationError(vErrs)
			c.JSON(configs.BadRequest, utils.FormatApiError(msg, configs.BadRequest, fieldErrors))
			return
		}
		c.JSON(configs.BadRequest, utils.FormatApiError(msg, configs.BadRequest, entities.E{"message": err.Error()}))
		return
	}

	user, ok := c.Get("user")
	if !ok {
		msg := "user account not found"
		c.JSON(configs.Unauthorized, utils.FormatApiError(msg, configs.Unauthorized, entities.E{"account": msg}))
		return
	}

	u := user.(dto.UserSession)
	span.SetAttributes(attribute.String("user.id", u.ID))
	span.SetAttributes(attribute.String("user.name", u.Name))

	ps := services.NewDefaultPartService(parts.NewPartRepoDb(s.Db))
	part, apiErr := ps.GetPartById(ctx, cr.PartId)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	service := services.NewDefaultCartService(
		carts.NewCartRepoDb(s.Db),
		account.NewUserRepoDb(s.Db),
		services.NewDefaultTaxService(taxes.NewTaxRepoDb(s.Db)),
	)
	_, apiErr = service.AddCartItem(u.ID, cr, *part)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{}, "Added to cart successfully", utils.Empty(), utils.Empty()))
}

func (s *Serve) GetCarts(c *gin.Context) {
	limit, offset, err := queryParams(c, entities.QueryPathParam{Limit: 10, Offset: 0})
	if err != nil {
		msg := "bad request"
		c.JSON(configs.BadRequest, utils.FormatApiError(msg, configs.BadRequest, entities.E{"search": msg}))
		return
	}

	cs := services.NewDefaultCartService(
		carts.NewCartRepoDb(s.Db),
		account.NewUserRepoDb(s.Db),
		services.NewDefaultTaxService(taxes.NewTaxRepoDb(s.Db)),
	)
	userId, _ := c.GetQuery("user")
	status, _ := c.GetQuery("status")
	res, apiErr := cs.GetCarts(limit, offset, userId, status)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	links := next(c, entities.QueryPathParam{Limit: limit, Offset: offset})
	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "carts",
		Data:    res,
		Links:   entities.E{"next": links.Next, "previous": links.Previous},
		Meta:    entities.E{},
	})
}

func (s *Serve) GetCartByID(c *gin.Context) {
	hid := c.Param("cartID")
	if hid == "" {
		msg := "bad request"
		c.JSON(configs.BadRequest, utils.FormatApiError(msg, configs.BadRequest, entities.E{"head_id": msg}))
		return
	}

	service := services.NewDefaultCartService(
		carts.NewCartRepoDb(s.Db),
		account.NewUserRepoDb(s.Db),
		services.NewDefaultTaxService(taxes.NewTaxRepoDb(s.Db)),
	)
	resp, apiErr := service.GetCartByID(hid)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "cart detail",
		Data:    resp,
		Links:   entities.E{},
		Meta:    entities.E{},
	})
}
