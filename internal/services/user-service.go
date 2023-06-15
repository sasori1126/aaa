package services

import (
	"encoding/json"
	"errors"
	"time"

	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal"
	"axis/ecommerce-backend/internal/actions"
	"axis/ecommerce-backend/internal/domain/account"
	"axis/ecommerce-backend/internal/dto"
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/internal/notification/mail"
	"axis/ecommerce-backend/internal/storage"
	"axis/ecommerce-backend/pkg/entities"
	"axis/ecommerce-backend/pkg/utils"

	"github.com/bugsnag/bugsnag-go/v2"
	"gorm.io/gorm"
)

type UserService interface {
	AddAddress(userHid string, request *dto.UserAddressRequest) (*dto.UserAddressResponse, *entities.ApiError)
	AdminGetUser(userId string) (*dto.UserResponse, *entities.ApiError)
	AskUserToResetPassword(request *dto.AskUserToResetPasswordReq) *entities.ApiError
	AssignUserRole(request *dto.AssignUserRoleRequest) *entities.ApiError
	CreateUser(userRequest *dto.RegisterRequest) *entities.ApiError
	DeleteUserAddresses(addressId string) *entities.ApiError
	FindUserId(id string) (*models.User, error)
	ForgotPassword(fp *dto.ForgotPassword) *entities.ApiError
	GetAddresses(userHid string) ([]dto.UserAddressResponse, *entities.ApiError)
	GetUser(userId string) (*dto.UserResponse, *entities.ApiError)
	GetUserByEmail(email string) (*dto.UserResponse, *entities.ApiError)
	GetUserToken(user *models.User) (*actions.TokenDetail, error)
	GetUsers(limit, offset int) ([]dto.UserResponse, *entities.ApiError)
	ResendEmailVerificationToken(email string) *entities.ApiError
	SearchUsers(limit, offset int, search string) ([]dto.UserResponse, *entities.ApiError)
	SetDefaultAddresses(addressId string) *entities.ApiError
	UpdateAddress(userHid string, addressId string, request *dto.UserAddressRequest) *entities.ApiError
	UpdateUser(userId string, updates []entities.UpdateFields) *entities.ApiError
	UpdateUserCurrency(request *dto.UserUpdateCurrencyRequest) *entities.ApiError
	UpdateUserOnAccount(request *dto.UserUpdateOnAccountRequest) *entities.ApiError
	UserLogin(ur *dto.LoginRequest) (*actions.TokenDetail, *entities.ApiError)
}

type DefaultUserService struct {
	repo account.UserRepo
}

func NewDefaultUserService(repo account.UserRepo) UserService {
	return DefaultUserService{repo: repo}
}

func (s DefaultUserService) AddAddress(
	userHid string,
	request *dto.UserAddressRequest,
) (*dto.UserAddressResponse, *entities.ApiError) {
	userId, err := models.DecodeHashId(userHid)
	if err != nil {
		bugsnag.Notify(err)
		errMsg := "invalid user id"
		return nil, utils.FormatApiError(errMsg, configs.BadRequest, entities.E{"user_id": errMsg})
	}
	user, err := s.repo.FindUserById(userId)
	if err != nil {
		bugsnag.Notify(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errMsg := "User not found"
			return nil, utils.FormatApiError(errMsg, configs.NotFound, entities.E{"user_id": errMsg})
		}
		return nil, utils.FormatApiError("Internal server error", configs.ServerError, entities.E{"error": err.Error()})
	}

	address := models.Address{
		City:       request.City,
		Country:    request.Country,
		Province:   request.Province,
		State:      request.State,
		StreetName: request.StreetName,
		ZipCode:    request.ZipCode,
	}
	contact := models.Contact{
		Email:        request.Email,
		GormModel:    configs.GormModel{},
		Name:         request.Name,
		Organisation: request.Organisation,
		Phone:        request.Phone,
	}
	userAddress := models.UserAddress{
		Address:          address,
		Contact:          contact,
		GormModel:        configs.GormModel{},
		IsDefaultAddress: true,
		Type:             internal.ShippingAddress,
		UserId:           user.ID,
		User:             *user,
	}
	createdAddress, err := s.repo.AddUserAddress(userAddress)
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError("Internal server error", configs.ServerError, entities.E{"error": err.Error()})
	}
	response := createdAddress.ToResponse()
	return &response, nil
}

func (s DefaultUserService) AdminGetUser(userId string) (*dto.UserResponse, *entities.ApiError) {
	ui, err := models.DecodeHashId(userId)
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError("Invalid user id", configs.NotFound, entities.E{})
	}

	user, err := s.repo.GetUserById(ui)
	if err != nil {
		bugsnag.Notify(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.FormatApiError("User not found", configs.NotFound, entities.E{})
		}
		return nil, utils.FormatApiError("Failed to process error: "+err.Error(), configs.ServerError, entities.E{})
	}
	return user.ToUserResponse(), nil
}

func (s DefaultUserService) AskUserToResetPassword(request *dto.AskUserToResetPasswordReq) *entities.ApiError {
	ui, err := models.DecodeHashId(request.UserId)
	if err != nil || ui == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError("Invalid user id", configs.NotFound, entities.E{})
	}

	user, err := s.repo.FindUserById(ui)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError("Failed to retrieve user", configs.ServerError, entities.E{})
	}

	data := map[string]string{
		"NAME":               user.Name,
		"RESET_PASSWORD_URL": configs.AppConfig.AppUrl + "/#/forgot-password",
	}
	if err = mail.SendEmail(
		dto.MailData{Name: user.Name, Email: user.Email},
		dto.MailData{Name: "no-reply", Email: "no-reply@axisforestry.com"},
		data,
		35,
	); err != nil {
		bugsnag.Notify(err)
		configs.Logger.Warn(err)
	}
	return nil
}

func (s DefaultUserService) AssignUserRole(request *dto.AssignUserRoleRequest) *entities.ApiError {
	userId, err := models.DecodeHashId(request.UserId)
	if err != nil || userId == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError("invalid user id", configs.BadRequest, entities.E{"user_id": "invalid user id"})
	}
	err = s.repo.UpdateUserByField(userId, models.FindByField{Field: "role", Value: request.Role})
	if err != nil {
		bugsnag.Notify(err)
		errMsg := "failed to update user account"
		return utils.FormatApiError(errMsg, configs.ServerError, entities.E{"user_role": errMsg})
	}
	return nil
}

func (s DefaultUserService) CreateUser(userRequest *dto.RegisterRequest) *entities.ApiError {
	user, _ := s.repo.FindUserByEmail(userRequest.Email)
	if user != nil {
		errMsg := "User already exist"
		return utils.FormatApiError(errMsg, configs.Conflict, entities.E{"email": errMsg})
	}

	pwd, err := actions.PasswordBcrypt(userRequest.Password)
	if err != nil {
		bugsnag.Notify(err)
		errMsg := "Failed to encrypt password"
		return utils.FormatApiError(errMsg, configs.ServerError, entities.E{"password": errMsg})
	}

	currency := "CAD"
	if userRequest.CountryCode != "1" {
		currency = "USD"
	}

	newUser := models.User{
		Country:         userRequest.Country,
		DefaultCurrency: currency,
		Email:           userRequest.Email,
		Name:            userRequest.Name,
		Password:        pwd,
		PhoneNumber:     userRequest.PhoneNumber,
	}
	if err = s.repo.CreateUser(newUser); err != nil {
		bugsnag.Notify(err)
		errMsg := "Internal server error, failed to creat user"
		return utils.FormatApiError(errMsg, configs.ServerError, entities.E{"server": errMsg})
	}

	token, err := actions.GenerateToken(newUser.Email)
	if err != nil {
		bugsnag.Notify(err)
		configs.Logger.Warn(err)
	}

	data := map[string]string{
		"NAME":                   newUser.Name,
		"EMAIL_CONFIRMATION_URL": configs.AppConfig.AppUrl + "/#/account/verification/token/" + token,
	}
	if err = mail.SendEmail(
		dto.MailData{Name: newUser.Name, Email: newUser.Email},
		dto.MailData{Name: "no-reply", Email: "no-reply@axisforestry.com"},
		data,
		6,
	); err != nil {
		bugsnag.Notify(err)
		configs.Logger.Warn(err)
	}
	if err = storage.Cache.SetWithTime(token, newUser.Email, 60*time.Minute); err != nil {
		bugsnag.Notify(err)
		configs.Logger.Warn(err)
	}
	return nil
}

func (s DefaultUserService) DeleteUserAddresses(addressId string) *entities.ApiError {
	aid, err := models.DecodeHashId(addressId)
	if err != nil || aid == 0 {
		bugsnag.Notify(err)
		errMsg := "invalid address id"
		return utils.FormatApiError(errMsg, configs.BadRequest, entities.E{"id": errMsg})
	}

	getAddress, err := s.repo.GetUserAddressByField(models.FindByField{Field: "id", Value: aid})
	if err != nil {
		bugsnag.Notify(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errMsg := "user address not found"
			return utils.FormatApiError(errMsg, configs.NotFound, entities.E{"address_id": errMsg})
		}
		errMsg := "failed to retrieve user address"
		return utils.FormatApiError(errMsg, configs.ServerError, entities.E{"address_id": errMsg})
	}

	if err = s.repo.DelUserAddress(getAddress); err != nil {
		bugsnag.Notify(err)
		errMsg := "failed to delete user address"
		return utils.FormatApiError(errMsg, configs.ServerError, entities.E{"address_id": errMsg})
	}

	return nil
}

func (s DefaultUserService) FindUserId(id string) (*models.User, error) {
	decodedId, err := models.DecodeHashId(id)
	if err != nil {
		bugsnag.Notify(err)
		return nil, err
	}
	return s.repo.FindUserById(decodedId)
}

func (s DefaultUserService) ForgotPassword(fp *dto.ForgotPassword) *entities.ApiError {
	user, err := s.repo.FindUserByEmail(fp.Email)
	if err != nil {
		bugsnag.Notify(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errMsg := "User not found, check email address"
			return utils.FormatApiError(errMsg, configs.NotFound, entities.E{"email": errMsg})
		}
		errMsg := "Internal server error: failed to get user"
		return utils.FormatApiError(errMsg, configs.ServerError, entities.E{"email": errMsg})
	}

	token, err := actions.GenerateToken(user.Email)
	if err != nil {
		bugsnag.Notify(err)
		errMsg := "Internal server error: failed to generate token"
		utils.FormatApiError(errMsg, configs.ServerError, entities.E{"email": errMsg})
	}

	err = storage.Cache.SetWithTime(token, user.Email, 60*time.Minute)
	if err != nil {
		bugsnag.Notify(err)
		errMsg := "Internal server error: failed to store token"
		utils.FormatApiError(errMsg, configs.ServerError, entities.E{"email": errMsg})
	}

	data := map[string]string{
		"NAME":               user.Name,
		"PASSWORD_RESET_URL": configs.AppConfig.AppUrl + "/#/auth/password/reset?password_reset_token=" + token,
	}
	if err = mail.SendEmail(
		dto.MailData{Name: user.Name, Email: user.Email},
		dto.MailData{Name: "no-reply", Email: "no-reply@axisforestry.com"},
		data,
		7,
	); err != nil {
		bugsnag.Notify(err)
		errMsg := "Internal server error: failed to send email"
		utils.FormatApiError(errMsg, configs.ServerError, entities.E{"email": errMsg})
	}
	return nil
}

func (s DefaultUserService) GetAddresses(userHid string) ([]dto.UserAddressResponse, *entities.ApiError) {
	userId, err := models.DecodeHashId(userHid)
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError("invalid user id", configs.BadRequest, entities.E{"user_id": "invalid user id"})
	}

	getAddresses, err := s.repo.GetUserAddresses(userId)
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError("Internal server error", configs.ServerError, entities.E{"error": err.Error()})
	}

	addresses := []dto.UserAddressResponse{}
	for _, address := range getAddresses {
		addresses = append(addresses, address.ToResponse())
	}
	return addresses, nil
}

func (s DefaultUserService) GetUser(userId string) (*dto.UserResponse, *entities.ApiError) {
	user, err := storage.Cache.Remember(
		userId,
		func(f interface{}) (string, error) {
			ui, err := models.DecodeHashId(userId)
			if err != nil {
				bugsnag.Notify(err)
				return "", err
			}

			user, err := s.repo.FindUserById(ui)
			if err != nil {
				bugsnag.Notify(err)
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return "", gorm.ErrRecordNotFound
				}
				return "", err
			}

			byteUser, err := json.Marshal(user)
			if err != nil {
				bugsnag.Notify(err)
				return "", err
			}
			return string(byteUser), nil
		},
		userId,
	)
	if err != nil {
		bugsnag.Notify(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.FormatApiError("User not found", configs.NotFound, entities.E{})
		}
		return nil, utils.FormatApiError("Failed to process error"+err.Error(), configs.ServerError, entities.E{})
	}

	u := &models.User{}
	err = json.Unmarshal([]byte(user), u)
	if err != nil {
		bugsnag.Notify(err)
		return nil, utils.FormatApiError("Failed to process error"+err.Error(), configs.ServerError, entities.E{})
	}
	return u.ToUserResponse(), nil
}

func (s DefaultUserService) GetUserByEmail(email string) (*dto.UserResponse, *entities.ApiError) {
	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		bugsnag.Notify(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.FormatApiError("User not found", configs.NotFound, entities.E{})
		}
		return nil, utils.FormatApiError("Failed to process error"+err.Error(), configs.ServerError, entities.E{})
	}
	return user.ToUserResponse(), nil
}

func (s DefaultUserService) GetUserToken(user *models.User) (*actions.TokenDetail, error) {
	hashedUserID, _ := models.EncodeHashId(user.ID)
	token, err := actions.IssueToken(&dto.UserJwtEntity{UserId: hashedUserID, Roles: []string{user.Role}})
	if err != nil {
		bugsnag.Notify(err)
		return nil, err
	}

	// store to redis
	if err = s.repo.CreateAuth(user, token); err != nil {
		bugsnag.Notify(err)
		return nil, err
	}
	return token, nil
}

func (s DefaultUserService) GetUsers(limit, offset int) ([]dto.UserResponse, *entities.ApiError) {
	users, err := s.repo.GetUsers(limit, offset)
	if err != nil {
		bugsnag.Notify(err)
		errMsg := "failed to retrieve users"
		return nil, utils.FormatApiError(errMsg, configs.ServerError, entities.E{"users": errMsg})
	}

	var usersResp []dto.UserResponse
	for _, user := range users {
		u := user.ToUserResponse()
		usersResp = append(usersResp, *u)
	}
	return usersResp, nil
}

func (s DefaultUserService) ResendEmailVerificationToken(email string) *entities.ApiError {
	user, apiErr := s.GetUserByEmail(email)
	if apiErr != nil {
		return apiErr
	}
	if user.Verified {
		return utils.FormatApiError("user already verified", configs.BadRequest, entities.E{})
	}

	token, err := actions.GenerateToken(user.Email)
	if err != nil {
		bugsnag.Notify(err)
		configs.Logger.Warn(err)
	}

	data := map[string]string{
		"NAME":                   user.Name,
		"EMAIL_CONFIRMATION_URL": configs.AppConfig.AppUrl + "/#/account/verification/token/" + token,
	}
	if err = mail.SendEmail(
		dto.MailData{Name: user.Name, Email: user.Email},
		dto.MailData{Name: "no-reply", Email: "no-reply@axisforestry.com"},
		data,
		6,
	); err != nil {
		bugsnag.Notify(err)
		configs.Logger.Warn(err)
	}

	if err = storage.Cache.SetWithTime(token, user.Email, 60*time.Minute); err != nil {
		bugsnag.Notify(err)
		configs.Logger.Warn(err)
	}
	return nil
}

func (s DefaultUserService) SearchUsers(limit, offset int, search string) ([]dto.UserResponse, *entities.ApiError) {
	users, err := s.repo.SearchUsers(limit, offset, search)
	if err != nil {
		bugsnag.Notify(err)
		errMsg := "failed to retrieve users"
		return nil, utils.FormatApiError(errMsg, configs.ServerError, entities.E{"users": errMsg})
	}

	var usersResp []dto.UserResponse
	for _, user := range users {
		u := user.ToUserResponse()
		usersResp = append(usersResp, *u)
	}
	return usersResp, nil
}

func (s DefaultUserService) SetDefaultAddresses(addressId string) *entities.ApiError {
	aid, err := models.DecodeHashId(addressId)
	if err != nil || aid == 0 {
		bugsnag.Notify(err)
		errMsg := "invalid address id"
		return utils.FormatApiError(errMsg, configs.BadRequest, entities.E{"id": errMsg})
	}

	getAddress, err := s.repo.GetUserAddressByField(models.FindByField{Field: "id", Value: aid})
	if err != nil {
		bugsnag.Notify(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errMsg := "user address not found"
			return utils.FormatApiError(errMsg, configs.NotFound, entities.E{"address_id": errMsg})
		}

		errMsg := "failed to retrieve user address"
		return utils.FormatApiError(errMsg, configs.ServerError, entities.E{"address_id": errMsg})
	}

	getAddress.IsDefaultAddress = true
	if err = s.repo.UpdateUserAddress(getAddress); err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError("Failed to process error "+err.Error(), configs.ServerError, entities.E{})
	}
	return nil
}

func (s DefaultUserService) UpdateAddress(
	userHid, addressId string,
	request *dto.UserAddressRequest,
) *entities.ApiError {
	userId, err := models.DecodeHashId(userHid)
	if err != nil || userId == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError("invalid user id", configs.BadRequest, entities.E{"user_id": "invalid user id"})
	}

	aid, err := models.DecodeHashId(addressId)
	if err != nil || aid == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError("invalid address id", configs.BadRequest, entities.E{"id": "invalid address id"})
	}

	user, err := s.repo.FindUserById(userId)
	if err != nil {
		bugsnag.Notify(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.FormatApiError("User not found", configs.NotFound, entities.E{})
		}
		return utils.FormatApiError("Failed to process error"+err.Error(), configs.ServerError, entities.E{})
	}

	getAddress, err := s.repo.GetUserAddressByField(models.FindByField{Field: "id", Value: aid})
	if err != nil {
		bugsnag.Notify(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errMsg := "user address not found"
			return utils.FormatApiError(errMsg, configs.NotFound, entities.E{"address_id": errMsg})
		}

		errMsg := "failed to retrieve user address"
		return utils.FormatApiError(errMsg, configs.ServerError, entities.E{"address_id": errMsg})
	}

	a := getAddress.Address
	a.City = request.City
	a.StreetName = request.StreetName
	a.Country = request.Country
	a.Province = request.Province
	a.State = request.State
	a.ZipCode = request.ZipCode

	c := getAddress.Contact
	c.Name = request.Name
	c.Email = request.Email
	c.Phone = request.Phone
	c.Organisation = request.Organisation

	getAddress.User = *user
	getAddress.Address = a
	getAddress.Contact = c
	getAddress.IsDefaultAddress = true
	if err = s.repo.UpdateUserAddress(getAddress); err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError("Failed to process error "+err.Error(), configs.ServerError, entities.E{})
	}
	return nil
}

func (s DefaultUserService) UpdateUserCurrency(request *dto.UserUpdateCurrencyRequest) *entities.ApiError {
	userId, err := models.DecodeHashId(request.UserId)
	if err != nil || userId == 0 {
		bugsnag.Notify(err)
		return utils.FormatApiError("invalid user id", configs.BadRequest, entities.E{"user_id": "invalid user id"})
	}

	if err := s.repo.UpdateUserByField(
		userId,
		models.FindByField{Field: "default_currency", Value: request.Currency},
	); err != nil {
		bugsnag.Notify(err)
		errMsg := "failed to update user currency"
		return utils.FormatApiError(errMsg, configs.ServerError, entities.E{"allow_on_account": errMsg})
	}
	return nil
}

func (s DefaultUserService) UpdateUser(userId string, updates []entities.UpdateFields) *entities.ApiError {
	user, apiErr := s.GetUser(userId)
	if apiErr != nil {
		return apiErr
	}

	ui, err := models.DecodeHashId(user.Id)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError("Invalid user id", configs.NotFound, entities.E{})
	}

	u := &models.User{}
	u.ID = ui
	data := make(map[string]interface{})
	for _, update := range updates {
		data[update.Field] = update.Value
	}

	err = s.repo.UpdateUser(u, data)
	if err != nil {
		bugsnag.Notify(err)
		return utils.FormatApiError("Failed to update user", configs.ServerError, entities.E{})
	}

	_, err = storage.Cache.Del(userId)
	if err != nil {
		bugsnag.Notify(err)
		configs.Logger.Warn("failed to delete cache key")
	}
	return nil
}

func (s DefaultUserService) UpdateUserOnAccount(request *dto.UserUpdateOnAccountRequest) *entities.ApiError {
	userId, err := models.DecodeHashId(request.UserId)
	if err != nil || userId == 0 {
		bugsnag.Notify(err)
		errMsg := "invalid user id"
		return utils.FormatApiError(errMsg, configs.BadRequest, entities.E{"user_id": errMsg})
	}

	if err = s.repo.UpdateUserByField(
		userId,
		models.FindByField{Field: "allow_on_account", Value: request.Allow},
	); err != nil {
		bugsnag.Notify(err)
		errMsg := "failed to update user allow on account"
		return utils.FormatApiError(errMsg, configs.ServerError, entities.E{"allow_on_account": errMsg})
	}
	return nil
}

func (s DefaultUserService) UserLogin(ur *dto.LoginRequest) (*actions.TokenDetail, *entities.ApiError) {
	user, err := s.repo.FindUserByEmail(ur.Email)
	if err != nil {
		bugsnag.Notify(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errMsg := "User not found, check email address"
			return nil, utils.FormatApiError(errMsg, configs.NotFound, entities.E{"email": errMsg})
		}
		errMsg := "Failed to retrieve user"
		return nil, utils.FormatApiError(errMsg, configs.ServerError, entities.E{"email": errMsg})
	}

	if !user.EmailVerifiedAt.Valid {
		errMsg := "user email not verified"
		return nil, utils.FormatApiError(errMsg, configs.Unauthorized, entities.E{"email": errMsg})
	}

	ok := user.VerifyPassword(ur.Password)
	if !ok {
		errMsg := "Incorrect password, check your credential"
		return nil, utils.FormatApiError(errMsg, configs.NotFound, entities.E{"password": errMsg})
	}

	hid, _ := models.EncodeHashId(user.ID)
	storage.Cache.Del(hid)
	token, err := s.GetUserToken(user)
	if err != nil {
		bugsnag.Notify(err)
		configs.Logger.Error(err)
		errMsg := "Interanl server error, failed to create auth token"
		return nil, utils.FormatApiError(errMsg, configs.ServerError, entities.E{"server": errMsg})
	}
	return token, nil
}
