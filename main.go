package main

import (
	"callme/database"
	"callme/lib"
	"callme/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"os"
)

func main() {
	database.Connect()
	app := fiber.New(fiber.Config{
		ErrorHandler: lib.CustomErrorHandler,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: true,
	}))

	app.Use(logger.New())

	api := app.Group("/api", func(ctx *fiber.Ctx) error {
		return ctx.Next()
	})

	file := app.Group("/uploads", func(ctx *fiber.Ctx) error {
		return ctx.Next()
	})

	routes.Setup(api)

	routes.File(file)

	app.All("*", func(c *fiber.Ctx) error {
		return c.Status(404).JSON(lib.CustomError(fiber.ErrNotFound, "Not Found"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	app.Listen(":" + port)

}
