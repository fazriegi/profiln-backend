package http

import (
	"net/http"
	"profiln-be/libs"
	"profiln-be/model"
	"profiln-be/package/homepage"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type IHomepageController interface {
	ListPosts(ctx *gin.Context)
	ListFollowsRecommendation(ctx *gin.Context)
}

type HomepageController struct {
	usecase homepage.IHomepageUsecase
}

func NewHomepageController(usecase homepage.IHomepageUsecase) IHomepageController {
	return &HomepageController{
		usecase,
	}
}

func (c *HomepageController) ListPosts(ctx *gin.Context) {
	var response model.Response

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))
	orderBy := ctx.Query("orderBy")

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

	if orderBy != "newest" && orderBy != "following" && orderBy != "popular" {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

		ctx.JSON(response.Status.Code, response)
		return
	}

	pagination := model.PaginationRequest{
		Page:    page,
		Limit:   limit,
		OrderBy: orderBy,
	}

	response = c.usecase.ListPosts(userId, pagination)
	ctx.JSON(response.Status.Code, response)
}

func (c *HomepageController) ListFollowsRecommendation(ctx *gin.Context) {
	var response model.Response
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

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

	response = c.usecase.ListFollowsRecommendation(userId, pagination)
	ctx.JSON(response.Status.Code, response)
}
