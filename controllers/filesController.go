package controllers

import (
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
