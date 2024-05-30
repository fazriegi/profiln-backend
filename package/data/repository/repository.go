package data

import (
	"context"
	"database/sql"
	db "profiln-be/db/sqlc"
	"profiln-be/model"
)

type IDataRepository interface {
	GetSchools(offset, limit int32) ([]model.School, int64, error)
	GetCompanies(offset, limit int32) ([]model.Company, int64, error)
	GetIssuingOrganizations(offset, limit int32) ([]model.IssuingOrganization, int64, error)
	GetSkills(offset, limit int32) ([]model.Skill, int64, error)
}

type DataRepository struct {
	dbConn *sql.DB
	query  *db.Queries
}

func NewDataRepository(dbConn *sql.DB) IDataRepository {
	return &DataRepository{
		dbConn: dbConn,
		query:  db.New(dbConn),
	}
}

func (r *DataRepository) GetSchools(offset, limit int32) ([]model.School, int64, error) {
	schools, err := r.query.GetSchools(context.Background(), db.GetSchoolsParams{
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		return nil, 0, err
	}

	var count int64
	if len(schools) > 0 {
		count = schools[0].TotalRows
	}

	data := make([]model.School, len(schools))
	for i, v := range schools {
		data[i] = model.School{
			ID:   v.ID,
			Name: v.Name,
		}
	}

	return data, count, nil

}

func (r *DataRepository) GetCompanies(offset, limit int32) ([]model.Company, int64, error) {
	companies, err := r.query.GetCompanies(context.Background(), db.GetCompaniesParams{
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		return nil, 0, err
	}

	var count int64
	if len(companies) > 0 {
		count = companies[0].TotalRows
	}

	data := make([]model.Company, len(companies))
	for i, v := range companies {
		data[i] = model.Company{
			ID:   v.ID,
			Name: v.Name,
		}
	}

	return data, count, nil
}

func (r *DataRepository) GetIssuingOrganizations(offset, limit int32) ([]model.IssuingOrganization, int64, error) {
	issuingOrganizations, err := r.query.GetIssuingOrganizations(context.Background(), db.GetIssuingOrganizationsParams{
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		return nil, 0, err
	}

	var count int64
	if len(issuingOrganizations) > 0 {
		count = issuingOrganizations[0].TotalRows
	}

	data := make([]model.IssuingOrganization, len(issuingOrganizations))
	for i, v := range issuingOrganizations {
		data[i] = model.IssuingOrganization{
			ID:   v.ID,
			Name: v.Name,
		}
	}

	return data, count, nil
}

func (r *DataRepository) GetSkills(offset, limit int32) ([]model.Skill, int64, error) {
	skills, err := r.query.GetSkills(context.Background(), db.GetSkillsParams{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		return []model.Skill{}, 0, err
	}

	var count int64
	if len(skills) > 0 {
		count = skills[0].TotalRows
	}

	data := make([]model.Skill, len(skills))
	for i, v := range skills {
		data[i] = model.Skill{
			ID:   v.ID,
			Name: v.Name,
		}
	}

	return data, count, nil
}
