package middlewares

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/storage"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"net/http"
)

func NoRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(configs.NotFound, &entities.ApiError{
			Code:    configs.NotFound,
			Message: "route not found",
			Errors:  entities.E{},
		})
	}
}

func UserRole(allowedRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionUser, ok := c.Get(configs.User)
		if !ok {
			c.JSON(configs.Forbidden, entities.E{
				"message": "unauthorized",
			})
			c.AbortWithStatus(configs.Forbidden)
		}

		u := sessionUser.(dto.UserSession)
		ok = slices.Contains(allowedRoles, u.Role)
		if !ok {
			_, err := storage.Cache.Del(u.AccessSessionId)
			if err != nil {
				c.JSON(configs.ServerError, utils.FormatApiError(
					"failed to clear user session",
					configs.ServerError,
					entities.E{},
				))
				return
			}

			_, err = storage.Cache.Del(u.RefreshSessionId)
			if err != nil {
				c.JSON(configs.ServerError, utils.FormatApiError(
					"failed to clear user session",
					configs.ServerError,
					entities.E{},
				))
				return
			}

			c.JSON(http.StatusForbidden, &entities.ApiError{
				Code:    http.StatusForbidden,
				Message: "You do not have permission to access this route",
				Errors:  entities.E{},
			})
			c.AbortWithStatus(http.StatusForbidden)
		}
	}
}

func NoMethod() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, &entities.ApiError{
			Code:    http.StatusMethodNotAllowed,
			Message: "method not allowed",
			Errors:  entities.E{},
		})
	}
}
