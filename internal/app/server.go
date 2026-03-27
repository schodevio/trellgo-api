package app

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/schodevio/trellgo/internal/config"
)

type Server struct {
	app       *fiber.App
	container *Container
	cfg       *config.Config
}

func NewServer(cfg *config.Config) *Server {
	app := fiber.New(fiber.Config{
		AppName:   "TrellGo",
		BodyLimit: 1 * 1024 * 1024, // 1MB
	})

	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(loggerMiddleware())
	app.Use(corsMiddleware())

	container := NewContainer(cfg)

	server := &Server{
		app:       app,
		container: container,
		cfg:       cfg,
	}

	setupRoutes(app, container)

	return server
}

func (s *Server) Start() {
	addr := fmt.Sprintf(":%s", s.cfg.Port)
	log.Printf("Starting TrellGo on %s", addr)

	if err := s.app.Listen(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
