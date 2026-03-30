package boards

import (
	paseto "aidanwoods.dev/go-paseto"
	"github.com/gofiber/fiber/v3"
	"github.com/schodevio/trellgo/db/sqlc"
	"github.com/schodevio/trellgo/internal/platform/middleware"
)

type Module struct {
	Service Service
	authKey paseto.V4SymmetricKey
	handler *handler
}

func New(queries *sqlc.Queries, authKey paseto.V4SymmetricKey) *Module {
	repo := newRepository(queries)
	service := newService(repo)
	handler := newHandler(service)

	return &Module{
		Service: service,
		authKey: authKey,
		handler: handler,
	}
}

func (m *Module) RegisterRoutes(router fiber.Router) {
	group := router.Group("/boards", middleware.Authenticate(m.authKey))

	group.Post("/", m.handler.Create)
	group.Get("/", m.handler.List)
}
