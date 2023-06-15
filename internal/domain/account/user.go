package account

import (
	"axis/ecommerce-backend/internal/actions"
	"axis/ecommerce-backend/internal/models"
)

type UserRepo interface {
	AddUserAddress(data models.UserAddress) (*models.UserAddress, error)
	CreateAuth(u *models.User, token *actions.TokenDetail) error
	CreateUser(u models.User) error
	DelUserAddress(userAddress *models.UserAddress) error
	FindUserByEmail(email string) (*models.User, error)
	FindUserById(id uint) (*models.User, error)
	GetUserAddresses(userId uint) ([]models.UserAddress, error)
	GetUserAddressByID(addressID uint) (*models.UserAddress, error)
	GetUserAddressByField(fv models.FindByField) (*models.UserAddress, error)
	GetUserById(id uint) (*models.User, error)
	GetUsers(limit, offset int) ([]models.User, error)
	SearchUsers(limit, offset int, search string) ([]models.User, error)
	SetUserDefaultAddress(data *models.UserAddress) error
	UpdateUser(u *models.User, data map[string]interface{}) error
	UpdateUserAddress(userAddress *models.UserAddress) error
	UpdateUserByField(userId uint, data models.FindByField) error
}
