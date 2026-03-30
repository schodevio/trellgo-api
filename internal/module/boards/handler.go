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

// Create godoc
// @Summary      Create a board
// @Tags         boards
// @Accept       json
// @Produce      json
// @Param        body  body      CreateBoardRequest  true  "Board data"
// @Success      201   {object}  SingleBoardResponse
// @Failure      400   {object}  apierrors.ErrorResponse
// @Failure      401   {object}  apierrors.ErrorResponse
// @Failure      422   {object}  apierrors.ErrorResponse
// @Security     BearerAuth
// @Router       /boards [post]
func (h *handler) Create(ctx fiber.Ctx) error {
	var req CreateBoardRequest

	if err := ctx.Bind().Body(&req); err != nil {
		return apierrors.BadRequest("invalid request body")
	}

	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return apierrors.Unauthorized("missing user identity")
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

// List godoc
// @Summary      List boards
// @Tags         boards
// @Produce      json
// @Success      200  {object}  ListBoardsResponse
// @Failure      401  {object}  apierrors.ErrorResponse
// @Security     BearerAuth
// @Router       /boards [get]
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

// Show godoc
// @Summary      Get a board
// @Tags         boards
// @Produce      json
// @Param        id   path      string  true  "Board ID"
// @Success      200  {object}  SingleBoardResponse
// @Failure      401  {object}  apierrors.ErrorResponse
// @Failure      404  {object}  apierrors.ErrorResponse
// @Security     BearerAuth
// @Router       /boards/{id} [get]
func (h *handler) Show(ctx fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return apierrors.Unauthorized("missing user identity")
	}

	id := ctx.Params("id")

	resp, err := h.service.GetBoard(id, userID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// Update godoc
// @Summary      Update a board
// @Tags         boards
// @Accept       json
// @Produce      json
// @Param        id    path      string              true  "Board ID"
// @Param        body  body      UpdateBoardRequest  true  "Board data"
// @Success      200   {object}  SingleBoardResponse
// @Failure      400   {object}  apierrors.ErrorResponse
// @Failure      401   {object}  apierrors.ErrorResponse
// @Failure      404   {object}  apierrors.ErrorResponse
// @Failure      422   {object}  apierrors.ErrorResponse
// @Security     BearerAuth
// @Router       /boards/{id} [patch]
func (h *handler) Update(ctx fiber.Ctx) error {
	var req UpdateBoardRequest

	if err := ctx.Bind().Body(&req); err != nil {
		return apierrors.BadRequest("invalid request body")
	}

	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return apierrors.Unauthorized("missing user identity")
	}

	req.ID = ctx.Params("id")
	req.UserID = userID

	if err := validator.Validate(&req); err != nil {
		return err
	}

	resp, err := h.service.UpdateBoard(&req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// Delete godoc
// @Summary      Delete a board
// @Tags         boards
// @Produce      json
// @Param        id   path  string  true  "Board ID"
// @Success      204
// @Failure      401  {object}  apierrors.ErrorResponse
// @Failure      404  {object}  apierrors.ErrorResponse
// @Security     BearerAuth
// @Router       /boards/{id} [delete]
func (h *handler) Delete(ctx fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return apierrors.Unauthorized("missing user identity")
	}

	id := ctx.Params("id")

	if err := h.service.DeleteBoard(id, userID); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
