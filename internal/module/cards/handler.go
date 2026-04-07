package cards

import (
	"github.com/gofiber/fiber/v3"
	"github.com/schodevio/trellgo/internal/module/shared"
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
// @Summary      Create a card in a list
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
		return apierrors.BadRequest(shared.REQUEST_VALIDATION_FAILED)
	}

	if err := validator.Validate(&req); err != nil {
		return err
	}

	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return apierrors.Unauthorized(shared.USER_AUTHENTICATION_FAILED)
	}

	listID := ctx.Params("list_id")

	resp, err := h.service.CreateCard(listID, userID, &req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

// List godoc
// @Summary      List cards of a list
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
		return apierrors.Unauthorized(shared.USER_AUTHENTICATION_FAILED)
	}

	listID := ctx.Params("list_id")

	resp, err := h.service.ListCards(listID, userID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// Update godoc
// @Summary      Update a card
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
		return apierrors.BadRequest(shared.REQUEST_VALIDATION_FAILED)
	}

	if err := validator.Validate(&req); err != nil {
		return err
	}

	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return apierrors.Unauthorized(shared.USER_AUTHENTICATION_FAILED)
	}

	id := ctx.Params("id")

	resp, err := h.service.UpdateCard(id, userID, &req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// Delete godoc
// @Summary      Delete a card
// @Tags         cards
// @Produce      json
// @Param        id   path      string  true  "Card ID"
// @Success      204
// @Failure      401  {object}  apierrors.ErrorResponse
// @Failure      404  {object}  apierrors.ErrorResponse
// @Security     BearerAuth
// @Router       /cards/{id} [delete]
func (h *handler) Delete(ctx fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return apierrors.Unauthorized(shared.USER_AUTHENTICATION_FAILED)
	}

	id := ctx.Params("id")

	if err := h.service.DeleteCard(id, userID); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
