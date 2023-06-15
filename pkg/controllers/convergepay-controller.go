package controllers

import (
	"net/http"

	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/services"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (s *Serve) GetConvergePayToken(c *gin.Context) {
	request := &dto.GetConvergePayTokenRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		msg := "bad request, failed to process request"
		if vErrs, ok := err.(validator.ValidationErrors); ok {
			c.JSON(configs.BadRequest, utils.FormatApiError(msg, configs.BadRequest, utils.FormatValidationError(vErrs)))
			return
		}
		c.JSON(configs.BadRequest, utils.FormatApiError(msg, configs.BadRequest, entities.E{"message": err.Error()}))
		return
	}

	if _, ok := c.Get("user"); !ok {
		msg := "user account not found"
		c.JSON(configs.Unauthorized, utils.FormatApiError(msg, configs.Unauthorized, entities.E{"account": msg}))
		return
	}

	token, apiError := services.GetPaymentToken(request)
	if apiError != nil {
		c.JSON(apiError.Code, apiError)
		return
	}

	response := utils.GetApiResponse(token, "", utils.Empty(), utils.Empty())
	c.JSON(http.StatusOK, response)
}
