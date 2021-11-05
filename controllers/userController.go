package controllers

import (
	"callme/database"
	"callme/lib"
	"callme/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func GetUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	user.Password = ""
	err := database.DB.PreloadFollowers(user)
	if err != nil {
		fmt.Println("GetUser - ", err)
		c.Status(fiber.ErrInternalServerError.Code).JSON(lib.CustomError(fiber.ErrInternalServerError, "Internal server error"))
	}
	err = database.DB.PreloadFollowings(user)
	if err != nil {
		fmt.Println("GetUser - ", err)
		c.Status(fiber.ErrInternalServerError.Code).JSON(lib.CustomError(fiber.ErrInternalServerError, "Internal server error"))
	}
	var followers int
	var followings int

	followers = len(user.Followers)
	followings = len(user.Followings)
	user.Followers = nil
	user.Followings = nil
	return c.JSON(fiber.Map{
		"name":             user.Name,
		"username":         user.Username,
		"email":            user.Email,
		"followers_count":  followers,
		"followings_count": followings,
		"born":             user.Born,
		"created_at":       user.CreatedAt,
		"bio":              user.Bio,
		"avatar":           user.Avatar,
		"city":             user.City,
		"country":          user.Country,
	})
}
