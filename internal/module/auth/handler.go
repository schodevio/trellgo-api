package auth

import (
	"github.com/gofiber/fiber/v3"
	"github.com/schodevio/trellgo/internal/platform/apierrors"
	"github.com/schodevio/trellgo/internal/platform/validator"
)

type handler struct {
	service Service
}

func newHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) SignUp(ctx fiber.Ctx) error {
	var req SignUpRequest

	if err := ctx.Bind().Body(&req); err != nil {
		return apierrors.BadRequest("invalid request body")
	}

	if err := validator.Validate(&req); err != nil {
		return err
	}

	resp, err := h.service.CreateUser(&req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)
}
