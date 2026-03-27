package auth

import (
	"github.com/gofiber/fiber/v3"
	"github.com/schodevio/trellgo/internal/platform/apierrors"
)

type handler struct {
	service Service
}

func newHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) SignUp(ctx fiber.Ctx) error {
	var req CreateUserRequest

	if err := ctx.Bind().Body(&req); err != nil {
		return apierrors.BadRequest("invalid request body")
	}

	resp, err := h.service.CreateUser(&req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)
}
