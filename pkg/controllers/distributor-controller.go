package controllers

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/distributor"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/services"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func (s *Serve) UpdateDistributor(c *gin.Context) {
	hid := c.Param("distributorId")
	if hid == "" {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"head_id": "invalid request"},
		))
		return
	}
	disRequest := &dto.DistributorRequest{}
	err := c.ShouldBindJSON(disRequest)
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
	dService := services.NewDistributorService(distributor.NewRepoDistributorDb(s.Db))
	apiErr := dService.UpdateDistributorById(hid, disRequest)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}
	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"distributor updated successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) CreateDistributor(c *gin.Context) {
	disRequest := &dto.DistributorRequest{}
	err := c.ShouldBindJSON(disRequest)
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
	dService := services.NewDistributorService(distributor.NewRepoDistributorDb(s.Db))
	apiErr := dService.CreateDistributor(disRequest)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"distributor created successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) DeleteDistributor(c *gin.Context) {
	hid := c.Param("distributorId")
	if hid == "" {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"head_id": "invalid request"},
		))
		return
	}
	dService := services.NewDistributorService(distributor.NewRepoDistributorDb(s.Db))
	apiErr := dService.DeleteDistributorById(hid)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, utils.GetApiResponse(gin.H{},
		"distributor deleted successfully",
		utils.Empty(), utils.Empty()))
}

func (s *Serve) GetDistributor(c *gin.Context) {
	hid := c.Param("distributorId")
	if hid == "" {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"head_id": "invalid request"},
		))
		return
	}
	dService := services.NewDistributorService(distributor.NewRepoDistributorDb(s.Db))
	getDistributorById, apiErr := dService.GetDistributorById(hid)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "distributor detail",
		Data:    getDistributorById,
		Links:   entities.E{},
		Meta:    entities.E{},
	})
}

func (s *Serve) GetDistributors(c *gin.Context) {
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
	dService := services.NewDistributorService(distributor.NewRepoDistributorDb(s.Db))
	dists, apiErr := dService.GetDistributors(limit, offset)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	links := next(c, entities.QueryPathParam{
		Limit:  limit,
		Offset: offset,
	})
	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "list of distributors",
		Data:    dists,
		Links: entities.E{
			"next":     links.Next,
			"previous": links.Previous,
		},
		Meta: entities.E{},
	})
}
