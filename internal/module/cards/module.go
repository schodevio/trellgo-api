package cards

import (
	paseto "aidanwoods.dev/go-paseto"
	"github.com/gofiber/fiber/v3"

	"github.com/schodevio/trellgo/db/sqlc"
	"github.com/schodevio/trellgo/internal/module/lists"
	"github.com/schodevio/trellgo/internal/platform/middleware"
)

type Module struct {
	Service Service
	authKey paseto.V4SymmetricKey
	handler *handler
}

func New(queries *sqlc.Queries, authKey paseto.V4SymmetricKey, listsService lists.Service) *Module {
	repo := newRepository(queries)
	service := newService(repo, listsService)
	handler := newHandler(service)

	return &Module{
		Service: service,
		authKey: authKey,
		handler: handler,
	}
}

func (m *Module) RegisterRoutes(router fiber.Router) {
	boardGroup := router.Group("/lists/:list_id/cards", middleware.Authenticate(m.authKey))

	boardGroup.Post("/", m.handler.Create)
	boardGroup.Get("/", m.handler.List)
}
