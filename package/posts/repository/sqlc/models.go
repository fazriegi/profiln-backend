// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package posts

import (
	"database/sql"
)

type Certificate struct {
	ID                    int64
	UserID                sql.NullInt64
	Name                  sql.NullString
	IssuingOrganizationID sql.NullInt64
	IssueDate             sql.NullTime
	ExpirationDate        sql.NullTime
	CredentialID          sql.NullString
	Url                   sql.NullString
	CreatedAt             sql.NullTime
	UpdatedAt             sql.NullTime
}

type Company struct {
	ID   int64
	Name string
}

type Education struct {
	ID           int64
	UserID       sql.NullInt64
	SchoolID     sql.NullInt64
	Degree       sql.NullString
	FieldOfStudy sql.NullString
	Gpa          sql.NullString
	StartDate    sql.NullTime
	FinishDate   sql.NullTime
	Description  sql.NullString
	DocumentUrl  sql.NullString
	CreatedAt    sql.NullTime
	UpdatedAt    sql.NullTime
}

type EducationSkill struct {
	ID          int64
	EducationID sql.NullInt64
	SkillID     sql.NullInt64
}

type EmploymentType struct {
	ID   int16
	Name string
}

type Following struct {
	ID           int64
	UserID       sql.NullInt64
	FollowUserID sql.NullInt64
}

type IssuingOrganization struct {
	ID   int64
	Name string
}

type LocationType struct {
	ID   int16
	Name string
}

type Post struct {
	ID             int64
	UserID         sql.NullInt64
	Content        sql.NullString
	ImageUrl       sql.NullString
	LikeCount      sql.NullInt32
	CommentCount   sql.NullInt32
	RepostCount    sql.NullInt32
	Repost         sql.NullBool
	OriginalPostID sql.NullInt64
	CreatedAt      sql.NullTime
	UpdatedAt      sql.NullTime
}

type PostComment struct {
	ID           int64
	UserID       sql.NullInt64
	PostID       sql.NullInt64
	Content      sql.NullString
	ImageUrl     sql.NullString
	LikeCount    sql.NullInt32
	ReplyCount   sql.NullInt32
	IsPostAuthor sql.NullBool
	CreatedAt    sql.NullTime
	UpdatedAt    sql.NullTime
}

type PostCommentReply struct {
	ID            int64
	UserID        sql.NullInt64
	PostCommentID sql.NullInt64
	Content       sql.NullString
	ImageUrl      sql.NullString
	LikeCount     sql.NullInt32
	IsPostAuthor  sql.NullBool
	CreatedAt     sql.NullTime
	UpdatedAt     sql.NullTime
}

type ReportedPost struct {
	ID      int64
	UserID  sql.NullInt64
	PostID  sql.NullInt64
	Reason  sql.NullString
	Message sql.NullString
}

type School struct {
	ID   int64
	Name string
}

type Skill struct {
	ID   int64
	Name string
}

type SocialLink struct {
	ID   int16
	Name string
}

type User struct {
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
}

type UserDetail struct {
	ID              int64
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

type UserEmploymentTypeInterest struct {
	ID               int64
	UserID           sql.NullInt64
	EmploymentTypeID sql.NullInt16
}

type UserJobInterest struct {
	ID       int64
	UserID   sql.NullInt64
	JobTitle sql.NullString
}

type UserLocationTypeInterest struct {
	ID             int64
	UserID         sql.NullInt64
	LocationTypeID sql.NullInt16
}

type UserOtp struct {
	ID     int64
	UserID sql.NullInt64
	Otp    sql.NullString
}

type UserSkill struct {
	ID        int64
	UserID    sql.NullInt64
	SkillID   sql.NullInt64
	MainSkill sql.NullBool
}

type UserSocialLink struct {
	ID           int64
	UserID       sql.NullInt64
	SocialLinkID sql.NullInt16
	Url          sql.NullString
}

type WorkExperience struct {
	ID               int64
	UserID           sql.NullInt64
	JobTitle         sql.NullString
	CompanyID        sql.NullInt64
	EmploymentTypeID sql.NullInt16
	Location         sql.NullString
	LocationTypeID   sql.NullInt16
	StartDate        sql.NullTime
	FinishDate       sql.NullTime
	Description      sql.NullString
	CreatedAt        sql.NullTime
	UpdatedAt        sql.NullTime
}

type WorkExperienceSkill struct {
	ID               int64
	WorkExperienceID sql.NullInt64
	SkillID          sql.NullInt64
}
