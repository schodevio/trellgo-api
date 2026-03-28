package health

import (
	paseto "aidanwoods.dev/go-paseto"
	"github.com/gofiber/fiber/v3"
	"github.com/schodevio/trellgo/internal/platform/middleware"
)

type Module struct {
	Service Service
	AuthKey paseto.V4SymmetricKey
	handler *handler
}

func New(authKey paseto.V4SymmetricKey) *Module {
	service := newService(authKey)
	handler := newHandler(service)

	return &Module{
		Service: service,
		AuthKey: authKey,
		handler: handler,
	}
}

func (m *Module) RegisterRoutes(router fiber.Router) {
	group := router.Group("/health")

	group.Get("/", m.handler.HealthCheck)

	protected := group.Group("", middleware.Authenticate(m.AuthKey))
	protected.Get("/protected", m.handler.ProtectedCheck)
}
