package model

import (
	"database/sql"
	"time"
)

type CompanyRequest struct {
	Name string `validate:"required"`
}

type IssuingOrganizationRequest struct {
	Name string `validate:"required"`
}

type SchoolRequest struct {
	Name string `validate:"required"`
}

type SkillRequest struct {
	Name string `validate:"required"`
}

type UserSkillRequest struct {
	UserID    int64  `validate:"required"`
	SkillID   int64  `validate:"required"`
	MainSkill bool   `validate:"required"`
	Skills    string `validate:"required"`
}

type CertificateRequest struct {
	Name                  string  `validate:"required"`
	IssuingOrganizationID int64   `validate:"required"`
	IssueDate             string  `validate:"required"`
	ExpirationDate        *string `validate:"omitempty"`
	CredentialID          string  `validate:"required"`
	Url                   string  `validate:"required"`
}

type WorkExperienceRequest struct {
	// UserID           int64        `validate:"required"`
	JobTitle       string       `validate:"required"`
	CompanyID      int64        `validate:"required"`
	EmploymentType string       `validate:"required"`
	Location       string       `validate:"required"`
	LocationType   string       `validate:"required"`
	StartDate      sql.NullTime `validate:"required"`
	FinishDate     sql.NullTime `validate:"required"`
	Description    string       `validate:"required"`
}

type EducationRequest struct {
	// UserID           int64        `validate:"required"`
	SchoolID     int64        `validate:"required"`
	Degree       string       `validate:"required"`
	FieldOfStudy string       `validate:"required"`
	Gpa          string       `validate:"required"`
	StartDate    sql.NullTime `validate:"required"`
	FinishDate   sql.NullTime `validate:"required"`
}

type UserDetailAboutRequest struct {
	About string `json:"about" validate:"required"`
	// UserID sql.NullInt64
}

type UserDetailRequest struct {
	PhoneNumber string `validate:"required"`
	Gender      string `validate:"required"`
}

type GetSkillsResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UpdateProfileRequest struct {
	UserId          int64         `json:"user_id" form:"user_id"`
	Fullname        string        `json:"fullname" form:"fullname" validate:"required"`
	HidePhoneNumber bool          `json:"hide_phone_number" form:"hide_phone_number" validate:"required"`
	MainSkills      []string      `json:"main_skills" form:"main_skills" validate:"required"`
	PhoneNumber     string        `json:"phone_number" form:"phone_number" validate:"required"`
	Gender          string        `json:"gender" form:"gender" validate:"required"`
	SocialLinks     []SocialLinks `json:"social_links" form:"social_links" validate:"required"`
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
	Name string `json:"name" form:"name" validate:"required"`
	URL  string `json:"url" form:"url" validate:"required"`
}

type UpdateCertificate struct {
	ID                  int64               `json:"id"`
	Name                string              `json:"name" validate:"required"`
	IssuingOrganization IssuingOrganization `json:"issuing_organization" validate:"required"`
	IssueDate           string              `json:"issue_date" validate:"required"`
	ExpirationDate      string              `json:"expiration_date"`
	CredentialID        string              `json:"credential_id"`
	Url                 string              `json:"url"`
}

type IssuingOrganization struct {
	ID   int64  `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type UpdateUserInformation struct {
	UserId       int64    `json:"user_id"`
	Skills       []string `json:"skills"`
	Location     string   `json:"location" validate:"required"`
	PortfolioUrl string   `json:"portfolio_url"`
}

type UserDetail struct {
	ID              int64
	UserId          int64
	PhoneNumber     string
	Gender          string
	Location        string
	PortfolioUrl    string
	About           string
	HidePhoneNumber bool
}

type UpdateEducationRequest struct {
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

type School struct {
	ID   int64  `json:"id"`
	Name string `json:"name" validate:"required"`
}

type Skill struct {
	ID   int64  `json:"id"`
	Name string `json:"name" validate:"required"`
}

type UpdateWorkExperience struct {
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

type Company struct {
	ID   int64  `json:"id"`
	Name string `json:"name" validate:"required"`
}

type AboutProfileResponse struct {
	ID        *int64     `json:"id"`
	UserID    *int64     `json:"user_id"`
	About     *string    `json:"about"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type UserProfileResponse struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"fullname"`
}

type UserAboutResponse struct {
	About AboutProfileResponse
	User  UserProfileResponse
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

type UserCertificatesResponse struct {
	User         UserProfileResponse
	Certificates any
}

type CertificatesResponse struct {
	ID             int64      `json:"id"`
	Name           string     `json:"name"`
	Organization   string     `json:"origanization"`
	IssueDate      time.Time  `json:"issue_date"`
	ExpirationDate *time.Time `json:"expiration_date"`
	CredentialID   string     `json:"credential_id"`
	Url            string     `json:"url"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

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

type UserSkillsLocationResponse struct {
	User       UserProfileResponse
	Skills     []SkillsResponse
	UserDetail UserDetailResponse
}
