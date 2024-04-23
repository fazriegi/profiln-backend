package usecase

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"profiln-be/libs"
	"profiln-be/model"
	"profiln-be/repository"
)

type UserUsecase struct {
	db         *sql.DB
	repository *repository.Queries
}

func NewUserUsecase(db *sql.DB) *UserUsecase {
	return &UserUsecase{
		db:         db,
		repository: repository.New(db),
	}
}

func (u *UserUsecase) Login(props *model.UserLoginRequest) (resp model.Response) {

	user, err := u.repository.GetUserByEmail(context.Background(), props.Email)

	if err != nil && err == sql.ErrNoRows {
		resp.Status = libs.CustomResponse(http.StatusUnauthorized, "Incorrect email or password!")

		return resp
	} else if err != nil {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")

		log.Printf("repository.GetUserByEmail: %v", err)
		return resp
	}

	if !libs.CheckPasswordHash(props.Password, user.Password.String) {
		resp.Status = libs.CustomResponse(http.StatusUnauthorized, "Incorrect email or password!")

		return resp
	}

	token, err := libs.GenerateJWTToken(int(user.ID), user.Email)
	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")

		log.Printf("libs.GenerateJWTToken: %v", err)
		return resp
	}

	resp.Status = libs.CustomResponse(http.StatusOK, "success login")
	resp.Data = map[string]string{
		"token": token,
	}

	return resp
}

func (u *UserUsecase) ResetPassword(props *model.UserResetPasswordRequest) (resp model.Response) {
	user, err := u.repository.GetUserByEmail(context.Background(), props.Email)
	if err != nil && err == sql.ErrNoRows {
		resp.Status = libs.CustomResponse(http.StatusOK, "Success")

		return
	} else if err != nil {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")

		log.Printf("repository.GetUserByEmail: %v", err)
		return
	}

	hashedPassword, err := libs.HashPassword(props.Password)
	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")

		log.Printf("libs.HashPassword: %v", err)
		return
	}

	arg := repository.UpdateUserPasswordParams{
		ID:       user.ID,
		Password: sql.NullString{String: hashedPassword, Valid: true},
	}

	err = u.repository.UpdateUserPassword(context.Background(), arg)
	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")

		log.Printf("repository.UpdateUserPassword: %v", err)
		return
	}

	resp.Status = libs.CustomResponse(http.StatusOK, "Success reset password")

	return
}
