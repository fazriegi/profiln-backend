package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	db "profiln-be/db/sqlc"
	"profiln-be/libs"
	email "profiln-be/libs/email"
	"profiln-be/model"
	repository "profiln-be/package/auth/repository"

	"github.com/sirupsen/logrus"
)

type IAuthUsecase interface {
	Login(loginType string, props *model.LoginRequest) (resp model.Response)
	ResetPassword(props *model.ResetPasswordRequest) (resp model.Response)
	Register(props *model.RegisterRequest, oauth string) (resp model.Response)
	UpdateVerifiedEmail(props *model.VerifiedEmailOTPRequest) (resp model.Response)
	SendResetPasswordEmail(props *model.ResetPasswordEmailRequest) (resp model.Response)
	GetUserOtpByEmail(email string) (resp model.Response)
	SendOTPEmail(props *model.OTPEmailRequest) (resp model.Response)
}

type AuthUsecase struct {
	repository repository.IAuthRepository
	email      email.IEmail
	log        *logrus.Logger
}

func NewAuthUsecase(repository repository.IAuthRepository, email email.IEmail, log *logrus.Logger) IAuthUsecase {
	return &AuthUsecase{
		repository,
		email,
		log,
	}
}

func (u *AuthUsecase) Login(loginType string, props *model.LoginRequest) (resp model.Response) {
	user, err := u.repository.GetUserByEmail(props.Email)

	if err != nil && err == sql.ErrNoRows {
		resp.Status = libs.CustomResponse(http.StatusUnauthorized, "Incorrect email or password!")

		return
	} else if err != nil {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")

		u.log.Errorf("repository.GetUserByEmail: %v", err)
		return
	}

	if loginType == "app" {
		if !libs.CheckPasswordHash(props.Password, user.Password.String) {
			resp.Status = libs.CustomResponse(http.StatusUnauthorized, "Incorrect email or password!")

			return
		}
	}

	token, err := libs.GenerateJWTToken(user.ID, user.Email, time.Hour*24)
	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")

		u.log.Errorf("libs.GenerateJWTToken: %v", err)
		return
	}

	resp.Status = libs.CustomResponse(http.StatusOK, "Success login")
	resp.Data = map[string]string{
		"token": token,
	}
	return
}

func (u *AuthUsecase) ResetPassword(props *model.ResetPasswordRequest) (resp model.Response) {
	user, err := u.repository.GetUserByEmail(props.Email)
	if err != nil && err == sql.ErrNoRows {
		resp.Status = libs.CustomResponse(http.StatusNotFound, "Email not found")

		return
	} else if err != nil {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")

		u.log.Errorf("repository.GetUserByEmail: %v", err)
		return
	}

	hashedPassword, err := libs.HashPassword(props.Password)
	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")

		u.log.Errorf("libs.HashPassword: %v", err)
		return
	}

	err = u.repository.UpdateUserPassword(user.ID, hashedPassword)
	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")

		u.log.Errorf("repository.UpdateUserPassword: %v", err)
		return
	}

	resp.Status = libs.CustomResponse(http.StatusOK, "Success reset password")
	return
}

func (u *AuthUsecase) Register(props *model.RegisterRequest, oauth string) (resp model.Response) {
	subject := "OTP Regristation Profiln"
	user, _ := u.repository.GetUserByEmail(props.Email)

	if user.Email != "" {
		resp.Status = libs.CustomResponse(http.StatusUnauthorized, "Email already been taken")
		return resp
	}

	hashedPassword, err := libs.HashPassword(props.Password)
	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Something went wrong")
		u.log.Errorf("libs.HashPassword : %v", err)
		return resp
	}

	var verifiedEmail bool

	if oauth == "true" {
		verifiedEmail = true
	} else {
		verifiedEmail = false
	}

	registerUserParams := db.InsertUserParams{
		Email:         props.Email,
		Password:      sql.NullString{String: hashedPassword, Valid: true},
		FullName:      props.Fullname,
		VerifiedEmail: sql.NullBool{Bool: verifiedEmail, Valid: true},
	}

	insertUser, err := u.repository.CreateUser(registerUserParams)

	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Something went wrong while register")
		u.log.Errorf("repository.CreateUser : %v", err)
		return resp
	}

	if oauth != "true" {
		otp, err := libs.GenerateOTP(6)
		if err != nil {
			resp.Status = libs.CustomResponse(http.StatusBadRequest, "Failed to generate otp")
			u.log.Errorf("libs.GenerateOTP : %v", err)
			return resp
		}

		_, err = u.repository.InsertOtp(insertUser.ID, otp)
		if err != nil {
			resp.Status = libs.CustomResponse(http.StatusBadRequest, "Failed to insert otp")
			u.log.Errorf("repository.InsertOtp: %v", err)
			return resp
		}

		DigitOne := string([]rune(otp)[0])
		DigitTwo := string([]rune(otp)[1])
		DigitThree := string([]rune(otp)[2])
		DigitFour := string([]rune(otp)[3])
		DigitFive := string([]rune(otp)[4])
		DigitSix := string([]rune(otp)[5])

		dataEmail := model.OTPEmail{
			Email:      props.Email,
			DigitOne:   DigitOne,
			DigitTwo:   DigitTwo,
			DigitThree: DigitThree,
			DigitFour:  DigitFour,
			DigitFive:  DigitFive,
			DigitSix:   DigitSix,
		}

		go func() {
			err := u.email.SendAuthEmail(subject, []string{props.Email}, dataEmail, "otp.html")
			if err != nil {
				u.log.Errorf("email.SendAuthEmail: %v", err)
			}
		}()

	}

	resp.Status = libs.CustomResponse(http.StatusCreated, "Success to register")
	userResponse := model.RegisterResponse{
		ID:            int(insertUser.ID),
		Email:         insertUser.Email,
		Fullname:      insertUser.FullName,
		VerifiedEmail: insertUser.VerifiedEmail.Bool,
	}

	resp.Data = userResponse

	return resp
}

func (u *AuthUsecase) UpdateVerifiedEmail(props *model.VerifiedEmailOTPRequest) (resp model.Response) {
	_, err := u.repository.GetUserOtpByOtp(props.Otp)
	if err != nil && err == sql.ErrNoRows {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "OTP doesnt exist")

		return
	} else if err != nil {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")
		u.log.Errorf("repository.GetUserOtpByOtp %v", err)
		return
	}

	user, err := u.repository.UpdateVerifiedEmail(props.Otp, props.Email)

	if err != nil && err == sql.ErrNoRows {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Failed verify email")

		return
	} else if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Failed verify email")
		u.log.Errorf("repository.UpdateVerifiedEmailByOTP %v", err)
		return
	}

	err = u.repository.DeleteOtp(props.Otp)

	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Failed verify email")
		u.log.Errorf("repository.DeleteOtp %v", err)
		return
	}

	token, err := libs.GenerateJWTToken(user.ID, user.Email, time.Hour*24)
	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")

		u.log.Errorf("libs.GenerateJWTToken: %v", err)
		return
	}

	resp.Status = libs.CustomResponse(http.StatusOK, "Success verify email")
	resp.Data = map[string]string{
		"token": token,
	}
	return
}

func (u *AuthUsecase) SendResetPasswordEmail(props *model.ResetPasswordEmailRequest) (resp model.Response) {
	resp.Status = libs.CustomResponse(http.StatusOK, "Success")
	user, err := u.repository.GetUserByEmail(props.Email)

	if err != nil && err == sql.ErrNoRows {
		resp.Status = libs.CustomResponse(http.StatusNotFound, "Email not found")

		return
	} else if err != nil {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")

		u.log.Errorf("repository.GetUserByEmail: %v", err)
		return
	}

	subject := "Permintaan Reset Password"
	resetPasswordUrl := os.Getenv("FRONTEND_RESET_PASSWORD_URL")
	jwtToken, err := libs.GenerateJWTToken(user.ID, props.Email, time.Minute*30)

	if err != nil {
		u.log.Errorf("libs.GenerateJWTToken: %v", err)
		return
	}

	redirectLink := fmt.Sprintf("%s?token=%s", resetPasswordUrl, jwtToken)

	// data for template html
	data := model.ResetPasswordEmail{
		Email: props.Email,
		URL:   redirectLink,
	}

	// send email asynchronously.
	// no matter it's success or not,
	// it always returns success to the client
	go func() {
		err := u.email.SendAuthEmail(subject, []string{props.Email}, data, "reset-password.html")
		if err != nil {
			u.log.Errorf("email.SendAuthEmail: %v", err)
		}
	}()

	return
}

func (u *AuthUsecase) GetUserOtpByEmail(email string) (resp model.Response) {
	_, err := u.repository.GetUserByEmail(email)

	if err != nil && err == sql.ErrNoRows {
		resp.Status = libs.CustomResponse(http.StatusNotFound, "Email not found")

		return
	} else if err != nil {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")

		u.log.Errorf("repository.GetUserByEmail: %v", err)
		return
	}

	userOtpByEmail, err := u.repository.GetUserOtpByEmail(email)

	if err != nil && err == sql.ErrNoRows {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "OTP doesnt exist")

		return
	} else if err != nil {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")
		u.log.Errorf("repository.GetUserOtpByOtp %v", err)
		return
	}

	userOtpByEmailResp := model.UserOTPByEmailResponse{
		ID:            int(userOtpByEmail.ID),
		Email:         userOtpByEmail.Email,
		Fullname:      userOtpByEmail.FullName,
		VerifiedEmail: userOtpByEmail.VerifiedEmail.Bool,
		UserID:        int(userOtpByEmail.ID_2),
		Otp:           userOtpByEmail.Otp.String,
	}

	resp.Status = libs.CustomResponse(http.StatusOK, "Success get user otp by email")
	resp.Data = userOtpByEmailResp
	return
}

func (u *AuthUsecase) SendOTPEmail(props *model.OTPEmailRequest) (resp model.Response) {
	resp.Status = libs.CustomResponse(http.StatusOK, "Success send otp")
	subject := "OTP Regristation Profiln"
	_, err := u.repository.GetUserByEmail(props.Email)

	if err != nil && err == sql.ErrNoRows {
		resp.Status = libs.CustomResponse(http.StatusNotFound, "Email not found")

		return
	} else if err != nil {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")

		u.log.Errorf("repository.GetUserByEmail: %v", err)
		return
	}

	userOtpByEmail, err := u.repository.GetUserOtpByEmail(props.Email)

	if err != nil && err == sql.ErrNoRows {
		resp.Status = libs.CustomResponse(http.StatusNotFound, "OTP doesnt exist")

		return
	} else if err != nil {
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured")
		u.log.Errorf("repository.GetUserOtpByEmail %v", err)
		return
	}

	DigitOne := string([]rune(userOtpByEmail.Otp.String)[0])
	DigitTwo := string([]rune(userOtpByEmail.Otp.String)[1])
	DigitThree := string([]rune(userOtpByEmail.Otp.String)[2])
	DigitFour := string([]rune(userOtpByEmail.Otp.String)[3])
	DigitFive := string([]rune(userOtpByEmail.Otp.String)[4])
	DigitSix := string([]rune(userOtpByEmail.Otp.String)[5])

	dataEmail := model.OTPEmail{
		Email:      props.Email,
		DigitOne:   DigitOne,
		DigitTwo:   DigitTwo,
		DigitThree: DigitThree,
		DigitFour:  DigitFour,
		DigitFive:  DigitFive,
		DigitSix:   DigitSix,
	}

	go func() {
		err := u.email.SendAuthEmail(subject, []string{props.Email}, dataEmail, "otp.html")
		if err != nil {
			u.log.Errorf("email.SendAuthEmail: %v", err)
		}
	}()

	return
}
