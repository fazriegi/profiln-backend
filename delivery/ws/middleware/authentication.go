package ws_middleware

import (
	"net/http"
	"profiln-be/libs"
	"profiln-be/model"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var response = model.Response{}

		token := ctx.Request.URL.Query().Get("token")

		verifiedToken, err := libs.VerifyJWTTOken(token)
		if err != nil {
			status := libs.CustomResponse(http.StatusUnauthorized, err.Error())
			response.Status = status

			ctx.AbortWithStatusJSON(status.Code, response)
			return
		}

		ctx.Set("userData", verifiedToken)
		ctx.Next()
	}
}
