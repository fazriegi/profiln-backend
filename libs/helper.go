package libs

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"html/template"
	"io"
	"profiln-be/model"
	"time"
)

func CustomResponse(code int, message string) model.Status {
	statuses := map[int]string{
		500: "internal server error",
		422: "unprocessable content",
		413: "request entity too large",
		404: "not found",
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

func GenerateOTP(max int) (string, error) {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, max)

	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		return "", err
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b), nil
}

func HTMLToString(filepath string, data any) (string, error) {
	// get email template
	t, err := template.ParseFiles(filepath)
	if err != nil {
		return "", err
	}

	buff := new(bytes.Buffer)

	err = t.Execute(buff, data)
	if err != nil {
		return "", err
	}

	return buff.String(), nil
}

func GenerateNewFilename(filename string) string {
	currentTime := time.Now().UnixNano()
	newFilename := fmt.Sprintf("%d_%s", currentTime, filename)
	return newFilename
}
