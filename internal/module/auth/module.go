package auth

import (
	paseto "aidanwoods.dev/go-paseto"
	"github.com/gofiber/fiber/v3"
	"github.com/schodevio/trellgo/db/sqlc"
)

type Module struct {
	Service Service
	AuthKey paseto.V4SymmetricKey
	handler *handler
}

func New(queries *sqlc.Queries, authKey paseto.V4SymmetricKey) *Module {
	repo := newRepository(queries)
	service := newService(repo, authKey)
	handler := newHandler(service)

	return &Module{
		Service: service,
		AuthKey: authKey,
		handler: handler,
	}
}

func (m *Module) RegisterRoutes(router fiber.Router) {
	group := router.Group("/auth")

	group.Post("/signin", m.handler.SignIn)
	group.Post("/signup", m.handler.SignUp)
}
