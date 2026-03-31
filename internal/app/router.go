package app

import (
	"github.com/gofiber/fiber/v3"
	"github.com/schodevio/trellgo/internal/platform/docs"
)

func setupRoutes(app *fiber.App, c *Container) {
	docs.RegisterRoutes(app)

	api := app.Group("/api/v1")

	c.Auth.RegisterRoutes(api)
	c.Boards.RegisterRoutes(api)
	c.Health.RegisterRoutes(api)
	c.Lists.RegisterRoutes(api)
}
