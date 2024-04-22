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
