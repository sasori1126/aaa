package dto

type RegisterRequest struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	PhoneNumber     string `json:"phone_number" binding:"required"`
	CountryCode     string `json:"country_code" binding:"required,numeric"`
	Country         string `json:"country"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type ForgotPassword struct {
	Email string `json:"email" binding:"required"`
}

type ResetPasswordRequest struct {
	Token           string `json:"token" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type AccountResetPasswordRequest struct {
	UserId          string `json:"user_id" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type UserJwtEntity struct {
	UserId string
	Asid   string
	Rsid   string
	Roles  []string
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AtExpires    int64  `json:"access_token_expires_in"`
	RtExpires    int64  `json:"refresh_token_expires_in"`
}
