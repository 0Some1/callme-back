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

	api.Get("/profile", controllers.GetUser)
	api.Patch("/Profile", controllers.UpdateUser)
	api.Delete("/Profile", controllers.DeleteUser)

	api.Put("/avatar", controllers.UpdateAvatar)

	api.Get("/posts", controllers.GetPosts)

	api.Post("/post", controllers.CreatePost)

}
