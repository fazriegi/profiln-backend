package data

import (
	"net/http"
	"profiln-be/libs"
	"profiln-be/model"
	repository "profiln-be/package/data/repository"

	"github.com/sirupsen/logrus"
)

type IDataUsecase interface {
	GetSchools(pagination model.PaginationRequest) model.Response
	GetCompanies(pagination model.PaginationRequest) model.Response
	GetIssuingOrganizations(pagination model.PaginationRequest) model.Response
	GetSkills(pagination model.PaginationRequest) model.Response
	GetJobPositions(pagination model.PaginationRequest) model.Response
}

type DataUsecase struct {
	repository repository.IDataRepository
	log        *logrus.Logger
}

func NewDataUsecase(repository repository.IDataRepository, log *logrus.Logger) IDataUsecase {
	return &DataUsecase{
		repository,
		log,
	}
}

func (u *DataUsecase) GetSchools(pagination model.PaginationRequest) model.Response {
	offset := (pagination.Page - 1) * pagination.Limit

	data, totalRows, err := u.repository.GetSchools(int32(offset), int32(pagination.Limit))
	if err != nil {
		u.log.Errorf("repository.GetSchools: %v", err)
		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
		}
	}

	totalPages := int((totalRows + int64(pagination.Limit) - 1) / int64(pagination.Limit))

	paginate := model.PaginationResponse{
		Page:             pagination.Page,
		TotalRows:        totalRows,
		TotalPages:       totalPages,
		CurrentRowsCount: len(data),
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusOK, "Success fetch schools"),
		Data: map[string]any{
			"pagination": paginate,
			"data":       data,
		},
	}
}

func (u *DataUsecase) GetCompanies(pagination model.PaginationRequest) model.Response {
	offset := (pagination.Page - 1) * pagination.Limit

	data, totalRows, err := u.repository.GetCompanies(int32(offset), int32(pagination.Limit))
	if err != nil {
		u.log.Errorf("repository.GetCompanies: %v", err)
		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
		}
	}

	totalPages := int((totalRows + int64(pagination.Limit) - 1) / int64(pagination.Limit))

	paginate := model.PaginationResponse{
		Page:             pagination.Page,
		TotalRows:        totalRows,
		TotalPages:       totalPages,
		CurrentRowsCount: len(data),
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusOK, "Success fetch companies"),
		Data: map[string]any{
			"pagination": paginate,
			"data":       data,
		},
	}
}

func (u *DataUsecase) GetIssuingOrganizations(pagination model.PaginationRequest) model.Response {
	offset := (pagination.Page - 1) * pagination.Limit

	data, totalRows, err := u.repository.GetIssuingOrganizations(int32(offset), int32(pagination.Limit))
	if err != nil {
		u.log.Errorf("repository.GetIssuingOrganizations: %v", err)
		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
		}
	}

	totalPages := int((totalRows + int64(pagination.Limit) - 1) / int64(pagination.Limit))

	paginate := model.PaginationResponse{
		Page:             pagination.Page,
		TotalRows:        totalRows,
		TotalPages:       totalPages,
		CurrentRowsCount: len(data),
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusOK, "Success fetch issuing organizations"),
		Data: map[string]any{
			"pagination": paginate,
			"data":       data,
		},
	}
}

func (u *DataUsecase) GetSkills(pagination model.PaginationRequest) model.Response {
	offset := (pagination.Page - 1) * pagination.Limit

	data, totalRows, err := u.repository.GetSkills(int32(offset), int32(pagination.Limit))
	if err != nil {
		u.log.Errorf("repository.GetSkills: %v", err)
		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
		}
	}

	totalPages := int((totalRows + int64(pagination.Limit) - 1) / int64(pagination.Limit))

	paginate := model.PaginationResponse{
		Page:             pagination.Page,
		TotalRows:        totalRows,
		TotalPages:       totalPages,
		CurrentRowsCount: len(data),
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusOK, "Success fetch skills"),
		Data: map[string]any{
			"pagination": paginate,
			"data":       data,
		},
	}
}

func (u *DataUsecase) GetJobPositions(pagination model.PaginationRequest) model.Response {
	offset := (pagination.Page - 1) * pagination.Limit

	data, totalRows, err := u.repository.GetJobPositions(int32(offset), int32(pagination.Limit))
	if err != nil {
		u.log.Errorf("repository.GetJobPositions: %v", err)
		return model.Response{
			Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
		}
	}

	totalPages := int((totalRows + int64(pagination.Limit) - 1) / int64(pagination.Limit))

	paginate := model.PaginationResponse{
		Page:             pagination.Page,
		TotalRows:        totalRows,
		TotalPages:       totalPages,
		CurrentRowsCount: len(data),
	}

	return model.Response{
		Status: libs.CustomResponse(http.StatusOK, "Success fetch job positions"),
		Data: map[string]any{
			"pagination": paginate,
			"data":       data,
		},
	}
}
