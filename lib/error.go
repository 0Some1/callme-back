package lib

import (
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Status      int    `json:"status"`
	Description string `json:"description"`
}

func CustomError(err *fiber.Error, message string) ErrorResponse {

	var customError ErrorResponse
	customError.Status = err.Code
	if message != "" {
		customError.Description = message
	} else {
		customError.Description = err.Message
	}

	return customError
}

func CustomErrorHandler(ctx *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError
	// Retrieve the custom status code if it's an fiber.*Error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	// Send custom error page
	err = ctx.Status(code).JSON(fiber.Map{
		"status":      code,
		"description": err.Error(),
	})
	if err != nil {
		return ctx.Status(code).JSON(fiber.Map{
			"status":      500,
			"description": err.Error(),
		})
	}
	// Return from handler
	return nil
}
