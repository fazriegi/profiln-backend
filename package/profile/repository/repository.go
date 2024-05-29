package profile

import (
	"context"
	"database/sql"
	"fmt"
	"profiln-be/model"
	profileSqlc "profiln-be/package/profile/repository/sqlc"
	"time"
)

type IProfileRepository interface {
	InsertUserDetailAbout(arg profileSqlc.InsertUserDetailAboutParams) (profileSqlc.UserDetail, error)
	InsertCertificate(arg profileSqlc.InsertCertificateParams) (profileSqlc.Certificate, error)
	InsertUserSkill(arg profileSqlc.InsertUserSkillParams) (profileSqlc.UserSkill, error)
	InsertSkill(name string) (profileSqlc.Skill, error)
	InsertWorkExperience(arg profileSqlc.InsertWorkExperienceParams) (profileSqlc.WorkExperience, error)
	InsertUserAvatar(arg profileSqlc.InsertUserAvatarParams) error
	GetUserById(id int64) (profileSqlc.User, error)
	UpdateUserDetailAbout(arg profileSqlc.UpdateUserDetailAboutParams) error
	GetSkills(offset, limit int32) ([]profileSqlc.GetSkillsRow, int64, error)
	UpdateProfile(avatar_url string, props *model.UpdateProfileRequest) error
	UpdateAboutMe(userId int64, aboutMe string) error
	UpdateUserCertificate(userId int64, props *model.UpdateCertificate) error
	GetUserAvatarById(id int64) (string, error)
	UpdateUserInformation(props *model.UpdateUserInformation) error
	UpdateUserEducation(props *model.UpdateEducationRequest) error
	GetEducationById(id int64) (profileSqlc.Education, error)
	GetUserEducationFileURLs(educationId int64) ([]string, error)
	GetWorkExperienceById(id int64) (profileSqlc.WorkExperience, error)
	UpdateUserWorkExperience(props *model.UpdateWorkExperience) error
	GetWorkExperienceFileURLs(workExperienceId int64) ([]string, error)
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

func (r *ProfileRepository) GetUserById(id int64) (profileSqlc.User, error) {
	user, err := r.query.GetUserById(context.Background(), id)

	if err != nil {
		return profileSqlc.User{}, err
	}

	return user, nil
}

func (r *ProfileRepository) GetSkills(offset, limit int32) ([]profileSqlc.GetSkillsRow, int64, error) {
	arg := profileSqlc.GetSkillsParams{
		Offset: offset,
		Limit:  limit,
	}

	skills, err := r.query.GetSkills(context.Background(), arg)
	if err != nil {
		return []profileSqlc.GetSkillsRow{}, 0, err
	}

	var count int64
	if len(skills) > 0 {
		count = skills[0].TotalRows
	}

	return skills, count, nil
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
	_, err = qtx.BatchInsertUserMainSkills(ctx, profileSqlc.BatchInsertUserMainSkillsParams{
		UserID: props.UserId,
		Names:  props.MainSkills,
	})
	if err != nil {
		return fmt.Errorf("could not batch insert user skills: %w", err)
	}

	// update or insert user social links
	for _, v := range props.SocialLinks {
		err := qtx.UpsertUserSocialLink(ctx, profileSqlc.UpsertUserSocialLinkParams{
			UserID:   sql.NullInt64{Int64: props.UserId, Valid: true},
			Platform: sql.NullString{String: v.Platform, Valid: true},
			Url:      sql.NullString{String: v.URL, Valid: true},
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

func (r *ProfileRepository) UpdateUserCertificate(userId int64, props *model.UpdateCertificate) error {
	var (
		issueDate      time.Time
		expirationDate time.Time
		err            error
	)

	ctx := context.Background()
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := r.query.WithTx(tx)

	if props.IssuingOrganization.ID < 1 {
		createdIssuingOrganization, err :=
			qtx.InsertIssuingOrganization(ctx, props.IssuingOrganization.Name)

		if err != nil {
			return err
		}

		props.IssuingOrganization.ID = createdIssuingOrganization.ID
	}

	issueDate, err = time.Parse("2006-01-02", props.IssueDate)
	if err != nil {
		return fmt.Errorf("error parsing expiration date: %w", err)
	}

	if props.ExpirationDate != "" {
		expirationDate, err = time.Parse("2006-01-02", props.ExpirationDate)
		if err != nil {
			return fmt.Errorf("error parsing expiration date: %w", err)
		}
	}

	arg := profileSqlc.UpdateUserCertificateParams{
		Name:                  props.Name,
		IssuingOrganizationID: props.IssuingOrganization.ID,
		IssueDate:             issueDate,
		ExpirationDate:        expirationDate,
		CredentialID:          props.CredentialID,
		Url:                   props.Url,
		ID:                    props.ID,
		UserID:                userId,
	}

	_, err = qtx.UpdateUserCertificate(ctx, arg)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

func (r *ProfileRepository) GetUserAvatarById(id int64) (string, error) {
	avatarUrl, err := r.query.GetUserAvatarById(context.Background(), id)
	if err != nil {
		return "", err
	}

	return avatarUrl.String, nil
}

func (r *ProfileRepository) UpdateUserInformation(props *model.UpdateUserInformation) error {
	ctx := context.Background()
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := r.query.WithTx(tx)

	currentUserDetail, err := qtx.GetUserDetail(ctx, props.UserId)
	if err != nil {
		return fmt.Errorf("could not get user detail: %w", err)
	}

	updateUserDetailArg := profileSqlc.UpdateUserDetailParams{
		UserID:          sql.NullInt64{Int64: props.UserId, Valid: true},
		PhoneNumber:     sql.NullString{String: currentUserDetail.PhoneNumber.String, Valid: true},
		Gender:          sql.NullString{String: currentUserDetail.Gender.String, Valid: true},
		Location:        sql.NullString{String: props.Location, Valid: true},
		PortfolioUrl:    sql.NullString{String: props.PortfolioUrl, Valid: true},
		About:           sql.NullString{String: currentUserDetail.About.String, Valid: true},
		HidePhoneNumber: sql.NullBool{Bool: currentUserDetail.HidePhoneNumber.Bool, Valid: true},
	}
	_, err = r.updateUserDetail(ctx, qtx, &updateUserDetailArg)
	if err != nil {
		return fmt.Errorf("could not update user detail: %w", err)
	}

	_, err = r.batchInsertUserSkills(ctx, qtx, props.UserId, props.Skills)
	if err != nil {
		return fmt.Errorf("could not batch insert user skills: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

func (r *ProfileRepository) UpdateUserEducation(props *model.UpdateEducationRequest) error {
	var (
		startDate  time.Time
		finishDate time.Time
		err        error
	)

	ctx := context.Background()
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := r.query.WithTx(tx)

	// Delete current user education skills by education id
	err = qtx.DeleteEducationFilesByEducationId(ctx, props.ID)
	if err != nil {
		return fmt.Errorf("could not delete user education files: %w", err)
	}

	// Delete current user education skills by education id
	_, err = qtx.DeleteEducationSkillsByEducation(ctx, props.ID)
	if err != nil {
		return fmt.Errorf("could not delete user education skills: %w", err)
	}

	// Batch insert new user skills
	userSkillIDs, err := r.batchInsertUserSkills(ctx, qtx, props.UserId, props.Skills)
	if err != nil {
		return fmt.Errorf("could not batch insert user skills: %w", err)
	}

	err = qtx.BatchInsertEducationSkills(ctx, profileSqlc.BatchInsertEducationSkillsParams{
		EducationID: props.ID,
		UserSkillID: userSkillIDs,
	})
	if err != nil {
		return fmt.Errorf("could not batch insert user education skills: %w", err)
	}

	// If the school doesn't exist, insert the school first
	if props.School.ID < 1 {
		school, err := qtx.InsertSchool(ctx, props.School.Name)
		if err != nil {
			return fmt.Errorf("could not insert school: %w", err)
		}

		props.School.ID = school.ID
	}

	startDate, err = time.Parse("2006-01-02", props.StartDate)
	if err != nil {
		return fmt.Errorf("error parsing expiration date: %w", err)
	}

	if props.FinishDate != "" {
		finishDate, err = time.Parse("2006-01-02", props.FinishDate)
		if err != nil {
			return fmt.Errorf("error parsing expiration date: %w", err)
		}
	}

	arg := profileSqlc.UpdateUserEducationParams{
		ID:           props.ID,
		SchoolID:     sql.NullInt64{Int64: props.School.ID, Valid: true},
		Degree:       sql.NullString{String: props.Degree, Valid: true},
		FieldOfStudy: sql.NullString{String: props.FieldOfStudy, Valid: true},
		Gpa:          sql.NullString{String: props.GPA, Valid: true},
		StartDate:    sql.NullTime{Time: startDate, Valid: !startDate.IsZero()},
		FinishDate:   sql.NullTime{Time: finishDate, Valid: !finishDate.IsZero()},
		Description:  sql.NullString{String: props.Description, Valid: true},
	}
	_, err = qtx.UpdateUserEducation(ctx, arg)
	if err != nil {
		return fmt.Errorf("could not update user education: %w", err)
	}

	if _, err = r.batchInsertEducationFiles(ctx, qtx, props.ID, props.FileURLs); err != nil {
		return fmt.Errorf("could not batch insert education files: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

func (r *ProfileRepository) GetEducationById(id int64) (profileSqlc.Education, error) {
	data, err := r.query.GetEducationById(context.Background(), id)
	if err != nil {
		return profileSqlc.Education{}, err
	}

	return data, nil
}

func (r *ProfileRepository) GetUserEducationFileURLs(educationId int64) ([]string, error) {
	data, err := r.query.GetUserEducationFileURLs(context.Background(), educationId)
	if err != nil {
		return nil, err
	}

	urls := make([]string, len(data))
	for i, v := range data {
		urls[i] = v.String
	}

	return urls, nil
}

func (r *ProfileRepository) GetWorkExperienceById(id int64) (profileSqlc.WorkExperience, error) {
	data, err := r.query.GetWorkExperienceById(context.Background(), id)
	if err != nil {
		return profileSqlc.WorkExperience{}, err
	}

	return data, nil
}

func (r *ProfileRepository) GetWorkExperienceFileURLs(workExperienceId int64) ([]string, error) {
	data, err := r.query.GetWorkExperienceFileURLs(context.Background(), workExperienceId)
	if err != nil {
		return nil, err
	}

	urls := make([]string, len(data))
	for i, v := range data {
		urls[i] = v.String
	}

	return urls, nil
}

func (r *ProfileRepository) UpdateUserWorkExperience(props *model.UpdateWorkExperience) error {
	var (
		startDate  time.Time
		finishDate time.Time
		err        error
	)

	ctx := context.Background()
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := r.query.WithTx(tx)

	// Delete current user work experience skills by work experience id
	err = qtx.DeleteWorkExperienceFilesByWorkExperienceId(ctx, props.ID)
	if err != nil {
		return fmt.Errorf("could not delete user work experience files: %w", err)
	}

	// Delete current user work experience skills by work experience id
	_, err = qtx.DeleteWorkExperienceSkillsByWorkExperience(ctx, props.ID)
	if err != nil {
		return fmt.Errorf("could not delete user work experience skills: %w", err)
	}

	// Batch insert new user skills
	userSkillIDs, err := r.batchInsertUserSkills(ctx, qtx, props.UserId, props.Skills)
	if err != nil {
		return fmt.Errorf("could not batch insert user skills: %w", err)
	}

	err = qtx.BatchInsertWorkExperienceSkills(ctx, profileSqlc.BatchInsertWorkExperienceSkillsParams{
		WorkExperienceID: props.ID,
		UserSkillID:      userSkillIDs,
	})
	if err != nil {
		return fmt.Errorf("could not batch insert user work experience skills: %w", err)
	}

	// If the company doesn't exist, insert the company first
	if props.Company.ID < 1 {
		company, err := qtx.InsertCompany(ctx, props.Company.Name)
		if err != nil {
			return fmt.Errorf("could not insert company: %w", err)
		}

		props.Company.ID = company.ID
	}

	startDate, err = time.Parse("2006-01-02", props.StartDate)
	if err != nil {
		return fmt.Errorf("error parsing expiration date: %w", err)
	}

	if props.FinishDate != "" {
		finishDate, err = time.Parse("2006-01-02", props.FinishDate)
		if err != nil {
			return fmt.Errorf("error parsing expiration date: %w", err)
		}
	}

	updateUserWorkExperienceArg := profileSqlc.UpdateUserWorkExperienceParams{
		ID:             props.ID,
		CompanyID:      sql.NullInt64{Int64: props.Company.ID, Valid: true},
		JobTitle:       sql.NullString{String: props.JobTitle, Valid: true},
		EmploymentType: sql.NullString{String: props.EmploymentType, Valid: true},
		Location:       sql.NullString{String: props.Location, Valid: true},
		LocationType:   sql.NullString{String: props.LocationType, Valid: true},
		StartDate:      sql.NullTime{Time: startDate, Valid: !startDate.IsZero()},
		FinishDate:     sql.NullTime{Time: finishDate, Valid: !finishDate.IsZero()},
		Description:    sql.NullString{String: props.Description, Valid: true},
	}
	_, err = qtx.UpdateUserWorkExperience(ctx, updateUserWorkExperienceArg)
	if err != nil {
		return fmt.Errorf("could not update user work experience: %w", err)
	}

	_, err = qtx.BatchInsertWorkExperienceFiles(ctx, profileSqlc.BatchInsertWorkExperienceFilesParams{
		WorkExperienceID: props.ID,
		Url:              props.FileURLs,
	})
	if err != nil {
		return fmt.Errorf("could not batch insert user work experience files: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

func (r *ProfileRepository) updateUserDetail(ctx context.Context, qtx *profileSqlc.Queries, props *profileSqlc.UpdateUserDetailParams) (profileSqlc.UpdateUserDetailRow, error) {
	data, err := qtx.UpdateUserDetail(ctx, *props)
	if err != nil {
		return data, err
	}

	return data, nil
}

// Batch insert to skills and user skills table (if not exist)
func (r *ProfileRepository) batchInsertUserSkills(ctx context.Context, qtx *profileSqlc.Queries, userId int64, skills []string) ([]int64, error) {
	var (
		userSkillIDs []int64
		err          error
	)
	// Insert skills
	if err := qtx.BatchInsertSkills(ctx, skills); err != nil {
		return nil, fmt.Errorf("could not batch insert skills: %w", err)
	}

	// Insert user skills
	_, err = qtx.BatchInsertUserSkills(ctx, profileSqlc.BatchInsertUserSkillsParams{
		UserID:      userId,
		Names:       skills,
		IsMainSkill: false,
	})
	if err != nil {
		return nil, fmt.Errorf("could not batch insert user skills: %w", err)
	}

	// If no new user skills were inserted, get the existing user skills
	userSkillIDs, err = qtx.GetUserSkillIDsByName(ctx, skills)
	if err != nil {
		return nil, fmt.Errorf("could not get user skills: %w", err)
	}

	return userSkillIDs, nil
}

func (r *ProfileRepository) batchInsertEducationFiles(ctx context.Context, qtx *profileSqlc.Queries, educationId int64, url []string) ([]profileSqlc.EducationFile, error) {
	arg := profileSqlc.BatchInsertEducationFilesParams{
		EducationID: educationId,
		Url:         url,
	}
	educationFiles, err := qtx.BatchInsertEducationFiles(ctx, arg)
	if err != nil {
		return nil, err
	}

	return educationFiles, nil
}
