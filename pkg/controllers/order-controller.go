package controllers

import (
	"log"
	"net/http"

	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/account"
	"axis/ecommerce-backend/internal/domain/carts"
	"axis/ecommerce-backend/internal/domain/orders"
	"axis/ecommerce-backend/internal/domain/payments"
	"axis/ecommerce-backend/internal/domain/taxes"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/services"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (s *Serve) UserOrders(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.JSON(configs.Unauthorized, utils.FormatApiError("user account not found", configs.Unauthorized, entities.E{
			"account": "account not found",
		}))
		return
	}
	u := user.(dto.UserSession)
	os := services.NewDefaultOrderService(
		orders.NewOrderRepoDb(s.Db),
		carts.NewCartRepoDb(s.Db),
		account.NewUserRepoDb(s.Db),
		services.NewDefaultTaxService(taxes.NewTaxRepoDb(s.Db)),
		payments.NewPaymentRepoDb(s.Db),
	)
	res, apiErr := os.GetUserOrders(u.ID)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "user orders",
		Data:    res,
		Links:   entities.E{},
		Meta:    entities.E{},
	})
}

func (s *Serve) AddOrderPayment(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.JSON(configs.Unauthorized, utils.FormatApiError("user account not found", configs.Unauthorized, entities.E{
			"account": "account not found",
		}))
		return
	}
	u := user.(dto.UserSession)
	paymentReq := &dto.AddOrderPayment{}
	err := c.ShouldBindJSON(paymentReq)
	if err != nil {
		vErrs, ok := err.(validator.ValidationErrors)
		if ok {
			fieldErrors := utils.FormatValidationError(vErrs)
			c.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, fieldErrors))
		} else {
			c.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, entities.E{"message": err.Error()}))
		}
		return
	}

	os := services.NewDefaultOrderService(
		orders.NewOrderRepoDb(s.Db),
		carts.NewCartRepoDb(s.Db),
		account.NewUserRepoDb(s.Db),
		services.NewDefaultTaxService(taxes.NewTaxRepoDb(s.Db)),
		payments.NewPaymentRepoDb(s.Db),
	)
	apiErr := os.AddOrderPayment(u.ID, paymentReq)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"order payment added successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) OrderUpdateStatus(c *gin.Context) {
	hid := c.Param("orderID")
	if hid == "" {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"order id": "invalid request"},
		))
		return
	}

	or := &dto.OrderUpdateStatusRequest{}
	err := c.ShouldBindJSON(or)
	if err != nil {
		vErrs, ok := err.(validator.ValidationErrors)
		if ok {
			fieldErrors := utils.FormatValidationError(vErrs)
			c.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, fieldErrors))
		} else {
			c.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, entities.E{"message": err.Error()}))
		}
		return
	}

	os := services.NewDefaultOrderService(
		orders.NewOrderRepoDb(s.Db),
		carts.NewCartRepoDb(s.Db),
		account.NewUserRepoDb(s.Db),
		services.NewDefaultTaxService(taxes.NewTaxRepoDb(s.Db)),
		payments.NewPaymentRepoDb(s.Db),
	)
	apiErr := os.OrderUpdateStatus(hid, *or)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"order updated successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) GetOrderByID(c *gin.Context) {
	hid := c.Param("orderID")
	if hid == "" {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"head_id": "invalid request"},
		))
		return
	}

	or := services.NewDefaultOrderService(
		orders.NewOrderRepoDb(s.Db),
		carts.NewCartRepoDb(s.Db),
		account.NewUserRepoDb(s.Db),
		services.NewDefaultTaxService(taxes.NewTaxRepoDb(s.Db)),
		payments.NewPaymentRepoDb(s.Db),
	)

	resp, apiErr := or.GetOrderByID(hid)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "order detail",
		Data:    resp,
		Links:   entities.E{},
		Meta:    entities.E{},
	})
}

func (s *Serve) GetOrders(c *gin.Context) {
	qp := entities.QueryPathParam{
		Limit:  10,
		Offset: 0,
	}

	limit, offset, err := queryParams(c, qp)
	if err != nil {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"search": "invalid request"},
		))
		return
	}

	os := services.NewDefaultOrderService(
		orders.NewOrderRepoDb(s.Db),
		carts.NewCartRepoDb(s.Db),
		account.NewUserRepoDb(s.Db),
		services.NewDefaultTaxService(taxes.NewTaxRepoDb(s.Db)),
		payments.NewPaymentRepoDb(s.Db),
	)

	userId, _ := c.GetQuery("user")
	res, apiErr := os.GetOrders(limit, offset, userId)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	links := next(c, entities.QueryPathParam{
		Limit:  limit,
		Offset: offset,
	})

	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "orders",
		Data:    res,
		Links: entities.E{
			"next":     links.Next,
			"previous": links.Previous,
		},
		Meta: entities.E{},
	})
}

func (s *Serve) SearchOrders(c *gin.Context) {
	qp := entities.QueryPathParam{
		Limit:  10,
		Offset: 0,
	}

	limit, offset, err := queryParams(c, qp)
	if err != nil {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"search": "invalid request"},
		))
		return
	}

	q, ok := c.GetQuery("q")
	if !ok {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"search": "invalid request"},
		))
		return
	}

	os := services.NewDefaultOrderService(
		orders.NewOrderRepoDb(s.Db),
		carts.NewCartRepoDb(s.Db),
		account.NewUserRepoDb(s.Db),
		services.NewDefaultTaxService(taxes.NewTaxRepoDb(s.Db)),
		payments.NewPaymentRepoDb(s.Db),
	)

	log.Println(os, q, limit, offset)
}
