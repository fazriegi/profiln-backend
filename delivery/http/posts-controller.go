package http

import (
	"net/http"
	"profiln-be/delivery/ws"
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
	ListNewestPostsByTargetUser(ctx *gin.Context)
	ListLikedPostsByTargetUser(ctx *gin.Context)
	ListRepostedPostsByTargetUser(ctx *gin.Context)
	InsertPost(ctx *gin.Context)
	UpdatePost(ctx *gin.Context)
	DeletePost(ctx *gin.Context)
	RepostPost(ctx *gin.Context)
	UnrepostPost(ctx *gin.Context)
	UploadFileForInsertPost(ctx *gin.Context)
	UploadFileForUpdatePost(ctx *gin.Context)
	InsertPostComment(ctx *gin.Context)
	LikePostComment(ctx *gin.Context)
	UnlikePostComment(ctx *gin.Context)
}

type PostsController struct {
	usecase posts.IPostsUsecase
	hub     *ws.Hub
}

func NewPostsController(usecase posts.IPostsUsecase, hub *ws.Hub) IPostsController {
	return &PostsController{
		usecase,
		hub,
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

func (c *PostsController) ListNewestPostsByTargetUser(ctx *gin.Context) {
	var response model.Response

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	targetUserId, err := strconv.ParseInt(ctx.Param("userId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

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

	response = c.usecase.ListNewestPostsByTargetUser(userId, targetUserId, pagination)
	ctx.JSON(response.Status.Code, response)
}

func (c *PostsController) ListLikedPostsByTargetUser(ctx *gin.Context) {
	var response model.Response

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	targetUserId, err := strconv.ParseInt(ctx.Param("userId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

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

	response = c.usecase.ListLikedPostsByTargetUser(userId, targetUserId, pagination)
	ctx.JSON(response.Status.Code, response)
}

func (c *PostsController) ListRepostedPostsByTargetUser(ctx *gin.Context) {
	var response model.Response

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	targetUserId, err := strconv.ParseInt(ctx.Param("userId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request query")

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

	response = c.usecase.ListRepostedPostsByTargetUser(userId, targetUserId, pagination)
	ctx.JSON(response.Status.Code, response)
}

func (c *PostsController) InsertPost(ctx *gin.Context) {
	var (
		reqBody  model.CreatePostRequest
		response model.Response
	)

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

	response = c.usecase.InsertPost(&reqBody)
	ctx.JSON(response.Status.Code, response)
}

func (c *PostsController) UpdatePost(ctx *gin.Context) {
	var (
		response model.Response
		reqBody  model.UpdatePostRequest
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

	reqBody.UserId = userId
	reqBody.ID = postId

	response = c.usecase.UpdatePost(&reqBody)
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

func (c *PostsController) UploadFileForInsertPost(ctx *gin.Context) {
	var response model.Response

	fileNames := ctx.MustGet("fileNames").([]string)
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	postId, err := strconv.ParseInt(ctx.Param("postId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

	response = c.usecase.UploadFileForInsertPost(userId, postId, fileNames)

	ctx.JSON(response.Status.Code, response)
}

func (c *PostsController) UploadFileForUpdatePost(ctx *gin.Context) {
	var response model.Response

	fileNames := ctx.MustGet("fileNames").([]string)
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	postId, err := strconv.ParseInt(ctx.Param("postId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

	response = c.usecase.UploadFileForUpdatePost(userId, postId, fileNames)

	ctx.JSON(response.Status.Code, response)
}

func (c *PostsController) InsertPostComment(ctx *gin.Context) {
	var (
		response model.Response
		reqBody  model.AddPostCommentReq
	)
	fileNames := ctx.MustGet("fileNames").([]string)
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
	reqBody.PostId = postId

	response = c.usecase.InsertPostComment(fileNames, &reqBody)

	if response.Status.IsSuccess {
		comment := ws.Message{
			PostId: postId,
			Data:   response.Data,
		}
		c.hub.Broadcast(comment)
	}

	ctx.JSON(response.Status.Code, response)
}

func (c *PostsController) LikePostComment(ctx *gin.Context) {
	var response model.Response

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	postCommentId, err := strconv.ParseInt(ctx.Param("postCommentId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

	response = c.usecase.LikePostComment(userId, postCommentId)
	ctx.JSON(response.Status.Code, response)
}

func (c *PostsController) UnlikePostComment(ctx *gin.Context) {
	var response model.Response

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := int64(userData["id"].(float64))

	postCommentId, err := strconv.ParseInt(ctx.Param("postCommentId"), 10, 64)
	if err != nil {
		response.Status =
			libs.CustomResponse(http.StatusBadRequest, "Invalid request param")

		ctx.JSON(response.Status.Code, response)
		return
	}

	response = c.usecase.UnlikePostComment(userId, postCommentId)
	ctx.JSON(response.Status.Code, response)
}
