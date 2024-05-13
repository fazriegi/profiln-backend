package http

import (
	"net/http"
	"profiln-be/libs"
	"profiln-be/model"
	"profiln-be/package/posts"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type IPostsController interface {
	ReportPost(ctx *gin.Context)
}

type PostsController struct {
	usecase posts.IPostsUsecase
}

func NewPostsController(usecase posts.IPostsUsecase) IPostsController {
	return &PostsController{
		usecase,
	}
}

func (c *PostsController) ReportPost(ctx *gin.Context) {
	var (
		reqBody  model.ReportPostRequest
		response model.Response
	)
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	postId, err := strconv.Atoi(ctx.Param("postId"))
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

	if err := ctx.ShouldBind(&reqBody); err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Error parsing request body")

		ctx.JSON(response.Status.Code, response)
		return
	}

	validationErr := libs.ValidateRequest(reqBody) // validate reqBody struct
	// if there is an error
	if len(validationErr) > 0 {
		errResponse := map[string]any{
			"errors": validationErr,
		}

		response.Status =
			libs.CustomResponse(http.StatusUnprocessableEntity, "Validation error")
		response.Data = errResponse

		ctx.JSON(response.Status.Code, response)
		return
	}

	reqBody.PostId = int64(postId)

	response = c.usecase.InsertReportedPost(userId, &reqBody)
	ctx.JSON(response.Status.Code, response)
}
