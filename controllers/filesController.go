package controllers

import (
	"callme/database"
	"callme/models"
	"github.com/gofiber/fiber/v2"
)

func GetProfileImage(c *fiber.Ctx) error {
	err := c.SendFile("./uploads/profile/" + c.Params("filename"))
	if err != nil {
		return fiber.ErrNotFound
	}
	return nil
}

func GetPostImage(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	post := new(models.Post)
	filename := c.Params("filename")
	post, err := database.DB.GetPostByPhotoName(filename)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Post not found")
	}
	if post.UserID == user.ID {
		err = c.SendFile("./uploads/post/" + filename)
		if err != nil {
			return fiber.ErrNotFound
		}
		return nil
	}

	err = database.DB.PreloadFollowings(user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}
	//it is not tested , and it might have issues
	if user.IsFollowing(post.UserID) {
		err = c.SendFile("./uploads/post/" + filename)
		if err != nil {
			return fiber.ErrNotFound
		}
		return nil
	}

	return fiber.NewError(fiber.StatusForbidden, "You are not allowed to see this post")

}
