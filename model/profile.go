package model

import "database/sql"

type CompanyRequest struct {
	Name string `validate:"required"`
}

type IssuingOrganizationRequest struct {
	Name string `validate:"required"`
}

type EmploymentTypeRequest struct {
	Name string `validate:"required"`
}

type LocationTypeRequest struct {
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
	// UserID                int64        `validate:"required"`
	Name                  string       `validate:"required"`
	IssuingOrganizationID int64        `validate:"required"`
	IssueDate             sql.NullTime `validate:"required"`
	ExpirationDate        sql.NullTime `validate:"required"`
	CredentialID          string       `validate:"required"`
	Url                   string       `validate:"required"`
}

type WorkExperienceRequest struct {
	// UserID           int64        `validate:"required"`
	JobTitle         string       `validate:"required"`
	CompanyID        int64        `validate:"required"`
	EmploymentTypeID int16        `validate:"required"`
	Location         string       `validate:"required"`
	LocationTypeID   int16        `validate:"required"`
	StartDate        sql.NullTime `validate:"required"`
	FinishDate       sql.NullTime `validate:"required"`
	Description      string       `validate:"required"`
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
	About string `validate:"required"`
	// UserID sql.NullInt64
}

type UserDetailRequest struct {
	PhoneNumber string `validate:"required"`
	Gender      string `validate:"required"`
}
