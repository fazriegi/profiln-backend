package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"profiln-be/libs"
	"profiln-be/model"
	"profiln-be/repository"
)

type EmailUsecase struct {
	db         *sql.DB
	repository *repository.Queries
}

func NewEmailUsecase(db *sql.DB) *EmailUsecase {
	return &EmailUsecase{
		db:         db,
		repository: repository.New(db),
	}
}

func (u *EmailUsecase) SendResetPasswordMail(props *model.SendResetPassEmailRequest) (resp model.Response) {
	_, err := u.repository.GetUserByEmail(context.Background(), props.Email)

	if err != nil && err == sql.ErrNoRows {
		resp.Status = libs.CustomResponse(http.StatusNotFound, "Email not found")

		return resp
	} else if err != nil {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")

		log.Printf("repository.GetUserByEmail: %v", err)
		return resp
	}

	resp.Status = libs.CustomResponse(http.StatusOK, "success")
	subject := "Permintaan Reset Password"
	resetPasswordUrl := os.Getenv("FRONTEND_RESET_PASSWORD_URL")
	jwtToken, err := libs.GenerateJWTToken(props.Email)

	if err != nil {
		log.Printf("libs.GenerateJWTToken: %v", err)
		return
	}

	redirectLink := fmt.Sprintf("%s?token=%s", resetPasswordUrl, jwtToken)

	// data for template html
	data := struct {
		Email string
		URL   string
	}{
		Email: props.Email,
		URL:   redirectLink,
	}

	// get working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Printf("os.Getwd: %v", err)
		return
	}

	filepath := fmt.Sprintf("%s/template/%s", dir, "reset-password.html")

	body, err := libs.HTMLToString(filepath, data)
	if err != nil {
		log.Printf("libs.HTMLToString: %v", err)
		return
	}

	// send email asynchronously.
	// no matter it's success or not,
	// it always returns success to the client
	go func() {
		err := libs.SendEmail(subject, []string{props.Email}, body)
		if err != nil {
			log.Printf("libs.SendMail: %v", err)
		}
	}()

	return resp
}
