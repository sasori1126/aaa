package controllers

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/manufacturers"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/services"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func (s *Serve) DeleteManufacturer(c *gin.Context) {
	hid := c.Param("manufacturerId")
	if hid == "" {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"head_id": "invalid request"},
		))
		return
	}
	ms := services.NewDefaultManufacturerService(manufacturers.NewManufacturerRepoDb(s.Db))
	apiErr := ms.DeleteManufacturerById(hid)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}
	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"manufacturer deleted successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) UpdateManufacturer(c *gin.Context) {
	hid := c.Param("manufacturerId")
	if hid == "" {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"head_id": "invalid request"},
		))
		return
	}

	request := &dto.ManufacturerRequest{}
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
	ms := services.NewDefaultManufacturerService(manufacturers.NewManufacturerRepoDb(s.Db))
	apiErr := ms.UpdateManufacturerById(hid, request)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}
	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"manufacturer updated successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) CreateManufacturer(c *gin.Context) {
	mr := &dto.ManufacturerRequest{}
	err := c.ShouldBindJSON(mr)
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

	ms := services.NewDefaultManufacturerService(manufacturers.NewManufacturerRepoDb(s.Db))
	apiErr := ms.CreateManufacturer(mr)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"manufacturer created successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) GetManufacturers(c *gin.Context) {
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

	ms := services.NewDefaultManufacturerService(manufacturers.NewManufacturerRepoDb(s.Db))
	getManufacturers, apiErr := ms.GetManufacturers(limit, offset)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	links := next(c, entities.QueryPathParam{
		Limit:  limit,
		Offset: offset,
	})
	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "list of manufacturers",
		Data:    getManufacturers,
		Links: entities.E{
			"next":     links.Next,
			"previous": links.Previous,
		},
		Meta: entities.E{},
	})
}
