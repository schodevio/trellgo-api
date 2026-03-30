package boards

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

func (h *handler) Create(ctx fiber.Ctx) error {
	var req CreateBoardRequest

	if err := ctx.Bind().Body(&req); err != nil {
		return apierrors.BadRequest("invalid request body")
	}

	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return apierrors.NotFound("user not found")
	}

	// Assign current user ID to the requests
	req.UserID = userID

	if err := validator.Validate(&req); err != nil {
		return err
	}

	resp, err := h.service.CreateBoard(&req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

func (h *handler) List(ctx fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return apierrors.Unauthorized("missing user identity")
	}

	resp, err := h.service.ListBoards(userID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}
