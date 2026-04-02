package lists

import (
	paseto "aidanwoods.dev/go-paseto"
	"github.com/gofiber/fiber/v3"

	"github.com/schodevio/trellgo/db/sqlc"
	"github.com/schodevio/trellgo/internal/module/boards"
	"github.com/schodevio/trellgo/internal/platform/middleware"
)

type Module struct {
	Service Service
	authKey paseto.V4SymmetricKey
	handler *handler
}

func New(queries *sqlc.Queries, authKey paseto.V4SymmetricKey, boardsService boards.Service) *Module {
	repo := newRepository(queries)
	service := newService(repo, boardsService)
	handler := newHandler(service)

	return &Module{
		Service: service,
		authKey: authKey,
		handler: handler,
	}
}

func (m *Module) RegisterRoutes(router fiber.Router) {
	boardGroup := router.Group("/boards/:board_id/lists", middleware.Authenticate(m.authKey))

	boardGroup.Post("/", m.handler.Create)
	boardGroup.Get("/", m.handler.List)

	group := router.Group("/lists", middleware.Authenticate(m.authKey))

	group.Patch("/:id", m.handler.Update)
}
