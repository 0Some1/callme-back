package middlewares

import (
	"callme/database"
	"callme/lib"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strings"
)

func IsAuthenticated(c *fiber.Ctx) error {
	cookie := c.Get("Authorization")
	if len(strings.Split(cookie, " ")) > 1 {
		cookie = strings.Split(cookie, " ")[1]
	}

	id, err := lib.ParseJwt(cookie)
	if err != nil {
		return c.Status(fiber.ErrUnauthorized.Code).JSON(lib.CustomError(fiber.ErrUnauthorized, "Unauthorized"))
	}

	user, err := database.DB.GetUserByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.ErrNotFound.Code).JSON(lib.CustomError(fiber.ErrNotFound, "User not found"))
	}
	if err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(lib.CustomError(fiber.ErrInternalServerError, "Internal server error"))
	}

	c.Locals("user", user)

	return c.Next()
}
