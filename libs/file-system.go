package libs

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

type IFileSystem interface {
	SaveFile(file *multipart.FileHeader, dst string) error
	RemoveFile(filepath string) error
	GenerateNewFilename(filename string) string
}

type FileSystem struct{}

func NewFileSystem() IFileSystem {
	return &FileSystem{}
}

func (f *FileSystem) SaveFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func (f *FileSystem) RemoveFile(filepath string) error {
	if err := os.Remove(filepath); err != nil {
		return err
	}

	return nil
}

func (f *FileSystem) GenerateNewFilename(filename string) string {
	currentTime := time.Now().Format("20060102_150405")

	// Clean the filename by removing or replacing invalid characters
	cleanFilename := f.cleanString(filename)

	newFilename := fmt.Sprintf("%s_%s", currentTime, cleanFilename)
	return newFilename
}

// cleanString removes invalid characters from a filename
func (f *FileSystem) cleanString(s string) string {
	// Sanitize the path
	cleaned := filepath.Clean(s)

	// Remove or replace any remaining invalid characters
	// and replace it with underscore
	var sb strings.Builder
	for _, r := range cleaned {
		if f.isAllowedRune(r) {
			sb.WriteRune(r)
		} else {
			sb.WriteRune('_')
		}
	}
	return sb.String()
}

// isAllowedRune checks if a rune is a valid filename character
func (f *FileSystem) isAllowedRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '-' || r == '_'
}
