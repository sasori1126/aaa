package controllers

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/domain/stats"
	"axis/ecommerce-backend/internal/services"
	"axis/ecommerce-backend/internal/storage"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func (s *Serve) CacheInvalidate(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"cache": "invalid request"},
		))
		return
	}

	n, err := storage.Cache.Del(name)
	if err != nil {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"bad request",
			configs.BadRequest,
			entities.E{"cache": "invalid request"},
		))
		return
	}

	c.JSON(configs.Ok, entities.ApiResponse{
		Message: strconv.Itoa(int(n)) + " invalidated",
		Links:   entities.E{},
		Meta:    entities.E{},
	})
}

func (s *Serve) DashboardStats(c *gin.Context) {
	limit := 10
	l, ok := c.GetQuery("limit")
	if ok {
		lt, err := strconv.Atoi(l)
		if err == nil {
			limit = lt
		}
	}
	statService := services.NewDefaultStatService(stats.NewDefaultStatDbRepo(s.Db))
	stat, apiErr := statService.TotalStat(limit)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(configs.Ok, entities.ApiResponse{
		Message: "dashboard stats",
		Data:    stat,
		Links:   entities.E{},
		Meta:    entities.E{},
	})
}
