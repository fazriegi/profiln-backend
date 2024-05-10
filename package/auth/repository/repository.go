package auth

import (
	"context"
	"database/sql"
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
	InsertUserDetail(arg authSqlc.InsertUserDetailParams) (authSqlc.UserDetail, error)
	InsertUserDetailAbout(arg authSqlc.InsertUserDetailAboutParams) error
	InsertCompany(name string) (authSqlc.Company, error)
	InsertEducation(arg authSqlc.InsertEducationParams) (authSqlc.Education, error)
	InsertEmploymentType(name string) (authSqlc.EmploymentType, error)
	InsertLocationType(name string) (authSqlc.LocationType, error)
	InsertSchool(name string) (authSqlc.School, error)
	InsertCertificate(arg authSqlc.InsertCertificateParams) (authSqlc.Certificate, error)
	InsertIssuingOrganization(name string) (authSqlc.IssuingOrganization, error)
	InsertUserSkill(arg authSqlc.InsertUserSkillParams) (authSqlc.UserSkill, error)
	InsertSkill(name string) (authSqlc.Skill, error)
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

func (r *AuthRepository) InsertUserDetail(arg authSqlc.InsertUserDetailParams) (authSqlc.UserDetail, error) {
	userDetail, err := r.query.InsertUserDetail(context.Background(), arg)

	if err != nil {
		return authSqlc.UserDetail{}, err
	}

	return userDetail, nil
}

func (r *AuthRepository) InsertUserAvatar(arg authSqlc.InsertUserAvatarParams) error {
	err := r.query.InsertUserAvatar(context.Background(), arg)

	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) InsertUserDetailAbout(arg authSqlc.InsertUserDetailAboutParams) error {
	err := r.query.InsertUserDetailAbout(context.Background(), arg)

	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) InsertCompany(name string) (authSqlc.Company, error) {
	arg := sql.NullString{String: name, Valid: true}
	company, err := r.query.InsertCompany(context.Background(), arg)

	if err != nil {
		return authSqlc.Company{}, err
	}

	return company, nil
}

func (r *AuthRepository) InsertEducation(arg authSqlc.InsertEducationParams) (authSqlc.Education, error) {
	education, err := r.query.InsertEducation(context.Background(), arg)

	if err != nil {
		return authSqlc.Education{}, err
	}

	return education, nil
}

func (r *AuthRepository) InsertEmploymentType(name string) (authSqlc.EmploymentType, error) {
	arg := sql.NullString{String: name, Valid: true}
	employmentType, err := r.query.InsertEmploymentType(context.Background(), arg)

	if err != nil {
		return authSqlc.EmploymentType{}, err
	}

	return employmentType, nil
}

func (r *AuthRepository) InsertLocationType(name string) (authSqlc.LocationType, error) {
	arg := sql.NullString{String: name, Valid: true}
	locationType, err := r.query.InsertLocationType(context.Background(), arg)

	if err != nil {
		return authSqlc.LocationType{}, err
	}

	return locationType, nil
}

func (r *AuthRepository) InsertSchool(name string) (authSqlc.School, error) {
	arg := sql.NullString{String: name, Valid: true}
	school, err := r.query.InsertSchool(context.Background(), arg)

	if err != nil {
		return authSqlc.School{}, err
	}

	return school, nil
}

func (r *AuthRepository) InsertWorkExperience(arg authSqlc.InsertWorkExperienceParams) (authSqlc.WorkExperience, error) {
	workExperience, err := r.query.InsertWorkExperience(context.Background(), arg)

	if err != nil {
		return authSqlc.WorkExperience{}, err
	}

	return workExperience, nil
}

func (r *AuthRepository) InsertCertificate(arg authSqlc.InsertCertificateParams) (authSqlc.Certificate, error) {
	certificate, err := r.query.InsertCertificate(context.Background(), arg)

	if err != nil {
		return authSqlc.Certificate{}, err
	}

	return certificate, nil
}

func (r *AuthRepository) InsertIssuingOrganization(name string) (authSqlc.IssuingOrganization, error) {
	arg := sql.NullString{String: name, Valid: true}

	issueOrganization, err := r.query.InsertIssuingOrganization(context.Background(), arg)

	if err != nil {
		return authSqlc.IssuingOrganization{}, err
	}

	return issueOrganization, nil
}

func (r *AuthRepository) InsertUserSkill(arg authSqlc.InsertUserSkillParams) (authSqlc.UserSkill, error) {
	userSkill, err := r.query.InsertUserSkill(context.Background(), arg)

	if err != nil {
		return authSqlc.UserSkill{}, err
	}

	return userSkill, nil
}

func (r *AuthRepository) InsertSkill(name string) (authSqlc.Skill, error) {
	arg := sql.NullString{String: name, Valid: true}

	skill, err := r.query.InsertSkill(context.Background(), arg)

	if err != nil {
		return authSqlc.Skill{}, err
	}

	return skill, nil
}
