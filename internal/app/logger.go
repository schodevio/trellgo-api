package app

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func loggerMiddleware() fiber.Handler {
	return logger.New(logger.Config{
		Format:     logger.DefaultFormat,
		TimeFormat: time.RFC3339Nano,
		TimeZone:   "UTC",
	})
}
