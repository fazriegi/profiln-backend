package model

type SendResetPassEmailRequest struct {
	Email string `validate:"required,email"`
}
