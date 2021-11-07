package controllers

import (
	"callme/database"
	"callme/lib"
	"callme/models"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func GetUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
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
func UpdateUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	var userUpdate models.User
	err := c.BodyParser(&userUpdate)
	if err != nil {
		fmt.Println("UpdateUser - ", err)
		return c.Status(fiber.ErrNotAcceptable.Code).JSON(lib.CustomError(fiber.ErrNotAcceptable, "can't parse body"))
	}

	if userUpdate.Name != "" {
		user.Name = userUpdate.Name
	}
	if userUpdate.Username != "" {
		user.Username = userUpdate.Username
	}
	if userUpdate.Email != "" {
		err = validator.New().StructPartial(userUpdate, "Email")
		if err != nil {
			return c.Status(fiber.ErrBadRequest.Code).JSON(lib.CustomError(fiber.ErrBadRequest, err.Error()))
		}
		user.Email = userUpdate.Email
	}
	if userUpdate.Bio != "" {
		user.Bio = userUpdate.Bio
	}
	if userUpdate.City != "" {
		user.City = userUpdate.City
	}
	if userUpdate.Country != "" {
		user.Country = userUpdate.Country
	}
	if userUpdate.Born != nil {
		user.Born = userUpdate.Born
	}
	err = database.DB.SaveUser(user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return c.Status(fiber.ErrConflict.Code).JSON(lib.CustomError(fiber.ErrConflict, "username or email already exists"))
		}
		fmt.Println("UpdateUser - saveUser", err)
		return c.Status(fiber.ErrInternalServerError.Code).JSON(lib.CustomError(fiber.ErrInternalServerError, "Internal server error"))
	}
	return c.Status(204).JSON(nil)
}
func DeleteUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	err := database.DB.DeleteUser(user)
	if err != nil {
		fmt.Println("DeleteUser - ", err)
		return c.Status(fiber.ErrInternalServerError.Code).JSON(lib.CustomError(fiber.ErrInternalServerError, "Internal server error"))
	}
	return c.Status(204).JSON(nil)
}
func UpdateAvatar(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	file, err := c.FormFile("avatar")

	if err != nil {
		fmt.Println("UpdateAvatar - file - ", err)
		return c.Status(fiber.ErrInternalServerError.Code).JSON(lib.CustomError(fiber.ErrInternalServerError, "Internal server error"))
	}

	err = lib.ImageValidation(file)
	if err != nil {
		fmt.Println("UpdateAvatar - image - ", err)
		return c.Status(fiber.ErrBadRequest.Code).JSON(lib.CustomError(fiber.ErrBadRequest, err.Error()))
	}

	file.Filename = lib.AvatarFileName(user.Username, file.Filename)

	err = c.SaveFile(file, fmt.Sprintf("./uploads/profile/%s", file.Filename))
	if err != nil {
		fmt.Println("UpdateAvatar - saveFile - ", err)
		return c.Status(fiber.ErrInternalServerError.Code).JSON(lib.CustomError(fiber.ErrInternalServerError, "Internal server error"))
	}

	user.Avatar = "/uploads/profile/" + file.Filename
	err = database.DB.SaveUser(user)
	if err != nil {
		fmt.Println("UpdateAvatar - saveUser - ", err)
		return c.Status(fiber.ErrInternalServerError.Code).JSON(lib.CustomError(fiber.ErrInternalServerError, "Internal server error"))
	}

	return c.Status(201).JSON(fiber.Map{
		"avatar": c.BaseURL() + user.Avatar,
	})
}
