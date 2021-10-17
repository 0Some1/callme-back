package routes

import (
	"callme/controllers"
	"callme/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Setup(app fiber.Router) {
	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)
	app.Use(middlewares.IsAuthenticated)
	app.Get("/user", controllers.AuthenticatedUser)
	app.Get("/logout", controllers.Logout)

}
