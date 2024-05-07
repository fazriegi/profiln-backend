package auth

import (
	"context"
	"database/sql"
	"fmt"
	authSqlc "profiln-be/package/auth/repository/sqlc"
)

type IAuthRepository interface {
	GetUserByEmail(email string) (authSqlc.User, error)
	UpdateUserPassword(id int64, hashedPassword string) error
	InsertUser(arg authSqlc.InsertUserParams) (authSqlc.User, error)
	UpdateVerifiedEmail(otp string, email string) error
	InsertOtp(id int64, otp string) (authSqlc.UserOtp, error)
	GetUserOtpByOtp(otp string) (authSqlc.UserOtp, error)
	DeleteOtp(otp string) error
}

type AuthRepository struct {
	db    *sql.DB
	query *authSqlc.Queries
}

func NewAuthRepository(db *sql.DB) IAuthRepository {
	return &AuthRepository{
		db:    db,
		query: authSqlc.New(db),
	}
}

func (r *AuthRepository) GetUserByEmail(email string) (authSqlc.User, error) {
	user, err := r.query.GetUserByEmail(context.Background(), email)

	if err != nil {
		return authSqlc.User{}, err
	}

	return user, nil
}

func (r *AuthRepository) UpdateUserPassword(id int64, hashedPassword string) error {
	arg := authSqlc.UpdateUserPasswordParams{
		ID:       id,
		Password: sql.NullString{String: hashedPassword, Valid: true},
	}

	err := r.query.UpdateUserPassword(context.Background(), arg)

	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) InsertUser(arg authSqlc.InsertUserParams) (authSqlc.User, error) {
	user, err := r.query.InsertUser(context.Background(), arg)

	if err != nil {
		return authSqlc.User{}, err
	}

	return user, nil
}

func (r *AuthRepository) UpdateVerifiedEmail(otp string, email string) error {
	updateVerfiedEmailParams := authSqlc.UpdateVerifiedEmailParams{
		Otp:   sql.NullString{String: otp, Valid: true},
		Email: email,
	}
	fmt.Println(updateVerfiedEmailParams)
	_, err := r.query.UpdateVerifiedEmail(context.Background(), updateVerfiedEmailParams)

	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) InsertOtp(id int64, otp string) (authSqlc.UserOtp, error) {
	insertOtpParams := authSqlc.InsertOtpParams{
		UserID: sql.NullInt64{Int64: id, Valid: true},
		Otp:    sql.NullString{String: otp, Valid: true},
	}

	userOtp, err := r.query.InsertOtp(context.Background(), insertOtpParams)

	if err != nil {
		return authSqlc.UserOtp{}, err
	}

	return userOtp, nil
}

func (r *AuthRepository) GetUserOtpByOtp(otp string) (authSqlc.UserOtp, error) {
	arg := sql.NullString{String: otp, Valid: true}
	userOtp, err := r.query.GetUserOtpByOtp(context.Background(), arg)

	if err != nil {
		return authSqlc.UserOtp{}, err
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
