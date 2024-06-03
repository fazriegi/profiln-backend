package middleware

import (
	"fmt"
	"net/http"
	"profiln-be/libs"
	"profiln-be/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ValidateFileUpload(maxBytes int64, maxTotalFile uint8, allowedExtensions []string, fs libs.IFileSystem, log *logrus.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		respMessageSize := fmt.Sprintf("File is too large. Maximum allowed size is %d bytes", maxBytes)
		respMessageCount := fmt.Sprintf("Too many files. Maximum allowed is %d files", maxTotalFile)
		respMessageFormat := fmt.Sprintf("File format not allowed. Allowed formats are: %v", allowedExtensions)

		form, err := ctx.MultipartForm()
		if err != nil {
			response := model.Response{
				Status: libs.CustomResponse(http.StatusBadRequest, "Error parsing form data"),
			}
			ctx.AbortWithStatusJSON(response.Status.Code, response)
			return
		}

		files := form.File["files"]
		if len(files) > int(maxTotalFile) {
			response := model.Response{
				Status: libs.CustomResponse(http.StatusBadRequest, respMessageCount),
			}
			ctx.AbortWithStatusJSON(response.Status.Code, response)
			return
		}

		fileNames := make([]string, len(files))

		for i, file := range files {
			if !libs.IsFileExtensionAllowed(allowedExtensions, file) {
				response := model.Response{
					Status: libs.CustomResponse(http.StatusUnsupportedMediaType, respMessageFormat),
				}
				ctx.AbortWithStatusJSON(response.Status.Code, response)
				return
			}

			if file.Size > maxBytes {
				response := model.Response{
					Status: libs.CustomResponse(http.StatusRequestEntityTooLarge, respMessageSize),
				}
				ctx.AbortWithStatusJSON(response.Status.Code, response)
				return
			}

			newFilename := fs.GenerateNewFilename(file.Filename)
			filePath := fmt.Sprintf("./storage/temp/file/%s", newFilename)

			if err := fs.SaveFile(file, filePath); err != nil {
				log.Errorf("fileSystem.SaveFile: %v", err)
				response := model.Response{
					Status: libs.CustomResponse(http.StatusInternalServerError, "Unexpected error occured"),
				}

				ctx.AbortWithStatusJSON(response.Status.Code, response)
				return
			}
			fileNames[i] = newFilename
		}

		ctx.Set("fileNames", fileNames)
		ctx.Next()
	}
}
