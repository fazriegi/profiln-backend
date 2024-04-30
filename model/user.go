package model

type UserLoginRequest struct {
	Email    string `validate:"required,email"`
	Password string
}

type UserResetPasswordRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,password"`
}

type UserRegisterRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
	Fullname string `validate:"required"`
}

type RegisterResponse struct {
	ID            int    `json:"id"`
	Email         string `json:"email"`
	Fullname      string `json:"fullname"`
	VerifiedEmail bool   `json:"verified_email"`
}

type VerifiedEmailOTPRequest struct {
	Otp string `validate:"required"`
}
