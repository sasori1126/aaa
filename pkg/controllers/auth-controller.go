package controllers

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/actions"
	"axis/ecommerce-backend/internal/domain/account"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/services"
	"axis/ecommerce-backend/internal/storage"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Serve) RegisterUser(c *gin.Context) {
	registerReq := &dto.RegisterRequest{}
	err := c.ShouldBindJSON(registerReq)
	if err != nil {
		ers, ok := err.(validator.ValidationErrors)
		if ok {
			reqErrors := utils.FormatValidationError(ers)
			c.JSON(http.StatusBadRequest, utils.FormatApiError(
				"failed to process user request",
				http.StatusBadRequest,
				reqErrors,
			))
		} else {
			c.JSON(http.StatusBadRequest, utils.FormatApiError(
				"failed to process user request",
				http.StatusBadRequest, entities.E{"binding": "failed to bind request"},
			))
		}
		return
	}

	userService := services.NewDefaultUserService(account.NewUserRepoDb(s.Db))
	apiErr := userService.CreateUser(registerReq)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(configs.Created, utils.GetApiResponse(gin.H{}, "user created successfully", utils.Empty(), utils.Empty()))
}

func (s *Serve) Logout(c *gin.Context) {
	user, exists := c.Get(configs.User)
	if !exists {
		c.JSON(configs.NotFound, utils.GeneralError("invalid user session", configs.NotFound))
		return
	}

	u := user.(dto.UserSession)

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

	c.JSON(http.StatusAccepted, utils.GetApiResponse(nil, "log out successfully", nil, nil))
}

func (s *Serve) Refresh(c *gin.Context) {
	rr := &dto.RefreshTokenRequest{}
	err := c.ShouldBindJSON(rr)
	if err != nil {
		ers, ok := err.(validator.ValidationErrors)
		if ok {
			fieldErs := utils.FormatValidationError(ers)
			c.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, fieldErs))
			return
		} else {
			c.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, entities.E{}))
			return
		}
	}

	tkn, err := actions.ValidateRefreshToken(rr.RefreshToken)
	if err != nil {
		c.JSON(configs.Forbidden, utils.FormatApiError("refresh token expired", configs.Forbidden, entities.E{}))
		return
	}

	if !tkn.Valid {
		c.JSON(configs.Forbidden, utils.FormatApiError("invalid refresh token", configs.Forbidden, entities.E{}))
		return
	}

	claims := tkn.Claims.(jwt.MapClaims)
	userJwt := claims["User"].(map[string]interface{})

	userId := userJwt["UserId"].(string)
	refreshUid, ok := claims["Uid"].(string)
	if !ok {
		c.JSON(configs.ServerError, utils.FormatApiError("failed to validate refresh token", configs.ServerError, entities.E{}))
		return
	}

	_, err = storage.Cache.Get(refreshUid)
	if err != nil {
		c.JSON(configs.NotFound, utils.FormatApiError("refresh token not found", configs.NotFound, entities.E{}))
		return
	}

	_, err = storage.Cache.Del(refreshUid)
	if err != nil {
		c.JSON(configs.ServerError, utils.FormatApiError("failed to validate refresh token", configs.ServerError, entities.E{}))
		return
	}

	userService := services.NewDefaultUserService(account.NewUserRepoDb(s.Db))
	user, err := userService.FindUserId(userId)
	if err != nil {
		c.JSON(configs.ServerError, utils.FormatApiError("failed to retrieve user account", configs.ServerError, entities.E{}))
		return
	}

	td, err := userService.GetUserToken(user)
	if err != nil {
		c.JSON(configs.ServerError, utils.FormatApiError("failed to retrieve user account", configs.ServerError, entities.E{}))
		return
	}
	res := dto.LoginResponse{
		AccessToken:  td.AccessToken,
		RefreshToken: td.RefreshToken,
		AtExpires:    td.AtExpires,
		RtExpires:    td.RtExpires,
	}

	c.JSON(configs.Ok, res)
}

func (s *Serve) RequestNewEmailVerificationToken(c *gin.Context) {
	er := &dto.EmailVerificationTokenRequest{}
	err := c.ShouldBindJSON(er)
	if err != nil {
		ers, ok := err.(validator.ValidationErrors)
		if ok {
			fieldErs := utils.FormatValidationError(ers)
			c.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, fieldErs))
			return
		} else {
			c.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, entities.E{}))
			return
		}
	}

	us := services.NewDefaultUserService(account.NewUserRepoDb(s.Db))
	apiErr := us.ResendEmailVerificationToken(er.Email)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusAccepted, entities.ApiResponse{
		Message: "token send successfully",
		Data:    entities.E{},
	})
}

func (s *Serve) ValidateAccountEmail(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(configs.BadRequest, utils.FormatApiError(
			"invalid request",
			configs.BadRequest,
			entities.E{"token": "invalid token request"},
		))
		return
	}

	getEmail, err := storage.Cache.Get(token)
	if err != nil {
		c.JSON(configs.NotFound, utils.FormatApiError(
			"validation token expired",
			configs.NoContent,
			entities.E{"token": "validation token expired"},
		))
		return
	}

	userService := services.NewDefaultUserService(account.NewUserRepoDb(s.Db))
	user, apiErr := userService.GetUserByEmail(getEmail)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	verifiedAt := time.Now().Format(time.RFC3339)

	data := []entities.UpdateFields{{Field: "email_verified_at", Value: verifiedAt}}
	apiErr = userService.UpdateUser(user.Id, data)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	_, err = storage.Cache.Del(token, user.Id)
	if err != nil {
		configs.Logger.Warn("failed to delete cache key")
	}

	c.JSON(http.StatusOK, entities.ApiResponse{
		Message: "account verified successfully",
		Data:    entities.E{"email": getEmail}})
}

func (s *Serve) ResetPassword(c *gin.Context) {
	pr := &dto.ResetPasswordRequest{}
	err := c.ShouldBindJSON(pr)
	if err != nil {
		validatorErrs, ok := err.(validator.ValidationErrors)
		if ok {
			formatErrs := utils.FormatValidationError(validatorErrs)
			c.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, formatErrs))
		}

		return
	}

	//confirm if password is eq to confirm_password
	if pr.Password != pr.ConfirmPassword {
		c.JSON(
			http.StatusBadRequest,
			utils.FormatApiError("confirm password is not equal to password",
				http.StatusBadRequest, entities.E{"confirm_password": "confirm password is not equal to password"}),
		)
		return
	}

	getEmail, err := storage.Cache.Get(pr.Token)
	if err != nil {
		c.JSON(configs.NotFound, utils.FormatApiError(
			"password reset token expired",
			configs.NoContent,
			entities.E{"token": "password reset  token expired"},
		))
		return
	}

	userService := services.NewDefaultUserService(account.NewUserRepoDb(s.Db))
	user, apiErr := userService.GetUserByEmail(getEmail)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	newPassword, err := actions.PasswordBcrypt(pr.Password)
	if err != nil {
		c.JSON(configs.NotFound, utils.FormatApiError(
			"password reset token expired",
			configs.NoContent,
			entities.E{"token": "password reset  token expired"},
		))
		return
	}

	data := []entities.UpdateFields{{Field: "password", Value: newPassword}}
	apiErr = userService.UpdateUser(user.Id, data)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	_, err = storage.Cache.Del(pr.Token)
	if err != nil {
		configs.Logger.Warn("failed to delete cache key")
	}

	c.JSON(http.StatusOK, entities.ApiResponse{
		Message: "password reset successfully",
		Data:    entities.E{"password": "password reset successfully"}})
}

func (s *Serve) AccountResetPassword(c *gin.Context) {
	pr := &dto.AccountResetPasswordRequest{}
	err := c.ShouldBindJSON(pr)
	if err != nil {
		validatorErrs, ok := err.(validator.ValidationErrors)
		if ok {
			formatErrs := utils.FormatValidationError(validatorErrs)
			c.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, formatErrs))
		}

		return
	}

	//confirm if password is eq to confirm_password
	if pr.Password != pr.ConfirmPassword {
		c.JSON(
			http.StatusBadRequest,
			utils.FormatApiError("confirm password is not equal to password",
				http.StatusBadRequest, entities.E{"confirm_password": "confirm password is not equal to password"}),
		)
		return
	}

	userService := services.NewDefaultUserService(account.NewUserRepoDb(s.Db))
	user, apiErr := userService.GetUser(pr.UserId)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	newPassword, err := actions.PasswordBcrypt(pr.Password)
	if err != nil {
		c.JSON(configs.NotFound, utils.FormatApiError(
			"password reset token expired",
			configs.NoContent,
			entities.E{"token": "password reset  token expired"},
		))
		return
	}

	data := []entities.UpdateFields{{Field: "password", Value: newPassword}}
	apiErr = userService.UpdateUser(user.Id, data)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusOK, entities.ApiResponse{
		Message: "password reset successfully",
		Data:    entities.E{"password": "password reset successfully"}})
}

func (s *Serve) ForgotPassword(c *gin.Context) {
	fp := &dto.ForgotPassword{}
	err := c.ShouldBindJSON(fp)
	if err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			fieldsErr := utils.FormatValidationError(validationErrors)
			c.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, fieldsErr))
		} else {
			c.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, entities.E{}))
		}
		return
	}
	us := services.NewDefaultUserService(account.NewUserRepoDb(s.Db))
	apiErr := us.ForgotPassword(fp)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	c.JSON(http.StatusAccepted, utils.GetApiResponse(nil, "Password reset email sent to your email", nil, nil))
}

func (s *Serve) Login(c *gin.Context) {
	lr := &dto.LoginRequest{}
	err := c.ShouldBindJSON(lr)
	if err != nil {
		vErrs, ok := err.(validator.ValidationErrors)
		if ok {
			fieldErrs := utils.FormatValidationError(vErrs)
			c.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, fieldErrs))
		} else {
			c.JSON(configs.BadRequest, utils.FormatApiError("bad request, failed to process request", configs.BadRequest, entities.E{}))
		}
		return
	}
	userService := services.NewDefaultUserService(account.NewUserRepoDb(s.Db))
	td, apiErr := userService.UserLogin(lr)
	if apiErr != nil {
		c.JSON(apiErr.Code, apiErr)
		return
	}

	res := dto.LoginResponse{
		AccessToken:  td.AccessToken,
		RefreshToken: td.RefreshToken,
		AtExpires:    td.AtExpires,
		RtExpires:    td.RtExpires,
	}

	c.JSON(http.StatusOK, res)
}
