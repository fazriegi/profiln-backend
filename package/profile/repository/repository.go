package profile

import (
	"context"
	"database/sql"
	"fmt"
	"profiln-be/model"
	profileSqlc "profiln-be/package/profile/repository/sqlc"
)

type IProfileRepository interface {
	InsertUserDetail(arg profileSqlc.InsertUserDetailParams) (profileSqlc.UserDetail, error)
	InsertUserDetailAbout(arg profileSqlc.InsertUserDetailAboutParams) (profileSqlc.UserDetail, error)
	InsertCompany(name string) (profileSqlc.Company, error)
	InsertEducation(arg profileSqlc.InsertEducationParams) (profileSqlc.Education, error)
	InsertEmploymentType(name string) (profileSqlc.EmploymentType, error)
	InsertLocationType(name string) (profileSqlc.LocationType, error)
	InsertSchool(name string) (profileSqlc.School, error)
	InsertCertificate(arg profileSqlc.InsertCertificateParams) (profileSqlc.Certificate, error)
	InsertIssuingOrganization(name string) (profileSqlc.IssuingOrganization, error)
	InsertUserSkill(arg profileSqlc.InsertUserSkillParams) (profileSqlc.UserSkill, error)
	InsertSkill(name string) (profileSqlc.Skill, error)
	InsertWorkExperience(arg profileSqlc.InsertWorkExperienceParams) (profileSqlc.WorkExperience, error)
	InsertUserAvatar(arg profileSqlc.InsertUserAvatarParams) error
	GetUserById(id int64) (profileSqlc.User, error)
	UpdateUserDetailAbout(arg profileSqlc.UpdateUserDetailAboutParams) error
	GetSkills() ([]profileSqlc.Skill, error)
	UpdateProfile(avatar_url string, props *model.UpdateProfileRequest) error
	UpdateAboutMe(userId int64, aboutMe string) error
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

func (r *ProfileRepository) InsertUserDetailAbout(arg profileSqlc.InsertUserDetailAboutParams) (profileSqlc.UserDetail, error) {
	userAbout, err := r.query.InsertUserDetailAbout(context.Background(), arg)

	if err != nil {
		return profileSqlc.UserDetail{}, err
	}

	return userAbout, nil
}

func (r *ProfileRepository) UpdateUserDetailAbout(arg profileSqlc.UpdateUserDetailAboutParams) error {
	err := r.query.UpdateUserDetailAbout(context.Background(), arg)

	if err != nil {
		return err
	}

	return nil
}

func (r *ProfileRepository) InsertCompany(name string) (profileSqlc.Company, error) {
	company, err := r.query.InsertCompany(context.Background(), name)

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
	employmentType, err := r.query.InsertEmploymentType(context.Background(), name)

	if err != nil {
		return profileSqlc.EmploymentType{}, err
	}

	return employmentType, nil
}

func (r *ProfileRepository) InsertLocationType(name string) (profileSqlc.LocationType, error) {
	locationType, err := r.query.InsertLocationType(context.Background(), name)

	if err != nil {
		return profileSqlc.LocationType{}, err
	}

	return locationType, nil
}

func (r *ProfileRepository) InsertSchool(name string) (profileSqlc.School, error) {
	school, err := r.query.InsertSchool(context.Background(), name)

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
	issueOrganization, err := r.query.InsertIssuingOrganization(context.Background(), name)

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
	skill, err := r.query.InsertSkill(context.Background(), name)

	if err != nil {
		return profileSqlc.Skill{}, err
	}

	return skill, nil
}

func (r *ProfileRepository) GetUserAbout(id int64) (profileSqlc.GetUserAboutRow, error) {
	about, err := r.query.GetUserAbout(context.Background(), id)

	if err != nil {
		return profileSqlc.GetUserAboutRow{}, err
	}

	return about, nil
}

func (r *ProfileRepository) GetUserById(id int64) (profileSqlc.User, error) {
	user, err := r.query.GetUserById(context.Background(), id)

	if err != nil {
		return profileSqlc.User{}, err
	}

	return user, nil
}

func (r *ProfileRepository) GetSkills() ([]profileSqlc.Skill, error) {
	skills, err := r.query.GetSkills(context.Background())
	if err != nil {
		return []profileSqlc.Skill{}, err
	}

	return skills, nil
}

func (r *ProfileRepository) UpdateProfile(avatar_url string, props *model.UpdateProfileRequest) error {
	ctx := context.Background()
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin edit profile transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := r.query.WithTx(tx)

	// update users table
	_, err = qtx.UpdateUser(ctx, profileSqlc.UpdateUserParams{
		ID:        props.UserId,
		FullName:  props.Fullname,
		AvatarUrl: sql.NullString{String: avatar_url, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("could not update user: %w", err)
	}

	// update user details table
	_, err = qtx.UpdateUserDetailByUserId(ctx, profileSqlc.UpdateUserDetailByUserIdParams{
		UserID:          sql.NullInt64{Int64: props.UserId, Valid: true},
		HidePhoneNumber: sql.NullBool{Bool: props.HidePhoneNumber, Valid: true},
		PhoneNumber:     sql.NullString{String: props.PhoneNumber, Valid: true},
		Gender:          sql.NullString{String: props.Gender, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("could not update user detail: %w", err)
	}

	// insert to skills table (if not exist)
	if err := qtx.BatchInsertSkills(ctx, props.MainSkills); err != nil {
		return fmt.Errorf("could not batch insert skills: %w", err)
	}

	// change all main skills to false
	if err := qtx.UpdateUserMainSkillToFalse(ctx, props.UserId); err != nil {
		return fmt.Errorf("could not update user main skills to false: %w", err)
	}

	// insert user main skills
	err = qtx.BatchInsertUserSkills(ctx, profileSqlc.BatchInsertUserSkillsParams{
		UserID:      props.UserId,
		Names:       props.MainSkills,
		IsMainSkill: true,
	})
	if err != nil {
		return fmt.Errorf("could not batch insert user skills: %w", err)
	}

	// update or insert user social links
	for _, v := range props.SocialLinks {
		err := qtx.UpsertUserSocialLink(ctx, profileSqlc.UpsertUserSocialLinkParams{
			UserID: sql.NullInt64{Int64: props.UserId, Valid: true},
			Name:   v.Name,
			Url:    sql.NullString{String: v.URL, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("could not upsert user social links: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit edit profile transaction: %w", err)
	}

	return nil
}

func (r *ProfileRepository) UpdateAboutMe(userId int64, aboutMe string) error {
	arg := profileSqlc.UpdateUserDetailAboutParams{
		UserID: userId,
		About:  aboutMe,
	}

	if err := r.query.UpdateUserDetailAbout(context.Background(), arg); err != nil {
		return err
	}

	return nil
}
