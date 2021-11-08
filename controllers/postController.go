package controllers

import (
	"callme/database"
	"callme/models"
	"fmt"
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
		fmt.Println("CreatePost - ", err)
		return fiber.ErrBadRequest
	}
	post.UserID = user.ID

	var photos []models.Photo

	if form, err := c.MultipartForm(); err == nil {
		files := form.File["photos"]
		for _, file := range files {

			photoTemp := models.Photo{
				Name: file.Filename,
				Path: "/uploads/post" + file.Filename,
			}
			photos = append(photos, photoTemp)
			err := database.DB.CreatePhoto(&photoTemp)
			if err != nil {
				return fiber.ErrInternalServerError
			}

		}

	}
	//err = database.DB.CreatePost(post)
	//if err != nil {
	//    fmt.Println("CreatePost - ", err)
	//    return fiber.ErrInternalServerError
	//}
	return c.JSON(post)
}
