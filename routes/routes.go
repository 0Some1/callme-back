package routes

import (
	"callme/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app fiber.Router) {
	app.Post("/register", controllers.Register)
	app.Post("/login",controllers.Login )
	app.Get("/user",controllers.AuthenticatedUser)
	app.Get("/logout",controllers.Logout)

}