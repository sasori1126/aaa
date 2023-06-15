package controllers

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/sales"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/services"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"github.com/bugsnag/bugsnag-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
)

func (s *Serve) ContactUs(c *gin.Context) {
	request := &dto.SupportRequest{}
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

	gs := services.NewDefaultGeneralService(sales.NewSaleRepoDb(s.Db))
	err = gs.SupportRequest(request)
	if err != nil {
		c.JSON(configs.ServerError,
			utils.FormatApiError("server error, failed to process request", configs.ServerError, entities.E{"message": err.Error()}))
		return
	}
	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "request send successfully",
		Data:    nil,
		Links:   nil,
		Meta:    nil,
	})
}

func (s Serve) S3Resource(c *gin.Context) {
	url, ok := c.GetQuery("url")
	if !ok {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"search": "invalid request"},
		))
		return
	}

	c.Header("Access-Control-Allow-Origin", "*")

	var client = http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		bugsnag.Notify(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	defer resp.Body.Close()

	t := resp.Header.Get("Content-Type")
	if t != "image/jpeg" && t != "image/png" && t != "image/gif" {
		bugsnag.Notify(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "invalid resource type",
		})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		bugsnag.Notify(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.Writer.Write(body)
}
