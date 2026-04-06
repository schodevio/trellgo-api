package cards

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

// Create godoc
// @Summary      Create a list card
// @Tags         cards
// @Accept       json
// @Produce      json
// @Param        list_id   path      string  true  "List ID"
// @Param        body  		 body      CreateCardRequest  true  "Card data"
// @Success      201   		 {object}  SingleCardResponse
// @Failure      400   		 {object}  apierrors.ErrorResponse
// @Failure      401   		 {object}  apierrors.ErrorResponse
// @Failure      422   		 {object}  apierrors.ErrorResponse
// @Security     BearerAuth
// @Router       /lists/{list_id}/cards [post]
func (h *handler) Create(ctx fiber.Ctx) error {
	var req CreateCardRequest

	if err := ctx.Bind().Body(&req); err != nil {
		return apierrors.BadRequest("invalid request body")
	}

	if err := validator.Validate(&req); err != nil {
		return err
	}

	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return apierrors.Unauthorized("missing user identity")
	}

	listID := ctx.Params("list_id")

	resp, err := h.service.CreateCard(listID, userID, &req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

// List godoc
// @Summary      List list cards
// @Tags         cards
// @Produce      json
// @Param        list_id  path      string  true  "List ID"
// @Produce      json
// @Success      200  {object}  ListCardsResponse
// @Failure      401  {object}  apierrors.ErrorResponse
// @Security     BearerAuth
// @Router       /lists/{list_id}/cards [get]
func (h *handler) List(ctx fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return apierrors.Unauthorized("missing user identity")
	}

	listID := ctx.Params("list_id")

	resp, err := h.service.ListCards(listID, userID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// Update godoc
// @Summary      Update a list card
// @Tags         cards
// @Accept       json
// @Produce      json
// @Param        id        path      string  true  "Card ID"
// @Param        body  		 body      UpdateCardRequest  true  "Card data"
// @Success      200   		 {object}  SingleCardResponse
// @Failure      400   		 {object}  apierrors.ErrorResponse
// @Failure      401   		 {object}  apierrors.ErrorResponse
// @Failure      422   		 {object}  apierrors.ErrorResponse
// @Security     BearerAuth
// @Router       /cards/{id} [put]
func (h *handler) Update(ctx fiber.Ctx) error {
	var req UpdateCardRequest

	if err := ctx.Bind().Body(&req); err != nil {
		return apierrors.BadRequest("invalid request body")
	}

	if err := validator.Validate(&req); err != nil {
		return err
	}

	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return apierrors.Unauthorized("missing user identity")
	}

	id := ctx.Params("id")

	resp, err := h.service.UpdateCard(id, userID, &req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}
