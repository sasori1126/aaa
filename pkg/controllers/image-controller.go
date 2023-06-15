package controllers

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/images"
	"axis/ecommerce-backend/internal/services"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Serve) DeleteImage(c *gin.Context) {
	hid := c.Param("imageId")
	if hid == "" {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"model_id": "invalid request"},
		))
		return
	}
	ms := services.NewDefaultImageService(images.NewImageRepoDb(s.Db))
	apiErr := ms.DeleteImage(hid)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}
	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"image deleted successfully",
		utils.Empty(), utils.Empty()))
}
