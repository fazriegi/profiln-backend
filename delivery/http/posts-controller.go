package http

import (
	"mime/multipart"
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
	GetDetailPost(ctx *gin.Context)
	GetPostComments(ctx *gin.Context)
	GetPostCommentReplies(ctx *gin.Context)
	LikePost(ctx *gin.Context)
	UnlikePost(ctx *gin.Context)
	ListNewestPostsByUserId(ctx *gin.Context)
	ListLikedPostsByUserId(ctx *gin.Context)
	ListRepostedPostsByUserId(ctx *gin.Context)
	InsertPost(ctx *gin.Context)
	UpdatePost(ctx *gin.Context)
	DeletePost(ctx *gin.Context)
	RepostPost(ctx *gin.Context)
	UnrepostPost(ctx *gin.Context)
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
		reqBody  model.ReportPost
		response model.Response
	)
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	postId, err := strconv.ParseInt(ctx.Param("postId"), 10, 64)
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

	reqBody.PostId = postId

	response = c.usecase.InsertReportedPost(userId, &reqBody)
	ctx.JSON(response.Status.Code, response)
}

func (c *PostsController) GetDetailPost(ctx *gin.Context) {
	var response model.Response

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	postId, err := strconv.ParseInt(ctx.Param("postId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

	response = c.usecase.GetDetailPost(postId, userId)
	ctx.JSON(response.Status.Code, response)
}

func (c *PostsController) GetPostComments(ctx *gin.Context) {
	var response model.Response

	postId, err := strconv.ParseInt(ctx.Param("postId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

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

	response = c.usecase.GetPostComments(postId, pagination)
	ctx.JSON(response.Status.Code, response)
}

func (c *PostsController) GetPostCommentReplies(ctx *gin.Context) {
	var response model.Response

	postId, err := strconv.ParseInt(ctx.Param("postId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

	postCommentId, err := strconv.ParseInt(ctx.Param("postCommentId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

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

	response = c.usecase.GetPostCommentReplies(postId, postCommentId, pagination)
	ctx.JSON(response.Status.Code, response)
}

func (c *PostsController) LikePost(ctx *gin.Context) {
	var response model.Response

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	postId, err := strconv.ParseInt(ctx.Param("postId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

	response = c.usecase.LikePost(userId, postId)
	ctx.JSON(response.Status.Code, response)
}

func (c *PostsController) UnlikePost(ctx *gin.Context) {
	var response model.Response

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	postId, err := strconv.ParseInt(ctx.Param("postId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

	response = c.usecase.UnlikePost(userId, postId)
	ctx.JSON(response.Status.Code, response)
}

func (c *PostsController) ListNewestPostsByUserId(ctx *gin.Context) {
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

	response = c.usecase.ListNewestPostsByUserId(userId, pagination)
	ctx.JSON(response.Status.Code, response)
}

func (c *PostsController) ListLikedPostsByUserId(ctx *gin.Context) {
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

	response = c.usecase.ListLikedPostsByUserId(userId, pagination)
	ctx.JSON(response.Status.Code, response)
}

func (c *PostsController) ListRepostedPostsByUserId(ctx *gin.Context) {
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

	response = c.usecase.ListRepostedPostsByUserId(userId, pagination)
	ctx.JSON(response.Status.Code, response)
}

func (c *PostsController) InsertPost(ctx *gin.Context) {
	var (
		reqBody  model.CreatePostRequest
		response model.Response
	)
	imageFiles := ctx.MustGet("files").([]*multipart.FileHeader)
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

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

	reqBody.UserId = userId

	response = c.usecase.InsertPost(imageFiles, &reqBody)
	ctx.JSON(response.Status.Code, response)
}

func (c *PostsController) UpdatePost(ctx *gin.Context) {
	var (
		response model.Response
		reqBody  model.UpdatePostRequest
	)
	files := ctx.MustGet("files").([]*multipart.FileHeader)
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	postId, err := strconv.ParseInt(ctx.Param("postId"), 10, 64)
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

	reqBody.UserId = userId
	reqBody.ID = postId

	response = c.usecase.UpdatePost(files, &reqBody)
	ctx.JSON(response.Status.Code, response)
}

func (c *PostsController) DeletePost(ctx *gin.Context) {
	var (
		response model.Response
	)

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	postId, err := strconv.ParseInt(ctx.Param("postId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

	response = c.usecase.DeletePost(userId, postId)
	ctx.JSON(response.Status.Code, response)
}

func (c *PostsController) RepostPost(ctx *gin.Context) {
	var response model.Response

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	postId, err := strconv.ParseInt(ctx.Param("postId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

	response = c.usecase.RepostPost(userId, postId)
	ctx.JSON(response.Status.Code, response)
}

func (c *PostsController) UnrepostPost(ctx *gin.Context) {
	var response model.Response

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	postId, err := strconv.ParseInt(ctx.Param("postId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

	response = c.usecase.UnrepostPost(userId, postId)
	ctx.JSON(response.Status.Code, response)
}
