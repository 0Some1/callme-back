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
	api.Get("/profile/:id", controllers.GetUserByID)
	api.Patch("/Profile", controllers.UpdateUser)
	api.Delete("/Profile", controllers.DeleteUser)
	api.Put("/avatar", controllers.UpdateAvatar)
	api.Get("/search", controllers.SearchUsers)
	api.Get("/unfollow/:id", controllers.UnfollowUser)
	api.Get("/followers", controllers.GetFollowers)
	api.Get("/followings", controllers.GetFollowings)

	//post routes
	api.Get("/posts", controllers.GetPosts)
	api.Post("/post", controllers.CreatePost)
	api.Delete("/post/:postID", controllers.DeletePost)
	api.Get("/posts/:userID", controllers.GetPostsByUserID)

	//request routes
	api.Get("/requests", controllers.GetRequests)
	api.Post("/request/:id", controllers.CreateRequest)
	api.Delete("/request/:userID", controllers.DeleteRequest)
	api.Get("/request/:id/accept", controllers.AcceptRequest)
	api.Get("/request/:id/decline", controllers.DeclineRequest)

}
