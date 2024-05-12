package model

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
