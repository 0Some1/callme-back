package controllers

import (
	"callme/database"
	"callme/lib"
	"callme/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func GetPosts(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	err := database.DB.PreloadPosts(user)
	if err != nil {
		fmt.Println("GetPosts - ", err)
		c.Status(fiber.ErrInternalServerError.Code).JSON(lib.CustomError(fiber.ErrInternalServerError, "Internal server error"))
	}
	return c.JSON(user.Posts)
}
