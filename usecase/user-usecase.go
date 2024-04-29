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

func (u *UserUsecase) Register(props *model.UserRegisterRequest, oauth string) (resp model.Response) {
	subject := "OTP Regristation Profiln"
	user, _ := u.repository.GetUserByEmail(context.Background(), props.Email)

	if user.Email != "" {
		resp.Status = libs.CustomResponse(http.StatusUnauthorized, "Email already been taken")
		return resp
	}

	hashedPassword, err := libs.HashPassword(props.Password)
	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Something went wrong")
		log.Printf("libs.HashPassword : %v", err)
		return resp
	}

	var verifiedEmail bool

	if oauth == "true" {
		verifiedEmail = true
	} else {
		verifiedEmail = false
	}

	registerUserParams := repository.InsertUserParams{
		Email:         props.Email,
		Password:      sql.NullString{String: hashedPassword, Valid: true},
		FullName:      props.Fullname,
		VerifiedEmail: sql.NullBool{Bool: verifiedEmail, Valid: true},
	}

	insertUser, err := u.repository.InsertUser(context.Background(), registerUserParams)

	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Something went wrong while register")
		log.Printf("Register : %v", err)
		return resp
	}

	if oauth != "true" {
		otp, err := libs.GenerateOTP(6)
		if err != nil {
			resp.Status = libs.CustomResponse(http.StatusBadRequest, "Failed to generate otp")
			log.Printf("libs.GenerateOTP : %v", err)
			return resp
		}

		insertOtpParams := repository.InsertOtpParams{
			UserID: sql.NullInt64{Int64: insertUser.ID, Valid: true},
			Otp:    sql.NullString{String: otp, Valid: true},
		}

		_, err = u.repository.InsertOtp(context.Background(), insertOtpParams)
		if err != nil {
			resp.Status = libs.CustomResponse(http.StatusBadRequest, "Failed to insert otp")
			return resp
		}

		DigitOne := string([]rune(otp)[0])
		DigitTwo := string([]rune(otp)[1])
		DigitThree := string([]rune(otp)[2])
		DigitFour := string([]rune(otp)[3])
		DigitFive := string([]rune(otp)[4])
		DigitSix := string([]rune(otp)[5])

		dataEmail := struct {
			Email      string
			DigitOne   string
			DigitTwo   string
			DigitThree string
			DigitFour  string
			DigitFive  string
			DigitSix   string
		}{
			Email:      props.Email,
			DigitOne:   DigitOne,
			DigitTwo:   DigitTwo,
			DigitThree: DigitThree,
			DigitFour:  DigitFour,
			DigitFive:  DigitFive,
			DigitSix:   DigitSix,
		}

		dir, err := os.Getwd()
		if err != nil {
			log.Printf("os.Getwd: %v", err)
			return
		}

		filepath := fmt.Sprintf("%s/template/%s", dir, "otp.html")

		body, err := libs.HTMLToString(filepath, dataEmail)
		if err != nil {
			log.Printf("libs.HTMLToString: %v", err)
			return
		}

		go func() {
			err := libs.SendEmail(subject, []string{props.Email}, body)
			if err != nil {
				log.Printf("libs.SendEmail: %v", err)
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

func (u *UserUsecase) UpdateVerifiedEmailByOTP(props *model.VerifiedEmailOTPRequest) (resp model.Response) {
	_, err := u.repository.GetUserOtpByOtp(context.Background(), sql.NullString{String: props.Otp, Valid: true})
	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "OTP doesnt exist")
		log.Printf("repository.GetUserOtpByOtp %v", err)
		return resp
	}

	err = u.repository.UpdateVerifiedEmailByOTP(context.Background(), sql.NullString{String: props.Otp, Valid: true})

	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Failed updated verify email")
		log.Printf("libs.CustomResponse %v", err)
		return resp
	}

	resp.Status = libs.CustomResponse(http.StatusOK, "Success verified email")
	return resp
}
