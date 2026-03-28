package auth

import (
	"log"

	paseto "aidanwoods.dev/go-paseto"
	"github.com/gofiber/fiber/v3"
	"github.com/schodevio/trellgo/db/sqlc"
)

type Module struct {
	Service Service
	handler *handler
}

func New(queries *sqlc.Queries, secretKey string) *Module {
	authKey, err := paseto.V4SymmetricKeyFromHex(secretKey)
	if err != nil {
		log.Fatalf("invalid SECRET_KEY: %v", err)
	}

	repo := newRepository(queries)
	service := newService(repo, authKey)
	handler := newHandler(service)

	return &Module{
		Service: service,
		handler: handler,
	}
}

func (m *Module) RegisterRoutes(router fiber.Router) {
	group := router.Group("/auth")

	group.Post("/signin", m.handler.SignIn)
	group.Post("/signup", m.handler.SignUp)
}
