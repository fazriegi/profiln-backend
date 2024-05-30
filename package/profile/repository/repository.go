package profile

import (
	"context"
	"database/sql"
	"fmt"
	db "profiln-be/db/sqlc"
	"profiln-be/model"
	"strings"
	"sync"
	"time"
)

type IProfileRepository interface {
	InsertUserDetailAbout(arg db.InsertUserDetailAboutParams) (db.UserDetail, error)
	InsertCertificate(arg db.InsertCertificateParams) (db.Certificate, error)
	InsertUserSkill(arg db.InsertUserSkillParams) (db.UserSkill, error)
	InsertSkill(name string) (db.Skill, error)
	InsertWorkExperience(arg db.InsertWorkExperienceParams) (db.WorkExperience, error)
	InsertUserAvatar(arg db.InsertUserAvatarParams) error
	GetUserById(id int64) (model.User, error)
	UpdateUserDetailAbout(arg db.UpdateUserDetailAboutParams) error
	UpdateProfile(avatar_url string, props *model.UpdateProfileRequest) error
	UpdateAboutMe(userId int64, aboutMe string) error
	UpdateUserCertificate(userId int64, props *model.UpdateCertificate) error
	AddUserOpenToWork(props *model.OpenToWork) error
	GetUserAvatarById(id int64) (string, error)
	UpdateUserInformation(props *model.UpdateUserInformation) error
	UpdateUserEducation(props *model.UpdateEducationRequest) error
	GetEducationById(id int64) (db.Education, error)
	GetUserEducationFileURLs(educationId int64) ([]string, error)
	GetWorkExperienceById(id int64) (db.WorkExperience, error)
	UpdateUserWorkExperience(props *model.UpdateWorkExperience) error
	GetWorkExperienceFileURLs(workExperienceId int64) ([]string, error)
	GetUserProfile(userId int64) (model.UserProfile, error)
	GetWorkExperiencesByUserId(userId int64, offset, limit int32) ([]model.WorkExperience, int64, error)
	GetEducationsByUserId(userId int64, offset, limit int32) ([]model.Education, int64, error)
	GetCertificatesByUserId(userId int64, offset, limit int32) ([]model.Certificate, int64, error)
	GetFollowedUsersByUserId(userId int64, offset, limit int32) ([]model.User, int64, error)
	DeleteUserOpenToWork(userId int64) error
	DeleteUserWorkExperienceById(userId, workExperienceId int64) error
}

type ProfileRepository struct {
	dbConn *sql.DB
	query  *db.Queries
}

func NewProfileRepository(dbConn *sql.DB) IProfileRepository {
	return &ProfileRepository{
		dbConn: dbConn,
		query:  db.New(dbConn),
	}
}

func (r *ProfileRepository) InsertUserAvatar(arg db.InsertUserAvatarParams) error {
	err := r.query.InsertUserAvatar(context.Background(), arg)

	if err != nil {
		return err
	}

	return nil
}

func (r *ProfileRepository) InsertUserDetailAbout(arg db.InsertUserDetailAboutParams) (db.UserDetail, error) {
	userAbout, err := r.query.InsertUserDetailAbout(context.Background(), arg)

	if err != nil {
		return db.UserDetail{}, err
	}

	return userAbout, nil
}

func (r *ProfileRepository) UpdateUserDetailAbout(arg db.UpdateUserDetailAboutParams) error {
	err := r.query.UpdateUserDetailAbout(context.Background(), arg)

	if err != nil {
		return err
	}

	return nil
}

func (r *ProfileRepository) InsertWorkExperience(arg db.InsertWorkExperienceParams) (db.WorkExperience, error) {
	workExperience, err := r.query.InsertWorkExperience(context.Background(), arg)

	if err != nil {
		return db.WorkExperience{}, err
	}

	return workExperience, nil
}

func (r *ProfileRepository) InsertCertificate(arg db.InsertCertificateParams) (db.Certificate, error) {
	certificate, err := r.query.InsertCertificate(context.Background(), arg)

	if err != nil {
		return db.Certificate{}, err
	}

	return certificate, nil
}

func (r *ProfileRepository) InsertUserSkill(arg db.InsertUserSkillParams) (db.UserSkill, error) {
	userSkill, err := r.query.InsertUserSkill(context.Background(), arg)

	if err != nil {
		return db.UserSkill{}, err
	}

	return userSkill, nil
}

func (r *ProfileRepository) InsertSkill(name string) (db.Skill, error) {
	skill, err := r.query.InsertSkill(context.Background(), name)

	if err != nil {
		return db.Skill{}, err
	}

	return skill, nil
}

func (r *ProfileRepository) GetUserById(id int64) (model.User, error) {
	user, err := r.query.GetUserById(context.Background(), id)

	if err != nil {
		return model.User{}, err
	}

	data := model.User{
		ID:         user.ID,
		Fullname:   user.FullName,
		AvatarUrl:  user.AvatarUrl.String,
		Bio:        user.Bio.String,
		OpenToWork: user.OpenToWork.Bool,
	}

	return data, nil
}

func (r *ProfileRepository) UpdateProfile(avatar_url string, props *model.UpdateProfileRequest) error {
	ctx := context.Background()
	tx, err := r.dbConn.Begin()
	if err != nil {
		return fmt.Errorf("could not begin edit profile transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := r.query.WithTx(tx)

	// update users table
	_, err = qtx.UpdateUser(ctx, db.UpdateUserParams{
		ID:        props.UserId,
		FullName:  props.Fullname,
		AvatarUrl: sql.NullString{String: avatar_url, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("could not update user: %w", err)
	}

	// update user details table
	_, err = qtx.UpdateUserDetailByUserId(ctx, db.UpdateUserDetailByUserIdParams{
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
	_, err = qtx.BatchInsertUserMainSkills(ctx, db.BatchInsertUserMainSkillsParams{
		UserID: props.UserId,
		Names:  props.MainSkills,
	})
	if err != nil {
		return fmt.Errorf("could not batch insert user skills: %w", err)
	}

	// update or insert user social links
	for _, v := range props.SocialLinks {
		err := qtx.UpsertUserSocialLink(ctx, db.UpsertUserSocialLinkParams{
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
	arg := db.UpdateUserDetailAboutParams{
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
	tx, err := r.dbConn.Begin()
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

	arg := db.UpdateUserCertificateParams{
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
	tx, err := r.dbConn.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := r.query.WithTx(tx)

	currentUserDetail, err := qtx.GetUserDetail(ctx, props.UserId)
	if err != nil {
		return fmt.Errorf("could not get user detail: %w", err)
	}

	updateUserDetailArg := db.UpdateUserDetailParams{
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
	tx, err := r.dbConn.Begin()
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

	err = qtx.BatchInsertEducationSkills(ctx, db.BatchInsertEducationSkillsParams{
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

	arg := db.UpdateUserEducationParams{
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

func (r *ProfileRepository) GetEducationById(id int64) (db.Education, error) {
	data, err := r.query.GetEducationById(context.Background(), id)
	if err != nil {
		return db.Education{}, err
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

func (r *ProfileRepository) GetWorkExperienceById(id int64) (db.WorkExperience, error) {
	data, err := r.query.GetWorkExperienceById(context.Background(), id)
	if err != nil {
		return db.WorkExperience{}, err
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
	tx, err := r.dbConn.Begin()
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

	err = qtx.BatchInsertWorkExperienceSkills(ctx, db.BatchInsertWorkExperienceSkillsParams{
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

	updateUserWorkExperienceArg := db.UpdateUserWorkExperienceParams{
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

	_, err = qtx.BatchInsertWorkExperienceFiles(ctx, db.BatchInsertWorkExperienceFilesParams{
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

func (r *ProfileRepository) GetUserProfile(userId int64) (model.UserProfile, error) {
	var (
		wg sync.WaitGroup
	)

	userChan := make(chan db.GetUserProfileRow, 1)
	socialLinksChan := make(chan []model.SocialLinks, 1)
	userSkillsChan := make(chan model.UserSkills, 1)
	errChan := make(chan error, 3)

	wg.Add(1)
	go func(userId int64) {
		defer wg.Done()
		data, err := r.query.GetUserProfile(context.Background(), userId)
		if err != nil {
			errChan <- err
			close(userChan)
			return
		}

		userChan <- data
		close(userChan)
	}(userId)

	wg.Add(1)
	go func(userId int64) {
		defer wg.Done()
		data, err := r.getUserSocialLinks(userId)
		if err != nil {
			errChan <- err
			close(socialLinksChan)
			return
		}

		socialLinksChan <- data
		close(socialLinksChan)
	}(userId)

	wg.Add(1)
	go func(userId int64) {
		defer wg.Done()
		data, err := r.getUserSkills(userId)
		if err != nil {
			errChan <- err
			close(userSkillsChan)
			return
		}

		userSkillsChan <- data
		close(userSkillsChan)
	}(userId)

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return model.UserProfile{}, err
		}
	}

	user := <-userChan
	socialLinks := <-socialLinksChan
	userSkills := <-userSkillsChan

	data := model.UserProfile{
		User: model.User{
			ID:         user.ID,
			AvatarUrl:  user.AvatarUrl.String,
			Fullname:   user.FullName,
			Bio:        user.Bio.String,
			OpenToWork: user.OpenToWork.Bool,
		},
		FollowingCount:  int64(user.FollowingsCount.Int32),
		SocialLinks:     socialLinks,
		Skills:          userSkills,
		Location:        user.Location.String,
		WebPortfolioUrl: user.PortfolioUrl.String,
		About:           user.About.String,
	}

	return data, nil
}

func (r *ProfileRepository) getUserSocialLinks(userId int64) ([]model.SocialLinks, error) {
	socialLinks, err := r.query.GetUserSocialLinks(context.Background(), userId)
	if err != nil {
		return nil, err
	}

	data := make([]model.SocialLinks, len(socialLinks))

	for i, v := range socialLinks {
		data[i] = model.SocialLinks{
			Platform: v.Platform.String,
			URL:      v.Url.String,
		}
	}

	return data, nil
}

func (r *ProfileRepository) getUserSkills(userId int64) (model.UserSkills, error) {
	userSkills, err := r.query.GetUserSkills(context.Background(), userId)
	if err != nil {
		return model.UserSkills{}, err
	}

	var (
		mainSkills  []string
		otherSkills []string
	)

	for _, userSkill := range userSkills {
		if userSkill.MainSkill.Bool {
			mainSkills = append(mainSkills, userSkill.Name.String)
			continue
		}

		otherSkills = append(otherSkills, userSkill.Name.String)
	}

	data := model.UserSkills{
		MainSkills:  mainSkills,
		OtherSkills: otherSkills,
	}

	return data, nil
}

func (r *ProfileRepository) updateUserDetail(ctx context.Context, qtx *db.Queries, props *db.UpdateUserDetailParams) (db.UpdateUserDetailRow, error) {
	data, err := qtx.UpdateUserDetail(ctx, *props)
	if err != nil {
		return data, err
	}

	return data, nil
}

// Batch insert to skills and user skills table (if not exist)
func (r *ProfileRepository) batchInsertUserSkills(ctx context.Context, qtx *db.Queries, userId int64, skills []string) ([]int64, error) {
	var (
		userSkillIDs []int64
		err          error
	)
	// Insert skills
	if err := qtx.BatchInsertSkills(ctx, skills); err != nil {
		return nil, fmt.Errorf("could not batch insert skills: %w", err)
	}

	// Insert user skills
	_, err = qtx.BatchInsertUserSkills(ctx, db.BatchInsertUserSkillsParams{
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

func (r *ProfileRepository) batchInsertEducationFiles(ctx context.Context, qtx *db.Queries, educationId int64, url []string) ([]db.EducationFile, error) {
	arg := db.BatchInsertEducationFilesParams{
		EducationID: educationId,
		Url:         url,
	}
	educationFiles, err := qtx.BatchInsertEducationFiles(ctx, arg)
	if err != nil {
		return nil, err
	}

	return educationFiles, nil
}

func (r *ProfileRepository) GetWorkExperiencesByUserId(userId int64, offset, limit int32) ([]model.WorkExperience, int64, error) {
	arg := db.GetWorkExperiencesByUserIdParams{
		Offset: offset,
		Limit:  limit,
		UserID: userId,
	}
	workExperiences, err := r.query.GetWorkExperiencesByUserId(context.Background(), arg)
	if err != nil {
		return nil, 0, err
	}

	data := make([]model.WorkExperience, len(workExperiences))
	finishDate := "now"

	var count int64
	if len(workExperiences) > 0 {
		count = workExperiences[0].TotalRows
	}

	for i, workExperience := range workExperiences {
		var fileUrls []string
		var skills []string
		startDate := workExperience.StartDate.Time.Format("2006-01-02")

		if !workExperience.FinishDate.Time.IsZero() {
			finishDate = workExperience.FinishDate.Time.Format("2006-01-02")
		}

		fileUrlsString := strings.Trim(string(workExperience.FileUrls.([]uint8)), "{}")
		if fileUrlsString != "NULL" {
			fileUrls = strings.Split(fileUrlsString, ",")
		}

		skillsString := strings.Trim(string(workExperience.Skills.([]uint8)), "{}")
		skillsString = strings.ReplaceAll(skillsString, "\"", "")
		if skillsString != "NULL" {
			skills = strings.Split(skillsString, ",")
		}

		data[i] = model.WorkExperience{
			ID:       workExperience.ID,
			JobTitle: workExperience.JobTitle.String,
			Company: model.Company{
				ID:   workExperience.CompanyID.Int64,
				Name: workExperience.CompanyName.String,
			},
			EmploymentType: workExperience.EmploymentType.String,
			Location:       workExperience.Location.String,
			LocationType:   workExperience.LocationType.String,
			StartDate:      startDate,
			FinishDate:     finishDate,
			Description:    workExperience.Description.String,
			FileURLs:       fileUrls,
			Skills:         skills,
		}
	}

	return data, count, nil
}

func (r *ProfileRepository) GetEducationsByUserId(userId int64, offset, limit int32) ([]model.Education, int64, error) {
	arg := db.GetEducationsByUserIdParams{
		Offset: offset,
		Limit:  limit,
		UserID: userId,
	}
	educations, err := r.query.GetEducationsByUserId(context.Background(), arg)
	if err != nil {
		return nil, 0, err
	}

	data := make([]model.Education, len(educations))
	finishDate := "now"

	var count int64
	if len(educations) > 0 {
		count = educations[0].TotalRows
	}

	for i, education := range educations {
		var fileUrls []string
		var skills []string
		startDate := education.StartDate.Time.Format("2006-01-02")

		if !education.FinishDate.Time.IsZero() {
			finishDate = education.FinishDate.Time.Format("2006-01-02")
		}

		fileUrlsString := strings.Trim(string(education.FileUrls.([]uint8)), "{}")
		if fileUrlsString != "NULL" {
			fileUrls = strings.Split(fileUrlsString, ",")
		}

		skillsString := strings.Trim(string(education.Skills.([]uint8)), "{}")
		skillsString = strings.ReplaceAll(skillsString, "\"", "")
		if skillsString != "NULL" {
			skills = strings.Split(skillsString, ",")
		}

		data[i] = model.Education{
			ID: education.ID,
			School: model.School{
				ID:   education.SchoolID.Int64,
				Name: education.SchoolName.String,
			},
			Degree:       education.Degree.String,
			FieldOfStudy: education.FieldOfStudy.String,
			StartDate:    startDate,
			FinishDate:   finishDate,
			GPA:          education.Gpa.String,
			Description:  education.Description.String,
			FileURLs:     fileUrls,
			Skills:       skills,
		}
	}

	return data, count, nil
}

func (r *ProfileRepository) GetCertificatesByUserId(userId int64, offset, limit int32) ([]model.Certificate, int64, error) {
	arg := db.GetCertificatesByUserIdParams{
		Offset: offset,
		Limit:  limit,
		UserID: userId,
	}
	certificates, err := r.query.GetCertificatesByUserId(context.Background(), arg)
	if err != nil {
		return nil, 0, err
	}

	data := make([]model.Certificate, len(certificates))
	expirationDate := ""

	var count int64
	if len(certificates) > 0 {
		count = certificates[0].TotalRows
	}

	for i, certificate := range certificates {
		issueDate := certificate.IssueDate.Time.Format("2006-01-02")

		if !certificate.ExpirationDate.Time.IsZero() {
			expirationDate = certificate.ExpirationDate.Time.Format("2006-01-02")
		}

		data[i] = model.Certificate{
			ID:             certificate.ID,
			Organization:   certificate.IssuingOrganizationName.String,
			Name:           certificate.Name.String,
			IssueDate:      issueDate,
			ExpirationDate: expirationDate,
			CredentialID:   certificate.CredentialID.String,
			Url:            certificate.Url.String,
		}
	}

	return data, count, nil
}

func (r *ProfileRepository) GetFollowedUsersByUserId(userId int64, offset, limit int32) ([]model.User, int64, error) {
	arg := db.GetFollowedUsersByUserIdParams{
		Offset: offset,
		Limit:  limit,
		UserID: userId,
	}
	followedUsers, err := r.query.GetFollowedUsersByUserId(context.Background(), arg)
	if err != nil {
		return nil, 0, err
	}

	data := make([]model.User, len(followedUsers))

	var count int64
	if len(followedUsers) > 0 {
		count = followedUsers[0].TotalRows
	}

	for i, followedUser := range followedUsers {
		data[i] = model.User{
			ID:         followedUser.ID.Int64,
			Fullname:   followedUser.FullName.String,
			AvatarUrl:  followedUser.AvatarUrl.String,
			Bio:        followedUser.Bio.String,
			OpenToWork: followedUser.OpenToWork.Bool,
		}
	}

	return data, count, nil
}

func (r *ProfileRepository) AddUserOpenToWork(props *model.OpenToWork) error {
	ctx := context.Background()
	tx, err := r.dbConn.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := r.query.WithTx(tx)

	err = r.deleteOpenToWorkDataByUserId(ctx, qtx, props.UserId)
	if err != nil {
		return err
	}

	_, err = qtx.UpdateUserOpenToWork(ctx, db.UpdateUserOpenToWorkParams{
		UserID:     props.UserId,
		OpenToWork: props.OpenToWork,
	})
	if err != nil {
		return fmt.Errorf("could not update user open to work: %w", err)
	}

	var jobPositionIds []int64
	for i, jobPosition := range props.JobPositions {
		if jobPosition.ID < 1 {
			createdJobPosition, err := qtx.InsertJobPosition(ctx, sql.NullString{String: jobPosition.Name, Valid: true})
			if err != nil {
				return fmt.Errorf("could not insert job position: %w", err)
			}
			props.JobPositions[i].ID = createdJobPosition.ID
		}

		jobPositionIds = append(jobPositionIds, props.JobPositions[i].ID)
	}

	err = qtx.BatchInsertUserJobInterests(ctx, db.BatchInsertUserJobInterestsParams{
		UserID:        props.UserId,
		JobPositionID: jobPositionIds,
	})
	if err != nil {
		return fmt.Errorf("could not batch insert user job interests: %w", err)
	}

	err = qtx.BatchInsertUserLocationTypeInterests(ctx, db.BatchInsertUserLocationTypeInterestsParams{
		UserID:       props.UserId,
		LocationType: props.LocationTypes,
	})
	if err != nil {
		return fmt.Errorf("could not batch insert user location type interests: %w", err)
	}

	err = qtx.BatchInsertUserEmploymentTypeInterests(ctx, db.BatchInsertUserEmploymentTypeInterestsParams{
		UserID:         props.UserId,
		EmploymentType: props.EmploymentTypes,
	})
	if err != nil {
		return fmt.Errorf("could not batch insert user employment type interests: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

func (r *ProfileRepository) DeleteUserOpenToWork(userId int64) error {
	ctx := context.Background()
	tx, err := r.dbConn.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := r.query.WithTx(tx)

	_, err = qtx.UpdateUserOpenToWork(ctx, db.UpdateUserOpenToWorkParams{
		UserID:     userId,
		OpenToWork: false,
	})
	if err != nil {
		return fmt.Errorf("could not update user open to work: %w", err)
	}

	err = r.deleteOpenToWorkDataByUserId(ctx, qtx, userId)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

func (r *ProfileRepository) deleteOpenToWorkDataByUserId(ctx context.Context, qtx *db.Queries, userId int64) error {
	errChan := make(chan error, 3)
	var wg sync.WaitGroup

	wg.Add(3)
	go func() {
		defer wg.Done()
		if err := qtx.BatchDeleteUserJobInterests(ctx, userId); err != nil {
			errChan <- fmt.Errorf("could not delete user job interests: %w", err)
		}
	}()
	go func() {
		defer wg.Done()
		if err := qtx.BatchDeleteUserLocationTypeInterests(ctx, userId); err != nil {
			errChan <- fmt.Errorf("could not delete user location type interests: %w", err)
		}
	}()
	go func() {
		defer wg.Done()
		if err := qtx.BatchDeleteUserEmploymentTypeInterests(ctx, userId); err != nil {
			errChan <- fmt.Errorf("could not delete user employment type interests: %w", err)
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		return err
	}

	return nil
}

func (r *ProfileRepository) DeleteUserWorkExperienceById(userId, workExperienceId int64) error {
	ctx := context.Background()
	tx, err := r.dbConn.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := r.query.WithTx(tx)

	err = qtx.DeleteWorkExperienceFilesByWorkExperienceId(ctx, workExperienceId)
	if err != nil {
		return fmt.Errorf("could not delete work experience files: %w", err)
	}

	_, err = qtx.DeleteWorkExperienceSkillsByWorkExperience(ctx, workExperienceId)
	if err != nil {
		return fmt.Errorf("could not delete work experience skills: %w", err)
	}

	err = qtx.DeleteWorkExperienceById(ctx, db.DeleteWorkExperienceByIdParams{
		UserID: userId,
		ID:     workExperienceId,
	})
	if err != nil {
		return fmt.Errorf("could not delete work experience: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}
