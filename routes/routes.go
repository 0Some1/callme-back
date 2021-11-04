package routes

import (
	"callme/controllers"
	"callme/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Setup(api fiber.Router) {
	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
	api.Use(middlewares.IsAuthenticated)

}
