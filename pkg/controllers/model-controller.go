package controllers

import (
	"axis/ecommerce-backend/configs"
	man_models "axis/ecommerce-backend/internal/domain/man-models"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/services"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func (s *Serve) DeleteModel(c *gin.Context) {
	hid := c.Param("modelId")
	if hid == "" {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"model_id": "invalid request"},
		))
		return
	}
	ms := services.NewDefaultModelService(man_models.NewModelRepoDb(s.Db))
	apiErr := ms.DeleteModel(hid)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}
	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"model deleted successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) GetModel(c *gin.Context) {
	hid := c.Param("modelId")
	if hid == "" {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"model_id": "invalid request"},
		))
		return
	}
	ms := services.NewDefaultModelService(man_models.NewModelRepoDb(s.Db))
	getDistributorById, apiErr := ms.GetModel(hid)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "model detail",
		Data:    getDistributorById,
		Links:   entities.E{},
		Meta:    entities.E{},
	})
}

func (s *Serve) GetModels(c *gin.Context) {
	defValues := entities.QueryPathParam{
		Limit:  -1,
		Offset: -1,
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
	ms := services.NewDefaultModelService(man_models.NewModelRepoDb(s.Db))
	getModels, apiErr := ms.GetModels(limit, offset)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	links := next(c, entities.QueryPathParam{
		Limit:  limit,
		Offset: offset,
	})
	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "list of models",
		Data:    getModels,
		Links: entities.E{
			"next":     links.Next,
			"previous": links.Previous,
		},
		Meta: entities.E{},
	})
}

func (s *Serve) UpdateModel(c *gin.Context) {
	hid := c.Param("modelId")
	if hid == "" {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"model_id": "invalid request"},
		))
		return
	}

	mr := &dto.ModelRequest{}
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
	ms := services.NewDefaultModelService(man_models.NewModelRepoDb(s.Db))
	apiErr := ms.UpdateModel(hid, mr)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"model update successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) CreateModel(c *gin.Context) {
	mr := &dto.ModelRequest{}
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

	ms := services.NewDefaultModelService(man_models.NewModelRepoDb(s.Db))
	apiErr := ms.CreateModel(mr)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"model created successfully",
		utils.Empty(), utils.Empty()))
}
