package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/taxes"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/internal/services"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (s *Serve) CreateTax(c *gin.Context) {
	request := &dto.TaxRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		message := "bad request, failed to process request"
		if vErrs, ok := err.(validator.ValidationErrors); !ok {
			c.JSON(configs.BadRequest, utils.FormatApiError(message, configs.BadRequest, entities.E{"message": err.Error()}))
		} else {
			fieldErrors := utils.FormatValidationError(vErrs)
			c.JSON(configs.BadRequest, utils.FormatApiError(message, configs.BadRequest, fieldErrors))
		}
		return
	}

	service := services.NewDefaultTaxService(taxes.NewTaxRepoDb(s.Db))
	if apiErr := service.AddTax(request); apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}
	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{}, "tax created successfully", utils.Empty(), utils.Empty()))
}

func (s *Serve) GetTaxesForTaxExemptions(c *gin.Context) {
	userHid := c.Param("userId")
	if userHid == "" {
		msg := "Missing userId"
		c.JSON(configs.BadRequest, utils.FormatApiError(msg, configs.BadRequest, entities.E{"state": msg}))
		return
	}
	userId, err := models.DecodeHashId(userHid)
	if err != nil {
		msg := "Invalid user id"
		c.JSON(configs.BadRequest, utils.FormatApiError(msg, configs.BadRequest, entities.E{"state": msg}))
		return
	}

	service := services.NewDefaultTaxService(taxes.NewTaxRepoDb(s.Db))

	taxesResult := []dto.TaxLine{}
	taxExemptionsResult := []dto.TaxLine{}
	taxes, taxExemptions, err := service.GetTaxesForTaxExemptions(userId)
	if err != nil {
		apiErr := utils.FormatApiError("failed to get taxes", configs.ServerError, entities.E{"error": err.Error()})
		c.JSON(apiErr.Code, apiErr)
		return
	}

	for i := range taxes {
		tax := taxes[i]
		taxesResult = append(taxesResult, dto.TaxLine{
			Country: tax.Country,
			Id:      fmt.Sprint(tax.ID),
			Name:    tax.Name,
			Rate:    tax.Rate,
			State:   tax.State,
		})
	}
	for i := range taxExemptions {
		tax := taxExemptions[i]
		taxExemptionsResult = append(taxExemptionsResult, dto.TaxLine{
			Country: tax.Country,
			Id:      fmt.Sprint(tax.ID),
			Name:    tax.Name,
			Rate:    tax.Rate,
			State:   tax.State,
		})
	}
	c.JSON(configs.Ok, entities.ApiResponse{
		Data:    map[string][]dto.TaxLine{"taxes": taxesResult, "tax_exemptions": taxExemptionsResult},
		Message: "List of taxes by country and state",
		Meta:    entities.E{},
		Links:   entities.E{},
	})
}

func (s *Serve) SaveTaxExemption(c *gin.Context) {
	userHid := c.Param("userId")
	if userHid == "" {
		msg := "Missing userId"
		c.JSON(configs.BadRequest, utils.FormatApiError(msg, configs.BadRequest, entities.E{"state": msg}))
		return
	}
	userId, err := models.DecodeHashId(userHid)
	if err != nil {
		msg := "Invalid user id"
		c.JSON(configs.BadRequest, utils.FormatApiError(msg, configs.BadRequest, entities.E{"state": msg}))
		return
	}
	taxId := c.Param("taxId")
	if taxId == "" {
		msg := "Missing taxId"
		c.JSON(configs.BadRequest, utils.FormatApiError(msg, configs.BadRequest, entities.E{"taxId": msg}))
		return
	}
	parsedTaxId, err := strconv.ParseUint(taxId, 10, 64)
	if err != nil {
		msg := "invalid tax id"
		c.JSON(configs.BadRequest, utils.FormatApiError(msg, configs.BadRequest, entities.E{"taxId": msg}))
		return
	}

	saveRequest := models.TaxExemption{TaxId: uint(parsedTaxId), UserId: userId}
	service := services.NewDefaultTaxService(taxes.NewTaxRepoDb(s.Db))
	if apiErr := service.SaveTaxExemption(saveRequest); apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}
	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{}, "tax exemptions saved", utils.Empty(), utils.Empty()))
}

func (s *Serve) DeleteTaxExemption(c *gin.Context) {
	userHid := c.Param("userId")
	if userHid == "" {
		msg := "Missing userId"
		c.JSON(configs.BadRequest, utils.FormatApiError(msg, configs.BadRequest, entities.E{"state": msg}))
		return
	}
	userId, err := models.DecodeHashId(userHid)
	if err != nil {
		msg := "Invalid user id"
		c.JSON(configs.BadRequest, utils.FormatApiError(msg, configs.BadRequest, entities.E{"state": msg}))
		return
	}
	taxId := c.Param("taxId")
	if taxId == "" {
		msg := "Missing taxId"
		c.JSON(configs.BadRequest, utils.FormatApiError(msg, configs.BadRequest, entities.E{"taxId": msg}))
		return
	}
	parsedTaxId, err := strconv.ParseUint(taxId, 10, 64)
	if err != nil {
		msg := "invalid tax id"
		c.JSON(configs.BadRequest, utils.FormatApiError(msg, configs.BadRequest, entities.E{"taxId": msg}))
		return
	}

	deleteRequest := models.TaxExemption{TaxId: uint(parsedTaxId), UserId: userId}
	service := services.NewDefaultTaxService(taxes.NewTaxRepoDb(s.Db))
	if apiErr := service.DeleteTaxExemption(deleteRequest); apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}
	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{}, "tax exemptions saved", utils.Empty(), utils.Empty()))
}
