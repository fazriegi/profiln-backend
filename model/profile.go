package model

type CompanyRequest struct {
	Name string `validate:"required"`
}
