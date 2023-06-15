package controllers

import (
	"net/http"

	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/account"
	"axis/ecommerce-backend/internal/domain/carts"
	"axis/ecommerce-backend/internal/domain/equipments"
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

func (s *Serve) UpdateUserProfile(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.JSON(configs.Unauthorized, utils.FormatApiError("user account not found", configs.Unauthorized, entities.E{
			"account": "account not found",
		}))
		return
	}
	u := user.(dto.UserSession)

	uur := &dto.UserUpdateRequest{}
	err := c.ShouldBindJSON(uur)
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

	us := services.NewDefaultUserService(account.NewUserRepoDb(s.Db))
	data := []entities.UpdateFields{
		{
			Field: "name",
			Value: uur.Name,
		},
		{
			Field: "phone_number",
			Value: uur.PhoneNumber,
		},
	}
	apiErr := us.UpdateUser(u.ID, data)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}
	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"user update successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) UserProfile(c *gin.Context) {
	sessionUser, exists := c.Get(configs.User)
	if !exists {
		s.Logout(c)
		return
	}
	u := sessionUser.(dto.UserSession)
	us := services.NewDefaultUserService(account.NewUserRepoDb(s.Db))
	user, apiErr := us.GetUser(u.ID)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}
	c.JSON(http.StatusOK, utils.GetApiResponse(user, "user detail", nil, nil))
}

func (s *Serve) GetUserEquipments(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.JSON(configs.Unauthorized, utils.FormatApiError("user account not found", configs.Unauthorized, entities.E{
			"account": "account not found",
		}))
		return
	}
	u := user.(dto.UserSession)

	es := services.NewDefaultEquipmentService(equipments.NewEquipmentRepoDb(s.Db))
	getEquipments, apiErr := es.GetUserEquipments(u.ID)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}
	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "list of user equipments",
		Data:    getEquipments,
		Links:   entities.E{},
		Meta:    entities.E{},
	})
}

func (s *Serve) UserAddresses(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.JSON(configs.Unauthorized, utils.FormatApiError("user account not found", configs.Unauthorized, entities.E{
			"account": "account not found",
		}))
		return
	}
	u := user.(dto.UserSession)
	userService := services.NewDefaultUserService(account.NewUserRepoDb(s.Db))
	response, apiErr := userService.GetAddresses(u.ID)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "user addresses",
		Data:    response,
		Links:   entities.E{},
		Meta:    entities.E{},
	})
}

func (s *Serve) DeleteAddress(c *gin.Context) {
	addressId := c.Param("addressId")
	if addressId == "" {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"head_id": "invalid request"},
		))
		return
	}

	userService := services.NewDefaultUserService(account.NewUserRepoDb(s.Db))
	apiErr := userService.DeleteUserAddresses(addressId)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"user address delete successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) AskUserToResetPassword(c *gin.Context) {
	req := &dto.AskUserToResetPasswordReq{}
	err := c.ShouldBindJSON(req)
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

	userService := services.NewDefaultUserService(account.NewUserRepoDb(s.Db))
	apiErr := userService.AskUserToResetPassword(req)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"user password reset sent",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) SetDefaultAddress(c *gin.Context) {
	addressId := c.Param("addressId")
	if addressId == "" {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"head_id": "invalid request"},
		))
		return
	}

	userService := services.NewDefaultUserService(account.NewUserRepoDb(s.Db))
	apiErr := userService.SetDefaultAddresses(addressId)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"user address update successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) UpdateAddress(c *gin.Context) {
	addressId := c.Param("addressId")
	if addressId == "" {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"head_id": "invalid request"},
		))
		return
	}

	user, ok := c.Get("user")
	if !ok {
		c.JSON(configs.Unauthorized, utils.FormatApiError("user account not found", configs.Unauthorized, entities.E{
			"account": "account not found",
		}))
		return
	}
	u := user.(dto.UserSession)

	request := &dto.UserAddressRequest{}
	err := c.ShouldBindJSON(request)
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

	userService := services.NewDefaultUserService(account.NewUserRepoDb(s.Db))
	apiErr := userService.UpdateAddress(u.ID, addressId, request)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"user address update successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) MakePayment(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.JSON(configs.Unauthorized, utils.FormatApiError("user account not found", configs.Unauthorized, entities.E{
			"account": "account not found",
		}))
		return
	}
	u := user.(dto.UserSession)

	request := &dto.PlaceOrderRequest{}
	err := c.ShouldBindJSON(request)
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
	apiErr := os.CreateUserOrder(u.ID, request)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"Order placed successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) AddAddress(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.JSON(configs.Unauthorized, utils.FormatApiError("user account not found", configs.Unauthorized, entities.E{
			"account": "account not found",
		}))
		return
	}
	u := user.(dto.UserSession)

	request := &dto.UserAddressRequest{}
	err := c.ShouldBindJSON(request)
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

	userService := services.NewDefaultUserService(account.NewUserRepoDb(s.Db))
	address, apiErr := userService.AddAddress(u.ID, request)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "created address",
		Data:    address,
		Links:   entities.E{},
		Meta:    entities.E{},
	})
}

func (s *Serve) AddEquipment(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.JSON(configs.Unauthorized, utils.FormatApiError("user account not found", configs.Unauthorized, entities.E{
			"account": "account not found",
		}))
		return
	}
	u := user.(dto.UserSession)

	er := &dto.EquipmentRequest{}
	err := c.ShouldBindJSON(er)
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

	es := services.NewDefaultEquipmentService(equipments.NewEquipmentRepoDb(s.Db))
	apiErr := es.AddEquipment(u.ID, er)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"Equipment added successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) AssignUserRoleOnAccount(c *gin.Context) {
	req := &dto.AssignUserRoleRequest{}
	err := c.ShouldBindJSON(req)
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
	userService := services.NewDefaultUserService(account.NewUserRepoDb(s.Db))
	apiErr := userService.AssignUserRole(req)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}
	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"User account updated successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) UpdateUserCurrency(c *gin.Context) {
	uor := &dto.UserUpdateCurrencyRequest{}
	err := c.ShouldBind(uor)
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
	userService := services.NewDefaultUserService(account.NewUserRepoDb(s.Db))
	apiErr := userService.UpdateUserCurrency(uor)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}
	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"User account updated successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) UpdateUserOnAccount(c *gin.Context) {
	uor := &dto.UserUpdateOnAccountRequest{}
	err := c.ShouldBind(uor)
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
	userService := services.NewDefaultUserService(account.NewUserRepoDb(s.Db))
	apiErr := userService.UpdateUserOnAccount(uor)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}
	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"User account updated successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) GetUser(c *gin.Context) {
	uid := c.Param("userId")
	if uid == "" {
		result := utils.FormatApiError("bad request", configs.BadRequest, entities.E{"user_id": "invalid request"})
		c.JSON(configs.BadRequest, result)
		return
	}

	userService := services.NewDefaultUserService(account.NewUserRepoDb(s.Db))
	user, apiErr := userService.AdminGetUser(uid)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}
	c.JSON(configs.Ok, entities.ApiResponse{Data: user, Links: entities.E{}, Message: "user detail", Meta: entities.E{}})
}

func (s *Serve) SearchUsers(c *gin.Context) {
	q, ok := c.GetQuery("q")
	if !ok {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"search": "invalid request"},
		))
		return
	}

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
	userService := services.NewDefaultUserService(account.NewUserRepoDb(s.Db))
	result, apiErr := userService.SearchUsers(limit, offset, q)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}
	links := next(c, entities.QueryPathParam{
		Limit:  limit,
		Offset: offset,
	})
	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "list of users",
		Data:    result,
		Links: entities.E{
			"next":     links.Next,
			"previous": links.Previous,
		},
		Meta: entities.E{},
	})
}

func (s *Serve) GetAllUsers(c *gin.Context) {
	defValues := entities.QueryPathParam{
		Limit:  10,
		Offset: 0,
	}
	limit, offset, err := queryParams(c, defValues)
	if err != nil {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"search": "invalid request"},
		))
		return
	}
	userService := services.NewDefaultUserService(account.NewUserRepoDb(s.Db))
	allUsers, apiErr := userService.GetUsers(limit, offset)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}
	links := next(c, entities.QueryPathParam{
		Limit:  limit,
		Offset: offset,
	})
	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "list of users",
		Data:    allUsers,
		Links: entities.E{
			"next":     links.Next,
			"previous": links.Previous,
		},
		Meta: entities.E{},
	})
}
