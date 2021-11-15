package routes

import (
	"callme/controllers"
	"callme/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Setup(api fiber.Router) {
	//authentication routes
	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
	api.Use(middlewares.IsAuthenticated)

	//profile or user routes
	api.Get("/profile", controllers.GetUser)
	api.Patch("/Profile", controllers.UpdateUser)
	api.Delete("/Profile", controllers.DeleteUser)
	api.Put("/avatar", controllers.UpdateAvatar)

	//post routes
	api.Get("/posts", controllers.GetPosts)
	api.Post("/post", controllers.CreatePost)

}
