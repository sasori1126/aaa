package controllers

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/controls"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/services"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func (s *Serve) GetControllers(c *gin.Context) {
	cs := services.NewDefaultControllerService(controls.NewControllerRepoDb(s.Db))
	getControllers, apiErr := cs.GetControllers(10, 0)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "list of controllers",
		Data:    getControllers,
		Links:   entities.E{},
		Meta:    entities.E{},
	})
}

func (s *Serve) DeleteController(c *gin.Context) {
	hid := c.Param("controllerId")
	if hid == "" {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"head_id": "invalid request"},
		))
		return
	}
	cs := services.NewDefaultControllerService(controls.NewControllerRepoDb(s.Db))
	apiErr := cs.DeleteController(hid)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}
	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"controller deleted successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) UpdateController(c *gin.Context) {
	hid := c.Param("controllerId")
	if hid == "" {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"controller_id": "invalid request"},
		))
		return
	}

	request := &dto.ControllerRequest{}
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

	cs := services.NewDefaultControllerService(controls.NewControllerRepoDb(s.Db))
	apiErr := cs.UpdateController(hid, request)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"controller updated successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) CreateController(c *gin.Context) {
	cr := &dto.ControllerRequest{}
	err := c.ShouldBindJSON(cr)
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

	cs := services.NewDefaultControllerService(controls.NewControllerRepoDb(s.Db))
	apiErr := cs.CreateController(cr)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"controller created successfully",
		utils.Empty(), utils.Empty()))
}
