package controllers

import "github.com/gofiber/fiber/v2"

func GetProfileImage(c *fiber.Ctx) error {
	err := c.SendFile("./uploads/profile/" + c.Params("filename"))
	if err != nil {
		return fiber.ErrNotFound
	}
	return nil
}
