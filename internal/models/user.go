package models

import (
	"database/sql"

	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/actions"
	"axis/ecommerce-backend/internal/dto"
)

type User struct {
	configs.GormModel
	Address         []UserAddress
	AllowOnAccount  bool
	BusinessName    string
	Country         string
	DefaultCurrency string `gorm:"default:CAD"`
	Email           string
	EmailVerifiedAt sql.NullTime `gorm:"type:TIMESTAMP NULL"`
	IsActive        bool
	Name            string
	Password        string
	PhoneNumber     string
	Role            string `gorm:"default:Member"`
}

type UserAddress struct {
	configs.GormModel
	Address          Address
	AddressId        uint
	Contact          Contact
	ContactId        uint
	IsDefaultAddress bool
	Type             string
	User             User
	UserId           uint
}

func (userAddress UserAddress) ToResponse() dto.UserAddressResponse {
	id, _ := EncodeHashId(userAddress.ID)
	userId, _ := EncodeHashId(userAddress.UserId)

	return dto.UserAddressResponse{
		Address:          userAddress.Address.ToResponse(),
		Contact:          userAddress.Contact.ToResponse(),
		Id:               id,
		IsDefaultAddress: userAddress.IsDefaultAddress,
		Type:             userAddress.Type,
		UserId:           userId,
	}
}

func (u *User) VerifyPassword(requestPwd string) bool {
	if len(u.Password) == 0 {
		return false
	}
	verified, _ := actions.VerifyPassword(requestPwd, u.Password)
	return verified
}

func (u *User) ToUserResponse() *dto.UserResponse {
	userId, _ := EncodeHashId(u.ID)
	return &dto.UserResponse{
		AllowOnAccount: u.AllowOnAccount,
		Currency:       u.DefaultCurrency,
		Email:          u.Email,
		Id:             userId,
		Name:           u.Name,
		Password:       u.Password,
		PhoneNumber:    u.PhoneNumber,
		Role:           u.Role,
		Verified:       u.EmailVerifiedAt.Valid,
	}
}
