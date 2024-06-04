package model

import (
	"time"
)

type SkillRequest struct {
	Name string `validate:"required"`
}

type UserDetailAboutRequest struct {
	About string `json:"about" validate:"required"`
	// UserID sql.NullInt64
}

type UserDetailRequest struct {
	PhoneNumber string `validate:"required"`
	Gender      string `validate:"required"`
}

type UpdateProfileRequest struct {
	UserId          int64         `json:"user_id" form:"user_id"`
	Fullname        string        `json:"fullname" form:"fullname" validate:"required"`
	HidePhoneNumber bool          `json:"hide_phone_number" form:"hide_phone_number" validate:"required,boolean"`
	MainSkills      []string      `json:"main_skills" form:"main_skills" validate:"required"`
	PhoneNumber     string        `json:"phone_number" form:"phone_number" validate:"required"`
	Gender          string        `json:"gender" form:"gender" validate:"required"`
	SocialLinks     []SocialLinks `json:"social_links" form:"social_links" validate:"required,isNotEmptyArray"`
}

type UpdateProfileResponse struct {
	UserId          int64         `json:"user_id"`
	Fullname        string        `json:"fullname"`
	AvatarUrl       string        `json:"avatar_url"`
	HidePhoneNumber bool          `json:"hide_phone_number"`
	MainSkills      []string      `json:"main_skills"`
	PhoneNumber     string        `json:"phone_number"`
	Gender          string        `json:"gender"`
	SocialLinks     []SocialLinks `json:"social_links"`
}

type SocialLinks struct {
	Platform string `json:"platform" form:"platform" validate:"required"`
	URL      string `json:"url" form:"url" validate:"required"`
}

type Certificate struct {
	ID                  int64               `json:"id"`
	UserId              int64               `json:"user_id"`
	Name                string              `json:"name" validate:"required"`
	IssuingOrganization IssuingOrganization `json:"issuing_organization" validate:"required"`
	IssueDate           string              `json:"issue_date" validate:"required"`
	ExpirationDate      string              `json:"expiration_date"`
	CredentialID        string              `json:"credential_id"`
	Url                 string              `json:"url"`
}

type UpdateUserInformation struct {
	UserId       int64    `json:"user_id"`
	Skills       []string `json:"skills"`
	Location     string   `json:"location" validate:"required"`
	PortfolioUrl string   `json:"portfolio_url"`
}

type Education struct {
	ID           int64    `json:"id"`
	UserId       int64    `json:"user_id"`
	School       School   `json:"school" form:"school" validate:"required"`
	Degree       string   `json:"degree" form:"degree" validate:"required"`
	FieldOfStudy string   `json:"field_of_study" form:"field_of_study" validate:"required"`
	StartDate    string   `json:"start_date" form:"start_date" validate:"required"`
	FinishDate   string   `json:"finish_date" form:"finish_date"`
	GPA          string   `json:"gpa" form:"gpa" validate:"required"`
	Description  string   `json:"description"  form:"description" validate:"required"`
	FileURLs     []string `json:"file_urls"`
	Skills       []string `json:"skills" form:"skills"`
}

type WorkExperience struct {
	ID             int64    `json:"id"`
	UserId         int64    `json:"user_id"`
	JobTitle       string   `json:"job_title" form:"job_title" validate:"required"`
	Company        Company  `json:"company" form:"company" validate:"required"`
	EmploymentType string   `json:"employment_type" form:"employment_type" validate:"required"`
	Location       string   `json:"location" form:"location" validate:"required"`
	LocationType   string   `json:"location_type" form:"location_type" validate:"required"`
	StartDate      string   `json:"start_date" form:"start_date" validate:"required"`
	FinishDate     string   `json:"finish_date" form:"finish_date"`
	Description    string   `json:"description"  form:"description" validate:"required"`
	FileURLs       []string `json:"file_urls"`
	Skills         []string `json:"skills" form:"skills"`
}

type OpenToWork struct {
	UserId          int64         `json:"user_id"`
	OpenToWork      bool          `json:"open_to_work"`
	JobPositions    []JobPosition `json:"job_positions" validate:"required,isNotEmptyArray"`
	LocationTypes   []string      `json:"location_types" validate:"required,isNotEmptyArray"`
	EmploymentTypes []string      `json:"employment_types" validate:"required,isNotEmptyArray"`
}

type AboutProfileResponse struct {
	ID        *int64     `json:"id"`
	UserID    *int64     `json:"user_id"`
	About     *string    `json:"about"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type InsertCertificateResponse struct {
	ID                    int64     `json:"id"`
	UserID                int64     `json:"user_id"`
	Name                  string    `json:"name"`
	IssuingOrganizationID int64     `json:"issuing_organization_id"`
	IssueDate             time.Time `json:"issue_date"`
	ExpirationDate        time.Time `json:"expiration_date"`
	CredentialID          string    `json:"credential_id"`
	Url                   string    `json:"url"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// type Certificate struct {
// 	ID             int64  `json:"id"`
// 	Name           string `json:"name"`
// 	Organization   string `json:"origanization"`
// 	IssueDate      string `json:"issue_date"`
// 	ExpirationDate string `json:"expiration_date"`
// 	CredentialID   string `json:"credential_id"`
// 	Url            string `json:"url"`
// }

type UserDetailResponse struct {
	ID              *int64     `json:"id"`
	UserID          *int64     `json:"user_id"`
	PhoneNumber     *string    `json:"phone_number"`
	Gender          *string    `json:"gender"`
	Location        *string    `json:"location"`
	PortfolioUrl    *string    `json:"portfolio_url"`
	About           *string    `json:"about"`
	HidePhoneNumber *bool      `json:"hide_phone_number"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

type SkillsResponse struct {
	ID   *int64  `json:"id"`
	Name *string `json:"name"`
}

type UserProfile struct {
	User            User          `json:"user"`
	FollowingCount  int64         `json:"following_count"`
	SocialLinks     []SocialLinks `json:"social_links"`
	Skills          UserSkills    `json:"skills"`
	Location        string        `json:"location"`
	WebPortfolioUrl string        `json:"web_portfolio_url"`
	About           string        `json:"about"`
}

type UserSkills struct {
	MainSkills  []string `json:"main_skills"`
	OtherSkills []string `json:"other_skills"`
}

// type WorkExperience struct {
// 	ID             int64    `json:"id"`
// 	JobTitle       string   `json:"job_title"`
// 	Company        Company  `json:"company"`
// 	EmploymentType string   `json:"employment_type"`
// 	Location       string   `json:"location"`
// 	LocationType   string   `json:"location_type"`
// 	StartDate      string   `json:"start_date"`
// 	FinishDate     string   `json:"finish_date"`
// 	Description    string   `json:"description"`
// 	FileURLs       []string `json:"file_urls"`
// 	Skills         []string `json:"skills"`
// }

// type Education struct {
// 	ID           int64    `json:"id"`
// 	School       School   `json:"school"`
// 	Degree       string   `json:"degree"`
// 	FieldOfStudy string   `json:"field_of_study"`
// 	StartDate    string   `json:"start_date"`
// 	FinishDate   string   `json:"finish_date"`
// 	GPA          string   `json:"gpa"`
// 	Description  string   `json:"description" `
// 	FileURLs     []string `json:"file_urls"`
// 	Skills       []string `json:"skills"`
// }
