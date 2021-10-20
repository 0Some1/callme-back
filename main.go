package main

import (
	"callme/database"
	"callme/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"os"
)

func main() {
	database.Connect()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: true,
	}))

	api := app.Group("/api", func(ctx *fiber.Ctx) error {
		return ctx.Next()
	})

	routes.Setup(api)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	app.Listen(":" + port)
}
