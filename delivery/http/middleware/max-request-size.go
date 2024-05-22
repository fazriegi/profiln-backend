package middleware

import (
	"fmt"
	"net/http"
	"profiln-be/libs"
	"profiln-be/model"

	"github.com/gin-gonic/gin"
)

func MaxReqSizeAllowed(maxBytes int64) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		respMessage := fmt.Sprintf("File is too large. Maximum allowed size is %d bytes", maxBytes)
		if ctx.Request.ContentLength > maxBytes {
			response := model.Response{
				Status: libs.CustomResponse(http.StatusRequestEntityTooLarge, respMessage),
			}

			ctx.AbortWithStatusJSON(response.Status.Code, response)
			return
		}
		ctx.Next()
	}
}
