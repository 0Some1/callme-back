package controllers

import (
	"callme/database"
	"callme/lib"
	"callme/models"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strings"
)

func GetUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	err := database.DB.PreloadFollowers(user)
	if err != nil {
		fmt.Println("GetUser - ", err)
		return fiber.ErrInternalServerError
	}
	err = database.DB.PreloadFollowings(user)
	if err != nil {
		fmt.Println("GetUser - ", err)
		return fiber.ErrInternalServerError
	}
	var followers int
	var followings int
	followers = len(user.Followers)
	followings = len(user.Followings)
	user.Followers = nil
	user.Followings = nil
	if user.Avatar != "" {
		user.Avatar = c.BaseURL() + user.Avatar
	}

	return c.JSON(fiber.Map{
		"id":               user.ID,
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
func GetUserByID(c *fiber.Ctx) error {
	user, err := database.DB.GetUserByID(c.Params("id"))
	if err != nil {
		fmt.Println("GetUserByID - ", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(404, "User not found")
		} else if strings.Contains(err.Error(), "invalid") {
			return fiber.NewError(400, "invalid user id")
		}
		return fiber.ErrInternalServerError

	}
	err = database.DB.PreloadFollowers(user)
	if err != nil {
		fmt.Println("GetUser - ", err)
		return fiber.ErrInternalServerError
	}
	err = database.DB.PreloadFollowings(user)
	if err != nil {
		fmt.Println("GetUser - ", err)
		return fiber.ErrInternalServerError
	}
	var followers int
	var followings int
	followers = len(user.Followers)
	followings = len(user.Followings)
	user.Followers = nil
	user.Followings = nil
	if user.Avatar != "" {
		user.Avatar = c.BaseURL() + user.Avatar
	}

	localUser := c.Locals("user").(*models.User)
	err = database.DB.PreloadFollowings(localUser)
	if err != nil {
		fmt.Println("GetUser - PreloadFollowings ", err)
		return fiber.ErrInternalServerError
	}
	followingStatus := "not_following"
	if localUser.IsFollowing(user.ID) {
		followingStatus = "following"
	}

	err = database.DB.PreloadRequests(user)
	if err != nil {
		fmt.Println("GetUserByID - GetRequests ", err)
		return fiber.ErrInternalServerError
	}
	if user.IsRequestedByUser(localUser.ID) {
		followingStatus = "requested"
	}

	return c.JSON(fiber.Map{
		"id":               user.ID,
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
		"following_status": followingStatus,
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
			return fiber.NewError(fiber.StatusConflict, "Username or Email already exists")
		}
		fmt.Println("UpdateUser - saveUser", err)
		return fiber.ErrInternalServerError
	}
	return c.Status(204).JSON(nil)
}
func DeleteUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	err := database.DB.DeleteUser(user)
	if err != nil {
		fmt.Println("DeleteUser - ", err)
		return fiber.ErrInternalServerError
	}
	return c.Status(204).JSON(nil)
}
func UpdateAvatar(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	file, err := c.FormFile("avatar")

	if err != nil {
		fmt.Println("UpdateAvatar - file - ", err)
		return fiber.ErrBadRequest
	}

	err = lib.ImageValidation(file)
	if err != nil {
		fmt.Println("UpdateAvatar - image - ", err)
		return c.Status(fiber.ErrUnprocessableEntity.Code).JSON(lib.CustomError(fiber.ErrUnprocessableEntity, err.Error()))
	}

	file.Filename = lib.GenFileName(file.Filename)

	err = c.SaveFile(file, fmt.Sprintf("./uploads/profile/%s", file.Filename))
	if err != nil {
		fmt.Println("UpdateAvatar - saveFile - ", err)
		return fiber.ErrInternalServerError
	}

	user.Avatar = "/uploads/profile/" + file.Filename
	err = database.DB.SaveUser(user)
	if err != nil {
		fmt.Println("UpdateAvatar - saveUser - ", err)
		return fiber.ErrInternalServerError
	}

	return c.Status(201).JSON(fiber.Map{
		"avatar": c.BaseURL() + user.Avatar,
	})
}
func SearchUsers(c *fiber.Ctx) error {
	q := c.Query("q")
	localUser := c.Locals("user").(*models.User)
	users, err := database.DB.SearchUsers(q)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	for i := 0; i < len(users); i++ {
		users[i].PrepareUser(c.BaseURL())
		if users[i].ID == localUser.ID {
			users = append(users[:i], users[i+1:]...)
			i--
		}
	}
	return c.JSON(users)
}
func UnfollowUser(c *fiber.Ctx) error {
	localUser := c.Locals("user").(*models.User)
	otherUser, err := database.DB.GetUserByID(c.Params("id"))
	if err != nil {
		fmt.Println("UnfollowUser - GetUserByID ", err)
		return fiber.NewError(404, "User not found")
	}
	err = database.DB.Unfollow(localUser, otherUser)
	if err != nil {
		fmt.Println("UnfollowUser - Unfollow -", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(204).JSON(nil)
}
func GetFollowers(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	err := database.DB.PreloadFollowers(user)
	if err != nil {
		fmt.Println("GetFollowers - ", err)
		return fiber.ErrInternalServerError
	}
	for i := 0; i < len(user.Followers); i++ {
		user.Followers[i].PrepareUser(c.BaseURL())
	}
	return c.JSON(user.Followers)
}

func GetFollowings(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	err := database.DB.PreloadFollowings(user)
	if err != nil {
		fmt.Println("GetFollowers - ", err)
		return fiber.ErrInternalServerError
	}
	for i := 0; i < len(user.Followings); i++ {
		user.Followings[i].PrepareUser(c.BaseURL())
	}
	return c.JSON(user.Followings)
}
