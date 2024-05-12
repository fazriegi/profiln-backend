package profile

import (
	"context"
	"database/sql"
	profileSqlc "profiln-be/package/profile/repository/sqlc"
)

type IProfileRepository interface {
	InsertUserDetail(arg profileSqlc.InsertUserDetailParams) (profileSqlc.UserDetail, error)
	InsertUserDetailAbout(arg profileSqlc.InsertUserDetailAboutParams) error
	InsertCompany(name string) (profileSqlc.Company, error)
	InsertEducation(arg profileSqlc.InsertEducationParams) (profileSqlc.Education, error)
	InsertEmploymentType(name string) (profileSqlc.EmploymentType, error)
	InsertLocationType(name string) (profileSqlc.LocationType, error)
	InsertSchool(name string) (profileSqlc.School, error)
	InsertCertificate(arg profileSqlc.InsertCertificateParams) (profileSqlc.Certificate, error)
	InsertIssuingOrganization(name string) (profileSqlc.IssuingOrganization, error)
	InsertUserSkill(arg profileSqlc.InsertUserSkillParams) (profileSqlc.UserSkill, error)
	InsertSkill(name string) (profileSqlc.Skill, error)
}

type ProfileRepository struct {
	db    *sql.DB
	query *profileSqlc.Queries
}

func NewProfileRepository(db *sql.DB) IProfileRepository {
	return &ProfileRepository{
		db:    db,
		query: profileSqlc.New(db),
	}
}

func (r *ProfileRepository) InsertUserDetail(arg profileSqlc.InsertUserDetailParams) (profileSqlc.UserDetail, error) {
	userDetail, err := r.query.InsertUserDetail(context.Background(), arg)

	if err != nil {
		return profileSqlc.UserDetail{}, err
	}

	return userDetail, nil
}

func (r *ProfileRepository) InsertUserAvatar(arg profileSqlc.InsertUserAvatarParams) error {
	err := r.query.InsertUserAvatar(context.Background(), arg)

	if err != nil {
		return err
	}

	return nil
}

func (r *ProfileRepository) InsertUserDetailAbout(arg profileSqlc.InsertUserDetailAboutParams) error {
	err := r.query.InsertUserDetailAbout(context.Background(), arg)

	if err != nil {
		return err
	}

	return nil
}

func (r *ProfileRepository) InsertCompany(name string) (profileSqlc.Company, error) {
	arg := sql.NullString{String: name, Valid: true}
	company, err := r.query.InsertCompany(context.Background(), arg)

	if err != nil {
		return profileSqlc.Company{}, err
	}

	return company, nil
}

func (r *ProfileRepository) InsertEducation(arg profileSqlc.InsertEducationParams) (profileSqlc.Education, error) {
	education, err := r.query.InsertEducation(context.Background(), arg)

	if err != nil {
		return profileSqlc.Education{}, err
	}

	return education, nil
}

func (r *ProfileRepository) InsertEmploymentType(name string) (profileSqlc.EmploymentType, error) {
	arg := sql.NullString{String: name, Valid: true}
	employmentType, err := r.query.InsertEmploymentType(context.Background(), arg)

	if err != nil {
		return profileSqlc.EmploymentType{}, err
	}

	return employmentType, nil
}

func (r *ProfileRepository) InsertLocationType(name string) (profileSqlc.LocationType, error) {
	arg := sql.NullString{String: name, Valid: true}
	locationType, err := r.query.InsertLocationType(context.Background(), arg)

	if err != nil {
		return profileSqlc.LocationType{}, err
	}

	return locationType, nil
}

func (r *ProfileRepository) InsertSchool(name string) (profileSqlc.School, error) {
	arg := sql.NullString{String: name, Valid: true}
	school, err := r.query.InsertSchool(context.Background(), arg)

	if err != nil {
		return profileSqlc.School{}, err
	}

	return school, nil
}

func (r *ProfileRepository) InsertWorkExperience(arg profileSqlc.InsertWorkExperienceParams) (profileSqlc.WorkExperience, error) {
	workExperience, err := r.query.InsertWorkExperience(context.Background(), arg)

	if err != nil {
		return profileSqlc.WorkExperience{}, err
	}

	return workExperience, nil
}

func (r *ProfileRepository) InsertCertificate(arg profileSqlc.InsertCertificateParams) (profileSqlc.Certificate, error) {
	certificate, err := r.query.InsertCertificate(context.Background(), arg)

	if err != nil {
		return profileSqlc.Certificate{}, err
	}

	return certificate, nil
}

func (r *ProfileRepository) InsertIssuingOrganization(name string) (profileSqlc.IssuingOrganization, error) {
	arg := sql.NullString{String: name, Valid: true}

	issueOrganization, err := r.query.InsertIssuingOrganization(context.Background(), arg)

	if err != nil {
		return profileSqlc.IssuingOrganization{}, err
	}

	return issueOrganization, nil
}

func (r *ProfileRepository) InsertUserSkill(arg profileSqlc.InsertUserSkillParams) (profileSqlc.UserSkill, error) {
	userSkill, err := r.query.InsertUserSkill(context.Background(), arg)

	if err != nil {
		return profileSqlc.UserSkill{}, err
	}

	return userSkill, nil
}

func (r *ProfileRepository) InsertSkill(name string) (profileSqlc.Skill, error) {
	arg := sql.NullString{String: name, Valid: true}

	skill, err := r.query.InsertSkill(context.Background(), arg)

	if err != nil {
		return profileSqlc.Skill{}, err
	}

	return skill, nil
}
