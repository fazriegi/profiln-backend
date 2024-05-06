package model

type LoginRequest struct {
	Email    string `validate:"required,email"`
	Password string
}

type ResetPasswordRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,password"`
}

type ResetPasswordEmailRequest struct {
	Email string `validate:"required,email"`
}

type RegisterRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,password"`
	Fullname string `validate:"required"`
}

type RegisterResponse struct {
	ID            int    `json:"id"`
	Email         string `json:"email"`
	Fullname      string `json:"fullname"`
	VerifiedEmail bool   `json:"verified_email"`
}

type VerifiedEmailOTPRequest struct {
	Otp   string `validate:"required"`
	Email string `validate:"required,email"`
}
