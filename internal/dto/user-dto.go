package dto

import (
	"time"
)

type AskUserToResetPasswordReq struct {
	UserId string `json:"user_id" binding:"required"`
}

type AssignUserRoleRequest struct {
	Role   string `json:"role" binding:"required,oneof='Admin' 'Member' 'Distributor' 'Staff'"`
	UserId string `json:"user_id" binding:"required"`
}

type EmbedUser struct {
	Id string `json:"id"`

	Email       string `json:"email"`
	IsActive    bool   `json:"is_active"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Verified    bool   `json:"verified"`
}

type EmailVerificationTokenRequest struct {
	Email string `json:"email" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UserAddressRequest struct {
	City         string `json:"city"`
	Country      string `json:"country" binding:"required"`
	Email        string `json:"email" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Organisation string `json:"organisation" binding:"required"`
	Phone        string `json:"phone" binding:"required"`
	Province     string `json:"province"`
	State        string `json:"state"`
	StreetName   string `json:"street_name"`
	ZipCode      string `json:"zip_code"`
}

type UserAddressResponse struct {
	Id string `json:"id"`

	Address          AddressResponse `json:"address"`
	Contact          ContactResponse `json:"contact"`
	IsDefaultAddress bool            `json:"is_default_address"`
	Type             string          `json:"type"`
	UserId           string          `json:"user_id"`
}

type UserResponse struct {
	Id string `json:"id"`

	AllowOnAccount bool   `json:"allow_on_account"`
	Currency       string `json:"currency"`
	Email          string `json:"email"`
	IsActive       bool   `json:"is_active"`
	Name           string `json:"name"`
	Password       string `json:"password"`
	PhoneNumber    string `json:"phone_number"`
	Role           string `json:"role"`
	Verified       bool   `json:"verified"`
}

type UserSession struct {
	ID string

	AccessSessionId  string
	CreatedAt        time.Time
	Email            string
	EmailVerifiedAt  time.Time
	IsActive         bool
	Name             string
	PhoneNumber      string
	RefreshSessionId string
	Role             string
	UpdatedAt        time.Time
}

type UserUpdateOnAccountRequest struct {
	Allow  bool   `json:"allow"`
	UserId string `json:"user_id" binding:"required"`
}

type UserUpdateCurrencyRequest struct {
	Currency string `json:"currency" binding:"required,oneof='USD' 'CAD'"`
	UserId   string `json:"user_id" binding:"required"`
}

type UserUpdateTaxExemptionRequest struct {
	UserId string `json:"user_id" binding:"required"`
}

type UserUpdateRequest struct {
	BusinessName string `json:"business_name"`
	Email        string `json:"email"`
	Name         string `json:"name" binding:"required"`
	PhoneNumber  string `json:"phone_number" binding:"required"`
}
