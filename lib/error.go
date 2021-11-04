package lib

import (
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Status      int    `json:"status"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func CustomError(err *fiber.Error, message string) ErrorResponse {

	var customError ErrorResponse
	customError.Status = err.Code
	customError.Title = err.Error()
	if message != "" {
		customError.Description = message
	} else {
		customError.Description = err.Message
	}

	return customError
}
