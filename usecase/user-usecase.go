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

func (u *UserUsecase) Register(props *model.UserRegisterRequest) (resp model.Response) {
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

	registerUserParams := repository.InsertUserParams{
		Email:    props.Email,
		Password: sql.NullString{String: hashedPassword, Valid: true},
		// Password:      libs.GetValidString(hashedPassword, user.Password),
		FullName:      props.Fullname,
		VerifiedEmail: sql.NullBool{Bool: false, Valid: true},
	}

	insertUser, err := u.repository.InsertUser(context.Background(), registerUserParams)

	if err != nil {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "Something went wrong while register")
		log.Printf("Register : %v", err)
		return resp
	}

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

	// todo: buat kirim email ke client

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
