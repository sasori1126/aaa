package middlewares

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/pkg/entities"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthorizeSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := sessions.Default(c)
		sessionId := s.Get("user_id")
		if sessionId == nil {
			c.JSON(configs.Forbidden, entities.E{
				"message": "unauthorized",
			})
			c.Abort()
		}
	}
}
