package lib

import (
	"errors"
	"mime/multipart"
	"strings"
	"time"
)

func ImageValidation(file *multipart.FileHeader) error {
	fileHeader := file.Header
	fileType := fileHeader.Get("Content-Type")
	if !strings.HasPrefix(fileType, "image/") {
		return errors.New("file is not an image")
	}
	return nil
}

func GenFileName(identifier string, fileName string) string {
	temp := strings.Split(fileName, ".")
	if len(temp) < 2 {
		return time.Now().Format("2006-01-02_15:04:05") + "_" + identifier + "_" + fileName
	}
	t := temp[len(temp)-1]
	return time.Now().Format("2006-01-02_15:04:05") + "_" + identifier + "." + t
}
