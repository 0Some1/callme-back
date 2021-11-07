package routes

import (
	"callme/controllers"
	"callme/middlewares"
	"github.com/gofiber/fiber/v2"
)

func File(file fiber.Router) {
	file.Use(middlewares.IsAuthenticated)
	file.Get("/profile/:filename", controllers.GetProfileImage)
}
