package middleware

import (
	"net/http"
	"strings"

	"profiln-be/libs"
	"profiln-be/model"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var response = model.Response{}
		header := ctx.Request.Header.Get("Authorization")
		isHasBearer := strings.HasPrefix(header, "Bearer")

		if !isHasBearer {
			status := libs.CustomResponse(http.StatusUnauthorized, "Sign in to proceed")
			response.Status = status

			ctx.AbortWithStatusJSON(status.Code, response)
			return
		}

		tokenString := strings.Split(header, " ")[1]

		verifiedToken, err := libs.VerifyJWTTOken(tokenString)
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
