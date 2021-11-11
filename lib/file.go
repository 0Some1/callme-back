package lib

import (
	"errors"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"mime/multipart"
	"strings"
)

func ImageValidation(file *multipart.FileHeader) error {
	fileHeader := file.Header
	fileType := fileHeader.Get("Content-Type")
	if !strings.HasPrefix(fileType, "image/") {
		return errors.New("file is not an image")
	}
	//check size of file
	if file.Size > 5000000 {
		return errors.New("file is too large")
	}
	return nil
}

func GenFileName(fileName string) string {
	temp := strings.Split(fileName, ".")
	fileExtention := ""
	if len(temp) > 1 {
		fileExtention = temp[len(temp)-1]
	}
	return uuid.New().String() + "." + fileExtention
}
