package health

import (
	"github.com/gofiber/fiber/v3"
)

type handler struct {
	service Service
}

func newHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) HealthCheck(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok", "protected": false})
}

func (h *handler) ProtectedCheck(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok", "protected": true})
}
