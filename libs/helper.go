package libs

import (
	"profiln-be/model"
)

func CustomResponse(code int, message string) model.Status {
	statuses := map[int]string{
		500: "internal server error",
		422: "unprocessable content",
		401: "unauthorized",
		400: "bad request",
		303: "redirect",
		201: "created",
		200: "success",
	}

	var status model.Status
	isSuccess := code >= 200 && code <= 299

	status.Code = code
	status.Message = message
	status.Status = statuses[code]
	status.IsSuccess = isSuccess

	return status
}
