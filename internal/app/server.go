package app

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	"github.com/schodevio/trellgo/internal/platform/config"
	"github.com/schodevio/trellgo/internal/platform/middleware"
)

type Server struct {
	app       *fiber.App
	container *Container
	cfg       *config.Config
}

func NewServer(cfg *config.Config) *Server {
	ctx := context.Background()

	app := fiber.New(fiber.Config{
		AppName:        "TrellGo",
		BodyLimit:      1 * 1024 * 1024, // 1MB
		ReadBufferSize: 16 * 1024,       // 16KB — handles large browser headers
		ErrorHandler:   ErrorHandler,
	})

	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())

	container := NewContainer(ctx, cfg)
	setupRoutes(app, container)

	return &Server{
		app:       app,
		container: container,
		cfg:       cfg,
	}
}

func (s *Server) Start() {
	addr := fmt.Sprintf(":%s", s.cfg.Port)
	log.Printf("Starting TrellGo on %s", addr)

	if err := s.app.Listen(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
