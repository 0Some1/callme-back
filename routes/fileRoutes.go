package routes

import (
	"callme/controllers"
	"github.com/gofiber/fiber/v2"
)

func File(file fiber.Router) {
	file.Get("/profile/:filename", controllers.GetProfileImage)
	file.Get("/post/:filename", controllers.GetPostImage)
}
