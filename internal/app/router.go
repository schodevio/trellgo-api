package app

import (
	"github.com/gofiber/fiber/v3"
)

func setupRoutes(app *fiber.App, c *Container) {
	api := app.Group("/api/v1")

	c.Auth.RegisterRoutes(api)
	c.Health.RegisterRoutes(api)
}
