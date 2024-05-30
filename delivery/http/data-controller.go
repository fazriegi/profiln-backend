package http

import (
	"net/http"
	"profiln-be/libs"
	"profiln-be/model"
	"profiln-be/package/data"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IDataController interface {
	GetSchools(ctx *gin.Context)
	GetCompanies(ctx *gin.Context)
	GetIssuingOrganizations(ctx *gin.Context)
	GetSkills(ctx *gin.Context)
}

type DataController struct {
	usecase data.IDataUsecase
}

func NewDataController(usecase data.IDataUsecase) IDataController {
	return &DataController{
		usecase,
	}
}

func (c *DataController) GetSchools(ctx *gin.Context) {
	var response model.Response

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

		ctx.JSON(response.Status.Code, response)
		return
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

		ctx.JSON(response.Status.Code, response)
		return
	}

	if page <= 0 || limit <= 0 {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

		ctx.JSON(response.Status.Code, response)
		return
	}

	pagination := model.PaginationRequest{
		Page:  page,
		Limit: limit,
	}
	response = c.usecase.GetSchools(pagination)
	ctx.JSON(response.Status.Code, response)
}

func (c *DataController) GetCompanies(ctx *gin.Context) {
	var response model.Response

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

		ctx.JSON(response.Status.Code, response)
		return
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

		ctx.JSON(response.Status.Code, response)
		return
	}

	if page <= 0 || limit <= 0 {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

		ctx.JSON(response.Status.Code, response)
		return
	}

	pagination := model.PaginationRequest{
		Page:  page,
		Limit: limit,
	}
	response = c.usecase.GetCompanies(pagination)
	ctx.JSON(response.Status.Code, response)
}

func (c *DataController) GetIssuingOrganizations(ctx *gin.Context) {
	var response model.Response

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

		ctx.JSON(response.Status.Code, response)
		return
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

		ctx.JSON(response.Status.Code, response)
		return
	}

	if page <= 0 || limit <= 0 {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

		ctx.JSON(response.Status.Code, response)
		return
	}

	pagination := model.PaginationRequest{
		Page:  page,
		Limit: limit,
	}
	response = c.usecase.GetIssuingOrganizations(pagination)
	ctx.JSON(response.Status.Code, response)
}

func (c *DataController) GetSkills(ctx *gin.Context) {
	var response model.Response

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

		ctx.JSON(response.Status.Code, response)
		return
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

		ctx.JSON(response.Status.Code, response)
		return
	}

	if page <= 0 || limit <= 0 {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

		ctx.JSON(response.Status.Code, response)
		return
	}

	pagination := model.PaginationRequest{
		Page:  page,
		Limit: limit,
	}
	response = c.usecase.GetSkills(pagination)
	ctx.JSON(response.Status.Code, response)
}
