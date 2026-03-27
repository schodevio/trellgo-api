package auth

import (
	"github.com/gofiber/fiber/v3"
	"github.com/schodevio/trellgo/db/sqlc"
)

type Module struct {
	Service Service
	handler *handler
}

func New(queries *sqlc.Queries) *Module {
	repo := newRepository(queries)
	service := newService(repo)
	handler := newHandler(service)

	return &Module{
		Service: service,
		handler: handler,
	}
}

func (m *Module) RegisterRoutes(router fiber.Router) {
	group := router.Group("/auth")

	group.Post("/signup", m.handler.SignUp)
}
