package docs

import (
	"github.com/gofiber/fiber/v3"
)

const swaggerHTML = `<!DOCTYPE html>
<html>
  <head>
    <title>TrellGo API — Swagger UI</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css">
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
    <script>
      SwaggerUIBundle({
        url: "/swagger/swagger.json",
        dom_id: "#swagger-ui",
        presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
        layout: "BaseLayout"
      })
    </script>
  </body>
</html>`

const redocHTML = `<!DOCTYPE html>
<html>
  <head>
    <title>TrellGo API Docs</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">
    <style>body { margin: 0; padding: 0; }</style>
  </head>
  <body>
    <redoc spec-url="/docs/swagger.json"></redoc>
    <script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"></script>
  </body>
</html>`

func RegisterRoutes(app *fiber.App) {
	app.Get("/swagger/swagger.json", func(ctx fiber.Ctx) error {
		return ctx.SendFile("./docs/swagger.json")
	})

	app.Get("/swagger", func(ctx fiber.Ctx) error {
		ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
		return ctx.SendString(swaggerHTML)
	})

	app.Get("/docs/swagger.json", func(ctx fiber.Ctx) error {
		return ctx.SendFile("./docs/swagger.json")
	})

	app.Get("/docs", func(ctx fiber.Ctx) error {
		ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
		return ctx.SendString(redocHTML)
	})
}
