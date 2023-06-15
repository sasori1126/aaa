package controllers

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/categories"
	"axis/ecommerce-backend/internal/domain/man-models"
	"axis/ecommerce-backend/internal/services"
	"axis/ecommerce-backend/pkg/entities"
	"github.com/gin-gonic/gin"
)

func (s *Serve) GetPartss(c *gin.Context) {
	cs := services.NewDefaultCatService(categories.NewCategoryRepoDb(s.Db))
	cats, apiErr := cs.GetCategories(10, 0)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	//get getModels
	ms := services.NewDefaultModelService(man_models.NewModelRepoDb(s.Db))
	getModels, apiErr := ms.GetModels(10, 0)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "list of parts, categories",
		Data: entities.E{
			"parts":           []string{"demo", "demo1", "demo3"},
			"categories":      cats,
			"user_equipments": []string{"equip1", "equip2"},
			"man-models":      getModels,
		},
		Meta:  entities.E{},
		Links: entities.E{},
	})
}
