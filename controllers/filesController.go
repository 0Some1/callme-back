package controllers

import (
	"callme/DTO"
	"callme/lib"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func GetProfileImage(c *fiber.Ctx) error {
	err := c.SendFile("./uploads/profile/" + c.Params("filename"))
	fmt.Println("getProfileImage - :", err)
	if err != nil {
		return fiber.ErrNotFound
	}
	return nil
}

func GetPostImage(c *fiber.Ctx) error {
	err := c.SendFile("./uploads/post/" + c.Params("filename"))
	fmt.Println("getPostImage - :", err)
	if err != nil {
		return fiber.ErrNotFound
	}
	return nil
}

func UploadImageToImageKit(filePath string, fileName string) (string, []error) {
	a := fiber.AcquireAgent()
	req := a.Request()
	req.Header.SetMethod("POST")
	a.SendFile(filePath, "file")
	args := fiber.AcquireArgs()
	args.Set("fileName", fileName)
	a.MultipartForm(args)
	req.SetRequestURI("https://upload.imagekit.io/api/v1/files/upload")
	a.BasicAuth(lib.IMAGEKIT_API_KEY, "")

	if err := a.Parse(); err != nil {
		return "", []error{err}
	}

	_, body, err := a.Bytes()
	if err != nil {
		return "", err
	}
	fiber.ReleaseArgs(args)

	imageKit := new(DTO.ImageKit)
	err2 := json.Unmarshal(body, &imageKit)
	if err2 != nil {
		return "", []error{err2}
	}

	return imageKit.URL, nil

}
