package usecase

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"profiln-be/libs"
	"profiln-be/model"
)

type EmailUsecase struct{}

func NewEmailUsecase() *EmailUsecase {
	return &EmailUsecase{}
}

func (u *EmailUsecase) SendResetPasswordMail(props *model.SendResetPassEmailRequest) (resp model.Response) {
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
		err := libs.SendMail(subject, []string{props.Email}, body)
		if err != nil {
			log.Printf("libs.SendMail: %v", err)
		}
	}()

	return resp
}
