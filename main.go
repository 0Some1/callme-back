package main

import (
	"callme/database"
	"callme/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

)

func main() {
	database.Connect()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	api := app.Group("/api", func(ctx *fiber.Ctx) error {
		return ctx.Next()
	})

	routes.Setup(api)

	app.Listen(":3000")
}