package model

type UserLoginRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}

type UserResetPasswordRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,password"`
}
