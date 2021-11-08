package controllers

import (
	"callme/database"
	"callme/lib"
	"callme/models"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func GetPosts(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	err := database.DB.PreloadPosts(user)
	if err != nil {
		fmt.Println("GetPosts - ", err)
		return fiber.ErrInternalServerError
	}
	return c.JSON(user.Posts)
}
func CreatePost(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	post := new(models.Post)
	err := c.BodyParser(post)
	if err != nil {
		fmt.Println("CreatePost - Bodyparser - ", err)
		return fiber.ErrBadRequest
	}
	//post validation must be done here!!!
	err = validator.New().Struct(post)
	if err != nil {
		fmt.Println("CreatePost - Validation - ", err)
		return fiber.ErrBadRequest
	}

	post.UserID = user.ID

	var photos []models.Photo

	if form, err := c.MultipartForm(); err == nil {
		files := form.File["photos"]
		for _, file := range files {
			err = lib.ImageValidation(file)
			if err != nil {
				fmt.Println("CreatePost - imageValidation -", err)
				return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
			}
			file.Filename = lib.GenFileName(post.Title, file.Filename)
			err = c.SaveFile(file, fmt.Sprintf("./uploads/post/%s", file.Filename))
			if err != nil {
				fmt.Println("CreatePost - saveFile ", err)
				return fiber.ErrInternalServerError
			}
			photoTemp := models.Photo{
				Name: file.Filename,
				Path: "/uploads/post/" + file.Filename,
			}
			photos = append(photos, photoTemp)
			err := database.DB.CreatePhoto(&photoTemp)
			if err != nil {
				fmt.Println("CreatePost - createPhoto -", err)
				return fiber.ErrInternalServerError
			}

		}

	}
	post.Photos = photos
	err = database.DB.CreatePost(post)
	if err != nil {
		fmt.Println("CreatePost - SavePost -", err)
		return fiber.ErrInternalServerError
	}
	return c.JSON(post)
}
