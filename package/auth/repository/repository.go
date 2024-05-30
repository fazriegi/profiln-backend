package auth

import (
	"context"
	"database/sql"
	db "profiln-be/db/sqlc"
)

type IAuthRepository interface {
	GetUserByEmail(email string) (db.User, error)
	UpdateUserPassword(id int64, hashedPassword string) error
	InsertUser(arg db.InsertUserParams) (db.User, error)
	UpdateVerifiedEmail(otp string, email string) (db.User, error)
	InsertOtp(id int64, otp string) (db.UserOtp, error)
	GetUserOtpByOtp(otp string) (db.UserOtp, error)
	DeleteOtp(otp string) error
	GetUserOtpByEmail(email string) (db.GetUserOtpByEmailRow, error)
}

type AuthRepository struct {
	dbConn *sql.DB
	query  *db.Queries
}

func NewAuthRepository(dbConn *sql.DB) IAuthRepository {
	return &AuthRepository{
		dbConn: dbConn,
		query:  db.New(dbConn),
	}
}

func (r *AuthRepository) GetUserByEmail(email string) (db.User, error) {
	user, err := r.query.GetUserByEmail(context.Background(), email)

	if err != nil {
		return db.User{}, err
	}

	return user, nil
}

func (r *AuthRepository) UpdateUserPassword(id int64, hashedPassword string) error {
	arg := db.UpdateUserPasswordParams{
		ID:       id,
		Password: sql.NullString{String: hashedPassword, Valid: true},
	}

	err := r.query.UpdateUserPassword(context.Background(), arg)

	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) InsertUser(arg db.InsertUserParams) (db.User, error) {
	user, err := r.query.InsertUser(context.Background(), arg)

	if err != nil {
		return db.User{}, err
	}

	return user, nil
}

func (r *AuthRepository) UpdateVerifiedEmail(otp string, email string) (db.User, error) {
	updateVerfiedEmailParams := db.UpdateVerifiedEmailParams{
		Otp:   sql.NullString{String: otp, Valid: true},
		Email: email,
	}

	user, err := r.query.UpdateVerifiedEmail(context.Background(), updateVerfiedEmailParams)

	if err != nil {
		return db.User{}, err
	}

	return db.User{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}

func (r *AuthRepository) InsertOtp(id int64, otp string) (db.UserOtp, error) {
	insertOtpParams := db.InsertOtpParams{
		UserID: sql.NullInt64{Int64: id, Valid: true},
		Otp:    sql.NullString{String: otp, Valid: true},
	}

	userOtp, err := r.query.InsertOtp(context.Background(), insertOtpParams)

	if err != nil {
		return db.UserOtp{}, err
	}

	return userOtp, nil
}

func (r *AuthRepository) GetUserOtpByOtp(otp string) (db.UserOtp, error) {
	arg := sql.NullString{String: otp, Valid: true}
	userOtp, err := r.query.GetUserOtpByOtp(context.Background(), arg)

	if err != nil {
		return db.UserOtp{}, err
	}

	return userOtp, nil
}

func (r *AuthRepository) DeleteOtp(otp string) error {
	arg := sql.NullString{String: otp, Valid: true}
	err := r.query.DeleteOtp(context.Background(), arg)

	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) GetUserOtpByEmail(email string) (db.GetUserOtpByEmailRow, error) {
	otpByEmail, err := r.query.GetUserOtpByEmail(context.Background(), email)

	if err != nil {
		return db.GetUserOtpByEmailRow{}, err
	}

	return otpByEmail, nil
}
