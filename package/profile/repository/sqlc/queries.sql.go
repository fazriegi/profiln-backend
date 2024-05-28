// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: queries.sql

package profile

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

const batchInsertEducationFiles = `-- name: BatchInsertEducationFiles :many
INSERT INTO education_files
  (education_id, url)
SELECT $1::bigint, UNNEST($2::text[])
RETURNING id, education_id, url
`

type BatchInsertEducationFilesParams struct {
	EducationID int64
	Url         []string
}

func (q *Queries) BatchInsertEducationFiles(ctx context.Context, arg BatchInsertEducationFilesParams) ([]EducationFile, error) {
	rows, err := q.db.QueryContext(ctx, batchInsertEducationFiles, arg.EducationID, pq.Array(arg.Url))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []EducationFile
	for rows.Next() {
		var i EducationFile
		if err := rows.Scan(&i.ID, &i.EducationID, &i.Url); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const batchInsertEducationSkills = `-- name: BatchInsertEducationSkills :exec
INSERT INTO education_skills (education_id, user_skill_id)
SELECT $1::bigint, unnest($2::bigint[])
ON CONFLICT (education_id, user_skill_id) DO NOTHING
`

type BatchInsertEducationSkillsParams struct {
	EducationID int64
	UserSkillID []int64
}

func (q *Queries) BatchInsertEducationSkills(ctx context.Context, arg BatchInsertEducationSkillsParams) error {
	_, err := q.db.ExecContext(ctx, batchInsertEducationSkills, arg.EducationID, pq.Array(arg.UserSkillID))
	return err
}

const batchInsertSkills = `-- name: BatchInsertSkills :exec
INSERT INTO skills (name)
SELECT unnest($1::text[])
ON CONFLICT (name) DO NOTHING
`

func (q *Queries) BatchInsertSkills(ctx context.Context, names []string) error {
	_, err := q.db.ExecContext(ctx, batchInsertSkills, pq.Array(names))
	return err
}

const batchInsertUserMainSkills = `-- name: BatchInsertUserMainSkills :many
WITH exist_skills AS (
    SELECT id, name
    FROM skills
    WHERE name = ANY($3::text[])
)
INSERT INTO user_skills (user_id, skill_id, main_skill)
SELECT
    $1::bigint,
    es.id,
    $2::boolean
FROM exist_skills es
WHERE es.name = ANY($3::text[])
ON CONFLICT (user_id, skill_id) DO UPDATE
SET main_skill = true
RETURNING id
`

type BatchInsertUserMainSkillsParams struct {
	UserID      int64
	IsMainSkill bool
	Names       []string
}

// start get skills id
// end get skills id
// start insert user skills if not exist
func (q *Queries) BatchInsertUserMainSkills(ctx context.Context, arg BatchInsertUserMainSkillsParams) ([]int64, error) {
	rows, err := q.db.QueryContext(ctx, batchInsertUserMainSkills, arg.UserID, arg.IsMainSkill, pq.Array(arg.Names))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const batchInsertUserSkills = `-- name: BatchInsertUserSkills :many

WITH exist_skills AS (
    SELECT id, name
    FROM skills
    WHERE name = ANY($3::text[])
)
INSERT INTO user_skills (user_id, skill_id, main_skill)
SELECT
    $1::bigint,
    es.id,
    $2::boolean
FROM exist_skills es
WHERE es.name = ANY($3::text[])
ON CONFLICT (user_id, skill_id) DO NOTHING
RETURNING id
`

type BatchInsertUserSkillsParams struct {
	UserID      int64
	IsMainSkill bool
	Names       []string
}

// end insert user skills if not exist
// start get skills id
// end get skills id
// start insert user skills if not exist
func (q *Queries) BatchInsertUserSkills(ctx context.Context, arg BatchInsertUserSkillsParams) ([]int64, error) {
	rows, err := q.db.QueryContext(ctx, batchInsertUserSkills, arg.UserID, arg.IsMainSkill, pq.Array(arg.Names))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const batchInsertWorkExperienceFiles = `-- name: BatchInsertWorkExperienceFiles :many
INSERT INTO work_experience_files
  (work_experience_id, url)
SELECT $1::bigint, UNNEST($2::text[])
RETURNING id, work_experience_id, url
`

type BatchInsertWorkExperienceFilesParams struct {
	WorkExperienceID int64
	Url              []string
}

func (q *Queries) BatchInsertWorkExperienceFiles(ctx context.Context, arg BatchInsertWorkExperienceFilesParams) ([]WorkExperienceFile, error) {
	rows, err := q.db.QueryContext(ctx, batchInsertWorkExperienceFiles, arg.WorkExperienceID, pq.Array(arg.Url))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []WorkExperienceFile
	for rows.Next() {
		var i WorkExperienceFile
		if err := rows.Scan(&i.ID, &i.WorkExperienceID, &i.Url); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const batchInsertWorkExperienceSkills = `-- name: BatchInsertWorkExperienceSkills :exec
INSERT INTO work_experience_skills (work_experience_id, user_skill_id)
SELECT $1::bigint, unnest($2::bigint[])
ON CONFLICT (work_experience_id, user_skill_id) DO NOTHING
`

type BatchInsertWorkExperienceSkillsParams struct {
	WorkExperienceID int64
	UserSkillID      []int64
}

func (q *Queries) BatchInsertWorkExperienceSkills(ctx context.Context, arg BatchInsertWorkExperienceSkillsParams) error {
	_, err := q.db.ExecContext(ctx, batchInsertWorkExperienceSkills, arg.WorkExperienceID, pq.Array(arg.UserSkillID))
	return err
}

const deleteEducationFilesByEducationId = `-- name: DeleteEducationFilesByEducationId :exec
DELETE FROM education_files
WHERE education_id = $1::bigint
`

func (q *Queries) DeleteEducationFilesByEducationId(ctx context.Context, educationID int64) error {
	_, err := q.db.ExecContext(ctx, deleteEducationFilesByEducationId, educationID)
	return err
}

const deleteEducationSkillsByEducation = `-- name: DeleteEducationSkillsByEducation :many
DELETE FROM education_skills
WHERE education_id = $1::bigint
RETURNING user_skill_id
`

func (q *Queries) DeleteEducationSkillsByEducation(ctx context.Context, educationID int64) ([]sql.NullInt64, error) {
	rows, err := q.db.QueryContext(ctx, deleteEducationSkillsByEducation, educationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []sql.NullInt64
	for rows.Next() {
		var user_skill_id sql.NullInt64
		if err := rows.Scan(&user_skill_id); err != nil {
			return nil, err
		}
		items = append(items, user_skill_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const deleteWorkExperienceFilesByWorkExperienceId = `-- name: DeleteWorkExperienceFilesByWorkExperienceId :exec
DELETE FROM work_experience_files
WHERE work_experience_id = $1::bigint
`

func (q *Queries) DeleteWorkExperienceFilesByWorkExperienceId(ctx context.Context, workExperienceID int64) error {
	_, err := q.db.ExecContext(ctx, deleteWorkExperienceFilesByWorkExperienceId, workExperienceID)
	return err
}

const deleteWorkExperienceSkillsByWorkExperience = `-- name: DeleteWorkExperienceSkillsByWorkExperience :many
DELETE FROM work_experience_skills
WHERE work_experience_id = $1::bigint
RETURNING user_skill_id
`

func (q *Queries) DeleteWorkExperienceSkillsByWorkExperience(ctx context.Context, workExperienceID int64) ([]sql.NullInt64, error) {
	rows, err := q.db.QueryContext(ctx, deleteWorkExperienceSkillsByWorkExperience, workExperienceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []sql.NullInt64
	for rows.Next() {
		var user_skill_id sql.NullInt64
		if err := rows.Scan(&user_skill_id); err != nil {
			return nil, err
		}
		items = append(items, user_skill_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEducationById = `-- name: GetEducationById :one
SELECT id, user_id, school_id, degree, field_of_study, gpa, start_date, finish_date, description, created_at, updated_at FROM educations
WHERE id = $1::bigint
LIMIT 1
`

func (q *Queries) GetEducationById(ctx context.Context, id int64) (Education, error) {
	row := q.db.QueryRowContext(ctx, getEducationById, id)
	var i Education
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SchoolID,
		&i.Degree,
		&i.FieldOfStudy,
		&i.Gpa,
		&i.StartDate,
		&i.FinishDate,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getProfile = `-- name: GetProfile :many
SELECT users.full_name, users.bio, user_social_links.url, social_links.name, user_skills.main_skill, skills.name, users.followers_count, users.followings_count
FROM users
LEFT JOIN user_social_links
ON users.id = user_social_links.user_id
LEFT JOIN social_links
ON user_social_links.social_link_id = social_links.id
LEFT JOIN user_skills
ON users.id = user_skills.user_id
LEFT JOIN skills
ON user_skills.skill_id = skills.id
WHERE user_skills.main_skill = TRUE AND users.id = $1
`

type GetProfileRow struct {
	FullName        string
	Bio             sql.NullString
	Url             sql.NullString
	Name            sql.NullString
	MainSkill       sql.NullBool
	Name_2          sql.NullString
	FollowersCount  sql.NullInt32
	FollowingsCount sql.NullInt32
}

func (q *Queries) GetProfile(ctx context.Context, id int64) ([]GetProfileRow, error) {
	rows, err := q.db.QueryContext(ctx, getProfile, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProfileRow
	for rows.Next() {
		var i GetProfileRow
		if err := rows.Scan(
			&i.FullName,
			&i.Bio,
			&i.Url,
			&i.Name,
			&i.MainSkill,
			&i.Name_2,
			&i.FollowersCount,
			&i.FollowingsCount,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSkills = `-- name: GetSkills :many
SELECT id, name, COUNT(id) OVER () AS total_rows
FROM skills
OFFSET $1
LIMIT $2
`

type GetSkillsParams struct {
	Offset int32
	Limit  int32
}

type GetSkillsRow struct {
	ID        int64
	Name      string
	TotalRows int64
}

func (q *Queries) GetSkills(ctx context.Context, arg GetSkillsParams) ([]GetSkillsRow, error) {
	rows, err := q.db.QueryContext(ctx, getSkills, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSkillsRow
	for rows.Next() {
		var i GetSkillsRow
		if err := rows.Scan(&i.ID, &i.Name, &i.TotalRows); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserAbout = `-- name: GetUserAbout :one
SELECT ud.id, ud.user_id, ud.about, ud.updated_at, ud.created_at, u.id, u.email, u.full_name
FROM users u
LEFT JOIN user_details ud
ON u.id = ud.user_id
WHERE u.id = $1
LIMIT 1
`

type GetUserAboutRow struct {
	ID        sql.NullInt64
	UserID    sql.NullInt64
	About     sql.NullString
	UpdatedAt sql.NullTime
	CreatedAt sql.NullTime
	ID_2      int64
	Email     string
	FullName  string
}

func (q *Queries) GetUserAbout(ctx context.Context, id int64) (GetUserAboutRow, error) {
	row := q.db.QueryRowContext(ctx, getUserAbout, id)
	var i GetUserAboutRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.About,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.ID_2,
		&i.Email,
		&i.FullName,
	)
	return i, err
}

const getUserAvatarById = `-- name: GetUserAvatarById :one
SELECT avatar_url
FROM users
WHERE users.id = $1
LIMIT 1
`

func (q *Queries) GetUserAvatarById(ctx context.Context, id int64) (sql.NullString, error) {
	row := q.db.QueryRowContext(ctx, getUserAvatarById, id)
	var avatar_url sql.NullString
	err := row.Scan(&avatar_url)
	return avatar_url, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, email, password, full_name, verified_email, avatar_url, bio, open_to_work, created_at, updated_at, deleted_at, followers_count, followings_count
FROM users
WHERE users.id = $1
LIMIT 1
`

func (q *Queries) GetUserById(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
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
		&i.FollowersCount,
		&i.FollowingsCount,
	)
	return i, err
}

const getUserCertificates = `-- name: GetUserCertificates :many
SELECT u.id, u.email, u.password, u.full_name, u.verified_email, u.avatar_url, u.bio, u.open_to_work, u.created_at, u.updated_at, u.deleted_at, u.followers_count, u.followings_count, c.id, c.name, c.issue_date, c.expiration_date, c.credential_id, c.url, i.name 
FROM users u 
LEFT JOIN certificates c 
ON u.id = c.user_id 
LEFT JOIN 
issuing_organizations i 
ON c.issuing_organization_id = i.id
WHERE u.id = $1
`

type GetUserCertificatesRow struct {
	ID              int64
	Email           string
	Password        sql.NullString
	FullName        string
	VerifiedEmail   sql.NullBool
	AvatarUrl       sql.NullString
	Bio             sql.NullString
	OpenToWork      sql.NullBool
	CreatedAt       sql.NullTime
	UpdatedAt       sql.NullTime
	DeletedAt       sql.NullTime
	FollowersCount  sql.NullInt32
	FollowingsCount sql.NullInt32
	ID_2            sql.NullInt64
	Name            sql.NullString
	IssueDate       sql.NullTime
	ExpirationDate  sql.NullTime
	CredentialID    sql.NullString
	Url             sql.NullString
	Name_2          sql.NullString
}

func (q *Queries) GetUserCertificates(ctx context.Context, id int64) ([]GetUserCertificatesRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserCertificates, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserCertificatesRow
	for rows.Next() {
		var i GetUserCertificatesRow
		if err := rows.Scan(
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
			&i.FollowersCount,
			&i.FollowingsCount,
			&i.ID_2,
			&i.Name,
			&i.IssueDate,
			&i.ExpirationDate,
			&i.CredentialID,
			&i.Url,
			&i.Name_2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserDetail = `-- name: GetUserDetail :one
SELECT id, user_id, phone_number, gender, location, portfolio_url, about, hide_phone_number, created_at, updated_at FROM user_details
WHERE user_id = $1::bigint
LIMIT 1
`

func (q *Queries) GetUserDetail(ctx context.Context, userID int64) (UserDetail, error) {
	row := q.db.QueryRowContext(ctx, getUserDetail, userID)
	var i UserDetail
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.PhoneNumber,
		&i.Gender,
		&i.Location,
		&i.PortfolioUrl,
		&i.About,
		&i.HidePhoneNumber,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserEducationFileURLs = `-- name: GetUserEducationFileURLs :many
SELECT url FROM education_files
WHERE education_id = $1::bigint
`

func (q *Queries) GetUserEducationFileURLs(ctx context.Context, educationID int64) ([]sql.NullString, error) {
	rows, err := q.db.QueryContext(ctx, getUserEducationFileURLs, educationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []sql.NullString
	for rows.Next() {
		var url sql.NullString
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		items = append(items, url)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserSkillIDsByName = `-- name: GetUserSkillIDsByName :many
SELECT us.id FROM user_skills us
JOIN skills s ON us.skill_id = s.id
WHERE s.name = ANY($1::text[])
`

func (q *Queries) GetUserSkillIDsByName(ctx context.Context, name []string) ([]int64, error) {
	rows, err := q.db.QueryContext(ctx, getUserSkillIDsByName, pq.Array(name))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserSkillsAndLocation = `-- name: GetUserSkillsAndLocation :many
SELECT u.id, u.email, u.full_name, s.id, s.name, ud.id, ud.user_id, ud.phone_number, ud.gender, ud.location, ud.portfolio_url, ud.about, ud.hide_phone_number, ud.created_at, ud.updated_at
FROM users u
LEFT JOIN user_skills us
ON u.id = us.user_id
LEFT JOIN skills s
ON us.skill_id = s.id
LEFT JOIN user_details ud
ON u.id = ud.user_id
WHERE u.id = $1
`

type GetUserSkillsAndLocationRow struct {
	ID              int64
	Email           string
	FullName        string
	ID_2            sql.NullInt64
	Name            sql.NullString
	ID_3            sql.NullInt64
	UserID          sql.NullInt64
	PhoneNumber     sql.NullString
	Gender          sql.NullString
	Location        sql.NullString
	PortfolioUrl    sql.NullString
	About           sql.NullString
	HidePhoneNumber sql.NullBool
	CreatedAt       sql.NullTime
	UpdatedAt       sql.NullTime
}

func (q *Queries) GetUserSkillsAndLocation(ctx context.Context, id int64) ([]GetUserSkillsAndLocationRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserSkillsAndLocation, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserSkillsAndLocationRow
	for rows.Next() {
		var i GetUserSkillsAndLocationRow
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.FullName,
			&i.ID_2,
			&i.Name,
			&i.ID_3,
			&i.UserID,
			&i.PhoneNumber,
			&i.Gender,
			&i.Location,
			&i.PortfolioUrl,
			&i.About,
			&i.HidePhoneNumber,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getWorkExperienceById = `-- name: GetWorkExperienceById :one
SELECT id, user_id, job_title, company_id, location, start_date, finish_date, description, created_at, updated_at, location_type, employment_type FROM work_experiences
WHERE id = $1::bigint
LIMIT 1
`

func (q *Queries) GetWorkExperienceById(ctx context.Context, id int64) (WorkExperience, error) {
	row := q.db.QueryRowContext(ctx, getWorkExperienceById, id)
	var i WorkExperience
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.JobTitle,
		&i.CompanyID,
		&i.Location,
		&i.StartDate,
		&i.FinishDate,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LocationType,
		&i.EmploymentType,
	)
	return i, err
}

const getWorkExperienceFileURLs = `-- name: GetWorkExperienceFileURLs :many
SELECT url FROM work_experience_files
WHERE work_experience_id = $1::bigint
`

func (q *Queries) GetWorkExperienceFileURLs(ctx context.Context, workExperienceID int64) ([]sql.NullString, error) {
	rows, err := q.db.QueryContext(ctx, getWorkExperienceFileURLs, workExperienceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []sql.NullString
	for rows.Next() {
		var url sql.NullString
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		items = append(items, url)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertCertificate = `-- name: InsertCertificate :one
INSERT INTO certificates (
  user_id, name, issuing_organization_id, issue_date, expiration_date, credential_id, url
) VALUES (
   $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, user_id, name, issuing_organization_id, issue_date, expiration_date, credential_id, url, created_at, updated_at
`

type InsertCertificateParams struct {
	UserID                sql.NullInt64
	Name                  sql.NullString
	IssuingOrganizationID sql.NullInt64
	IssueDate             sql.NullTime
	ExpirationDate        sql.NullTime
	CredentialID          sql.NullString
	Url                   sql.NullString
}

func (q *Queries) InsertCertificate(ctx context.Context, arg InsertCertificateParams) (Certificate, error) {
	row := q.db.QueryRowContext(ctx, insertCertificate,
		arg.UserID,
		arg.Name,
		arg.IssuingOrganizationID,
		arg.IssueDate,
		arg.ExpirationDate,
		arg.CredentialID,
		arg.Url,
	)
	var i Certificate
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.IssuingOrganizationID,
		&i.IssueDate,
		&i.ExpirationDate,
		&i.CredentialID,
		&i.Url,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const insertCompany = `-- name: InsertCompany :one
INSERT INTO companies (name)
VALUES ($1)
ON CONFLICT (name) DO NOTHING
RETURNING id, name
`

func (q *Queries) InsertCompany(ctx context.Context, name string) (Company, error) {
	row := q.db.QueryRowContext(ctx, insertCompany, name)
	var i Company
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const insertEducation = `-- name: InsertEducation :one
INSERT INTO educations (
  user_id, school_id, degree, field_of_study, gpa, start_date, finish_date
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, user_id, school_id, degree, field_of_study, gpa, start_date, finish_date, description, created_at, updated_at
`

type InsertEducationParams struct {
	UserID       sql.NullInt64
	SchoolID     sql.NullInt64
	Degree       sql.NullString
	FieldOfStudy sql.NullString
	Gpa          sql.NullString
	StartDate    sql.NullTime
	FinishDate   sql.NullTime
}

func (q *Queries) InsertEducation(ctx context.Context, arg InsertEducationParams) (Education, error) {
	row := q.db.QueryRowContext(ctx, insertEducation,
		arg.UserID,
		arg.SchoolID,
		arg.Degree,
		arg.FieldOfStudy,
		arg.Gpa,
		arg.StartDate,
		arg.FinishDate,
	)
	var i Education
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SchoolID,
		&i.Degree,
		&i.FieldOfStudy,
		&i.Gpa,
		&i.StartDate,
		&i.FinishDate,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const insertIssuingOrganization = `-- name: InsertIssuingOrganization :one
INSERT INTO issuing_organizations (name)
VALUES ($1)
ON CONFLICT (name) DO NOTHING
RETURNING id, name
`

func (q *Queries) InsertIssuingOrganization(ctx context.Context, name string) (IssuingOrganization, error) {
	row := q.db.QueryRowContext(ctx, insertIssuingOrganization, name)
	var i IssuingOrganization
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const insertSchool = `-- name: InsertSchool :one
INSERT INTO schools (name)
VALUES ($1)
ON CONFLICT (name) DO NOTHING
RETURNING id, name
`

func (q *Queries) InsertSchool(ctx context.Context, name string) (School, error) {
	row := q.db.QueryRowContext(ctx, insertSchool, name)
	var i School
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const insertSkill = `-- name: InsertSkill :one
INSERT INTO skills (name)
VALUES ($1)
ON CONFLICT (name) DO NOTHING
RETURNING id, name
`

func (q *Queries) InsertSkill(ctx context.Context, name string) (Skill, error) {
	row := q.db.QueryRowContext(ctx, insertSkill, name)
	var i Skill
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const insertUserAvatar = `-- name: InsertUserAvatar :exec
UPDATE users
SET avatar_url = $1
WHERE id = $2
RETURNING id, email, password, full_name, verified_email, avatar_url, bio, open_to_work, created_at, updated_at, deleted_at, followers_count, followings_count
`

type InsertUserAvatarParams struct {
	AvatarUrl sql.NullString
	ID        int64
}

func (q *Queries) InsertUserAvatar(ctx context.Context, arg InsertUserAvatarParams) error {
	_, err := q.db.ExecContext(ctx, insertUserAvatar, arg.AvatarUrl, arg.ID)
	return err
}

const insertUserDetail = `-- name: InsertUserDetail :one
INSERT INTO user_details (
  user_id, phone_number, gender
) VALUES (
  $1, $2, $3
)
RETURNING id, user_id, phone_number, gender, location, portfolio_url, about, hide_phone_number, created_at, updated_at
`

type InsertUserDetailParams struct {
	UserID      sql.NullInt64
	PhoneNumber sql.NullString
	Gender      sql.NullString
}

func (q *Queries) InsertUserDetail(ctx context.Context, arg InsertUserDetailParams) (UserDetail, error) {
	row := q.db.QueryRowContext(ctx, insertUserDetail, arg.UserID, arg.PhoneNumber, arg.Gender)
	var i UserDetail
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.PhoneNumber,
		&i.Gender,
		&i.Location,
		&i.PortfolioUrl,
		&i.About,
		&i.HidePhoneNumber,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const insertUserDetailAbout = `-- name: InsertUserDetailAbout :one
INSERT INTO user_details (
  user_id, about
) VALUES (
  $1, $2
)
RETURNING id, user_id, phone_number, gender, location, portfolio_url, about, hide_phone_number, created_at, updated_at
`

type InsertUserDetailAboutParams struct {
	UserID sql.NullInt64
	About  sql.NullString
}

func (q *Queries) InsertUserDetailAbout(ctx context.Context, arg InsertUserDetailAboutParams) (UserDetail, error) {
	row := q.db.QueryRowContext(ctx, insertUserDetailAbout, arg.UserID, arg.About)
	var i UserDetail
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.PhoneNumber,
		&i.Gender,
		&i.Location,
		&i.PortfolioUrl,
		&i.About,
		&i.HidePhoneNumber,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const insertUserSkill = `-- name: InsertUserSkill :one
INSERT INTO user_skills (
  user_id, skill_id, main_skill
) VALUES (
   $1, $2, $3
)
RETURNING id, user_id, skill_id, main_skill
`

type InsertUserSkillParams struct {
	UserID    sql.NullInt64
	SkillID   sql.NullInt64
	MainSkill sql.NullBool
}

func (q *Queries) InsertUserSkill(ctx context.Context, arg InsertUserSkillParams) (UserSkill, error) {
	row := q.db.QueryRowContext(ctx, insertUserSkill, arg.UserID, arg.SkillID, arg.MainSkill)
	var i UserSkill
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SkillID,
		&i.MainSkill,
	)
	return i, err
}

const insertWorkExperience = `-- name: InsertWorkExperience :one
INSERT INTO work_experiences (
  user_id, job_title, company_id, employment_type, location, location_type, start_date, finish_date, description
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING id, user_id, job_title, company_id, location, start_date, finish_date, description, created_at, updated_at, location_type, employment_type
`

type InsertWorkExperienceParams struct {
	UserID         sql.NullInt64
	JobTitle       sql.NullString
	CompanyID      sql.NullInt64
	EmploymentType sql.NullString
	Location       sql.NullString
	LocationType   sql.NullString
	StartDate      sql.NullTime
	FinishDate     sql.NullTime
	Description    sql.NullString
}

func (q *Queries) InsertWorkExperience(ctx context.Context, arg InsertWorkExperienceParams) (WorkExperience, error) {
	row := q.db.QueryRowContext(ctx, insertWorkExperience,
		arg.UserID,
		arg.JobTitle,
		arg.CompanyID,
		arg.EmploymentType,
		arg.Location,
		arg.LocationType,
		arg.StartDate,
		arg.FinishDate,
		arg.Description,
	)
	var i WorkExperience
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.JobTitle,
		&i.CompanyID,
		&i.Location,
		&i.StartDate,
		&i.FinishDate,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LocationType,
		&i.EmploymentType,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one

UPDATE users
SET full_name = $1,
    avatar_url = $2
WHERE id = $3
RETURNING full_name, avatar_url
`

type UpdateUserParams struct {
	FullName  string
	AvatarUrl sql.NullString
	ID        int64
}

type UpdateUserRow struct {
	FullName  string
	AvatarUrl sql.NullString
}

// end insert user skills if not exist
func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (UpdateUserRow, error) {
	row := q.db.QueryRowContext(ctx, updateUser, arg.FullName, arg.AvatarUrl, arg.ID)
	var i UpdateUserRow
	err := row.Scan(&i.FullName, &i.AvatarUrl)
	return i, err
}

const updateUserCertificate = `-- name: UpdateUserCertificate :one
UPDATE certificates 
SET name = $1::text,
    issuing_organization_id = $2::bigint,
    issue_date = $3::date, 
    expiration_date = $4::date, 
    credential_id = $5::text, 
    url = $6::text
WHERE id = $7::bigint AND user_id = $8::bigint
RETURNING id
`

type UpdateUserCertificateParams struct {
	Name                  string
	IssuingOrganizationID int64
	IssueDate             time.Time
	ExpirationDate        time.Time
	CredentialID          string
	Url                   string
	ID                    int64
	UserID                int64
}

func (q *Queries) UpdateUserCertificate(ctx context.Context, arg UpdateUserCertificateParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, updateUserCertificate,
		arg.Name,
		arg.IssuingOrganizationID,
		arg.IssueDate,
		arg.ExpirationDate,
		arg.CredentialID,
		arg.Url,
		arg.ID,
		arg.UserID,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const updateUserDetail = `-- name: UpdateUserDetail :one
UPDATE user_details
SET phone_number = $2,
    gender = $3,
    location = $4,
    portfolio_url = $5,
    about = $6,
    hide_phone_number = $7
WHERE user_id = $1
RETURNING id, phone_number, gender, location, portfolio_url, about, hide_phone_number
`

type UpdateUserDetailParams struct {
	UserID          sql.NullInt64
	PhoneNumber     sql.NullString
	Gender          sql.NullString
	Location        sql.NullString
	PortfolioUrl    sql.NullString
	About           sql.NullString
	HidePhoneNumber sql.NullBool
}

type UpdateUserDetailRow struct {
	ID              int64
	PhoneNumber     sql.NullString
	Gender          sql.NullString
	Location        sql.NullString
	PortfolioUrl    sql.NullString
	About           sql.NullString
	HidePhoneNumber sql.NullBool
}

func (q *Queries) UpdateUserDetail(ctx context.Context, arg UpdateUserDetailParams) (UpdateUserDetailRow, error) {
	row := q.db.QueryRowContext(ctx, updateUserDetail,
		arg.UserID,
		arg.PhoneNumber,
		arg.Gender,
		arg.Location,
		arg.PortfolioUrl,
		arg.About,
		arg.HidePhoneNumber,
	)
	var i UpdateUserDetailRow
	err := row.Scan(
		&i.ID,
		&i.PhoneNumber,
		&i.Gender,
		&i.Location,
		&i.PortfolioUrl,
		&i.About,
		&i.HidePhoneNumber,
	)
	return i, err
}

const updateUserDetailAbout = `-- name: UpdateUserDetailAbout :exec
UPDATE user_details
SET about = $1::text
WHERE user_id = $2::bigint
`

type UpdateUserDetailAboutParams struct {
	About  string
	UserID int64
}

func (q *Queries) UpdateUserDetailAbout(ctx context.Context, arg UpdateUserDetailAboutParams) error {
	_, err := q.db.ExecContext(ctx, updateUserDetailAbout, arg.About, arg.UserID)
	return err
}

const updateUserDetailByUserId = `-- name: UpdateUserDetailByUserId :one
UPDATE user_details
SET hide_phone_number = $2,
    phone_number = $3,
    gender = $4
WHERE user_id = $1
RETURNING hide_phone_number, phone_number, gender
`

type UpdateUserDetailByUserIdParams struct {
	UserID          sql.NullInt64
	HidePhoneNumber sql.NullBool
	PhoneNumber     sql.NullString
	Gender          sql.NullString
}

type UpdateUserDetailByUserIdRow struct {
	HidePhoneNumber sql.NullBool
	PhoneNumber     sql.NullString
	Gender          sql.NullString
}

func (q *Queries) UpdateUserDetailByUserId(ctx context.Context, arg UpdateUserDetailByUserIdParams) (UpdateUserDetailByUserIdRow, error) {
	row := q.db.QueryRowContext(ctx, updateUserDetailByUserId,
		arg.UserID,
		arg.HidePhoneNumber,
		arg.PhoneNumber,
		arg.Gender,
	)
	var i UpdateUserDetailByUserIdRow
	err := row.Scan(&i.HidePhoneNumber, &i.PhoneNumber, &i.Gender)
	return i, err
}

const updateUserEducation = `-- name: UpdateUserEducation :one
UPDATE educations
SET school_id = $2,
    degree = $3,
    field_of_study = $4,
    gpa = $5,
    start_date = $6,
    finish_date = $7,
    description = $8
WHERE id = $1
RETURNING id, user_id, school_id, degree, field_of_study, gpa, start_date, finish_date, description, created_at, updated_at
`

type UpdateUserEducationParams struct {
	ID           int64
	SchoolID     sql.NullInt64
	Degree       sql.NullString
	FieldOfStudy sql.NullString
	Gpa          sql.NullString
	StartDate    sql.NullTime
	FinishDate   sql.NullTime
	Description  sql.NullString
}

func (q *Queries) UpdateUserEducation(ctx context.Context, arg UpdateUserEducationParams) (Education, error) {
	row := q.db.QueryRowContext(ctx, updateUserEducation,
		arg.ID,
		arg.SchoolID,
		arg.Degree,
		arg.FieldOfStudy,
		arg.Gpa,
		arg.StartDate,
		arg.FinishDate,
		arg.Description,
	)
	var i Education
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SchoolID,
		&i.Degree,
		&i.FieldOfStudy,
		&i.Gpa,
		&i.StartDate,
		&i.FinishDate,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUserMainSkillToFalse = `-- name: UpdateUserMainSkillToFalse :exec
UPDATE user_skills 
SET main_skill = false 
WHERE user_id = $1::bigint
AND main_skill = true
`

func (q *Queries) UpdateUserMainSkillToFalse(ctx context.Context, userID int64) error {
	_, err := q.db.ExecContext(ctx, updateUserMainSkillToFalse, userID)
	return err
}

const updateUserWorkExperience = `-- name: UpdateUserWorkExperience :one
UPDATE work_experiences
SET job_title = $2,
    company_id = $3,
    employment_type = $4,
    location = $5,
    location_type = $6,
    start_date = $7,
    finish_date = $8,
    description = $9
WHERE id = $1
RETURNING id, user_id, job_title, company_id, location, start_date, finish_date, description, created_at, updated_at, location_type, employment_type
`

type UpdateUserWorkExperienceParams struct {
	ID             int64
	JobTitle       sql.NullString
	CompanyID      sql.NullInt64
	EmploymentType sql.NullString
	Location       sql.NullString
	LocationType   sql.NullString
	StartDate      sql.NullTime
	FinishDate     sql.NullTime
	Description    sql.NullString
}

func (q *Queries) UpdateUserWorkExperience(ctx context.Context, arg UpdateUserWorkExperienceParams) (WorkExperience, error) {
	row := q.db.QueryRowContext(ctx, updateUserWorkExperience,
		arg.ID,
		arg.JobTitle,
		arg.CompanyID,
		arg.EmploymentType,
		arg.Location,
		arg.LocationType,
		arg.StartDate,
		arg.FinishDate,
		arg.Description,
	)
	var i WorkExperience
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.JobTitle,
		&i.CompanyID,
		&i.Location,
		&i.StartDate,
		&i.FinishDate,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LocationType,
		&i.EmploymentType,
	)
	return i, err
}

const upsertUserSocialLink = `-- name: UpsertUserSocialLink :exec
WITH social_link AS (
    SELECT id
    FROM social_links
    WHERE name = $2
    LIMIT 1
)
INSERT INTO user_social_links (user_id, social_link_id, url)
SELECT $1, sl.id, $3
FROM social_link sl
ON CONFLICT (user_id, social_link_id) DO UPDATE
SET url = EXCLUDED.url
`

type UpsertUserSocialLinkParams struct {
	UserID sql.NullInt64
	Name   string
	Url    sql.NullString
}

func (q *Queries) UpsertUserSocialLink(ctx context.Context, arg UpsertUserSocialLinkParams) error {
	_, err := q.db.ExecContext(ctx, upsertUserSocialLink, arg.UserID, arg.Name, arg.Url)
	return err
}
