package libs

import (
	"github.com/go-playground/validator"
)

type ValidationErrResponse struct {
	FailedField string `json:"failed_field"`
	Tag         string `json:"tag"`
	TagValue    string `json:"tag_value"`
}

func ValidateRequest(data any) []ValidationErrResponse {
	var validationErrors []ValidationErrResponse

	validate := validator.New()

	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var validationErr ValidationErrResponse
			validationErr.FailedField = err.Field()
			validationErr.Tag = err.Tag()
			validationErr.TagValue = err.Param()

			validationErrors = append(validationErrors, validationErr)
		}
	}

	return validationErrors
}
