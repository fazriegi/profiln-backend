// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: queries.sql

package auth

import (
	"context"
	"database/sql"
)

const deleteOtp = `-- name: DeleteOtp :exec
DELETE FROM user_otps WHERE otp = $1
`

func (q *Queries) DeleteOtp(ctx context.Context, otp sql.NullString) error {
	_, err := q.db.ExecContext(ctx, deleteOtp, otp)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, password, full_name, verified_email, avatar_url, bio, open_to_work, created_at, updated_at, deleted_at FROM users
WHERE email = $1
LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.FullName,
		&i.VerifiedEmail,
		&i.AvatarUrl,
		&i.Bio,
		&i.OpenToWork,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getUserOtpByOtp = `-- name: GetUserOtpByOtp :one
SELECT id, user_id, otp 
FROM user_otps 
WHERE otp = $1 
LIMIT 1
`

func (q *Queries) GetUserOtpByOtp(ctx context.Context, otp sql.NullString) (UserOtp, error) {
	row := q.db.QueryRowContext(ctx, getUserOtpByOtp, otp)
	var i UserOtp
	err := row.Scan(&i.ID, &i.UserID, &i.Otp)
	return i, err
}

const insertOtp = `-- name: InsertOtp :one
INSERT INTO user_otps (
  user_id, otp
) VALUES (
  $1, $2
)
RETURNING id, user_id, otp
`

type InsertOtpParams struct {
	UserID sql.NullInt64
	Otp    sql.NullString
}

func (q *Queries) InsertOtp(ctx context.Context, arg InsertOtpParams) (UserOtp, error) {
	row := q.db.QueryRowContext(ctx, insertOtp, arg.UserID, arg.Otp)
	var i UserOtp
	err := row.Scan(&i.ID, &i.UserID, &i.Otp)
	return i, err
}

const insertUser = `-- name: InsertUser :one
INSERT INTO users (
  email, password, full_name, verified_email
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, email, password, full_name, verified_email, avatar_url, bio, open_to_work, created_at, updated_at, deleted_at
`

type InsertUserParams struct {
	Email         string
	Password      sql.NullString
	FullName      string
	VerifiedEmail sql.NullBool
}

func (q *Queries) InsertUser(ctx context.Context, arg InsertUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, insertUser,
		arg.Email,
		arg.Password,
		arg.FullName,
		arg.VerifiedEmail,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.FullName,
		&i.VerifiedEmail,
		&i.AvatarUrl,
		&i.Bio,
		&i.OpenToWork,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const updateUserPassword = `-- name: UpdateUserPassword :exec
UPDATE users
SET password = $2
WHERE id = $1
RETURNING id, email, password, full_name, verified_email, avatar_url, bio, open_to_work, created_at, updated_at, deleted_at
`

type UpdateUserPasswordParams struct {
	ID       int64
	Password sql.NullString
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.db.ExecContext(ctx, updateUserPassword, arg.ID, arg.Password)
	return err
}

const updateVerifiedEmail = `-- name: UpdateVerifiedEmail :one
UPDATE users
SET verified_email = TRUE
FROM user_otps 
WHERE users.id = user_otps.user_id AND user_otps.otp = $1 AND users.email = $2
RETURNING users.id
`

type UpdateVerifiedEmailParams struct {
	Otp   sql.NullString
	Email string
}

func (q *Queries) UpdateVerifiedEmail(ctx context.Context, arg UpdateVerifiedEmailParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, updateVerifiedEmail, arg.Otp, arg.Email)
	var id int64
	err := row.Scan(&id)
	return id, err
}
