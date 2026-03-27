package app

import "github.com/gofiber/fiber/v3"

func setupRoutes(app *fiber.App, cntr *Container) {
	app.Get("/health", func(ctx fiber.Ctx) error {
		return ctx.JSON(fiber.Map{"status": "ok"})
	})

	api := app.Group("/api/v1")
	cntr.Auth.RegisterRoutes(api)
}
