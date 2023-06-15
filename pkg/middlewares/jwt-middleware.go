package middlewares

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/actions"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/services"
	"axis/ecommerce-backend/internal/storage"
	"axis/ecommerce-backend/pkg/utils"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
)

func AuthorizeJwt(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		const AuthorizationHeader = "Authorization"
		bearerToken := c.GetHeader(AuthorizationHeader)
		strArr := strings.Split(bearerToken, " ")

		if len(strArr) != 2 {
			c.JSON(configs.Unauthorized, utils.GeneralError(
				"unauthorized",
				configs.Unauthorized,
			))
			c.Abort()
		} else {
			authorization := strArr[0]
			if authorization != "Bearer" {
				configs.Logger.Error(errors.New("wrong authorization"))
				c.JSON(configs.Unauthorized, utils.GeneralError(
					"invalid authorization method",
					configs.Unauthorized,
				))
				c.Abort()
				return
			}

			token := strArr[1]
			verified, err := actions.ValidateToken(token)
			if err != nil {
				configs.Logger.Error(err)
				c.JSON(configs.Unauthorized, utils.GeneralError(
					err.Error(),
					configs.Unauthorized,
				))
				c.Abort()
				return
			}

			if verified.Valid {
				claims := verified.Claims.(jwt.MapClaims)
				sessionUuid := claims["Uid"].(string)
				_, err := storage.Cache.Get(sessionUuid)
				if err != nil {
					c.JSON(configs.Unauthorized, utils.GeneralError(
						"token expired",
						configs.Unauthorized,
					))
					c.Abort()
					return
				}
				userJwt := claims["User"].(map[string]interface{})
				userId := userJwt["UserId"].(string)
				asid := userJwt["Asid"].(string)
				rsid := userJwt["Rsid"].(string)

				user, err := userService.FindUserId(userId)
				if err != nil {
					configs.Logger.Error(err)
					c.JSON(configs.ServerError, utils.GeneralError(
						"failed saving user session",
						configs.ServerError,
					))
					c.Abort()
					return
				}

				userSession := dto.UserSession{
					ID:               userId,
					Name:             user.Name,
					Email:            user.Email,
					PhoneNumber:      user.PhoneNumber,
					IsActive:         user.IsActive,
					AccessSessionId:  asid,
					RefreshSessionId: rsid,
					EmailVerifiedAt:  user.EmailVerifiedAt.Time,
					CreatedAt:        user.CreatedAt,
					UpdatedAt:        user.UpdatedAt,
					Role:             user.Role,
				}

				c.Set(configs.User, userSession)
				c.Next()
			} else {
				configs.Logger.Error(err)
				c.JSON(configs.Unauthorized, utils.GeneralError(
					"invalid token",
					configs.Unauthorized,
				))
				c.Abort()
				return
			}
		}
	}
}
