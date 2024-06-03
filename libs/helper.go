package libs

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"html/template"
	"io"
	"mime/multipart"
	"path/filepath"
	"profiln-be/model"
	"strings"
	"time"
)

func CustomResponse(code int, message string) model.Status {
	statuses := map[int]string{
		500: "internal server error",
		422: "unprocessable content",
		415: "unsupported media type",
		413: "request entity too large",
		404: "not found",
		401: "unauthorized",
		400: "bad request",
		303: "redirect",
		204: "no content",
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

func IsFileSizeAllowed(allowedSize int64, files ...*multipart.FileHeader) bool {
	for _, file := range files {
		if file.Size > allowedSize {
			return false
		}
	}

	return true
}

func IsFileExtensionAllowed(allowedExtensions []string, file *multipart.FileHeader) bool {
	fileTypes := map[string]string{
		".pdf":  "application/pdf",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".png":  "image/png",
		".jpg":  "image/jpeg",
	}

	fileExt := strings.ToLower(filepath.Ext(file.Filename))
	fileType := file.Header.Get("Content-Type")

	// Check if the file extension and type are allowed
	for _, ext := range allowedExtensions {
		if ext == fileExt && fileTypes[ext] == fileType {
			return true
		}
	}

	return false
}

func ParseTimeWithNill(dateString *string) (sql.NullTime, error) {
	if dateString == nil || *dateString == "" {
		return sql.NullTime{Valid: false}, nil
	}

	date, err := time.Parse("2006-01-02", *dateString)
	if err != nil {
		return sql.NullTime{}, err
	}

	return sql.NullTime{Time: date, Valid: true}, nil
}

func ParseTime(dateString string) (time.Time, error) {
	date, err := time.Parse("2006-01-02", dateString)

	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}
