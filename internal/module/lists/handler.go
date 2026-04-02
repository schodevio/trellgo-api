package lists

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
// @Summary      Create a board list
// @Tags         lists
// @Accept       json
// @Produce      json
// @Param        board_id  path      string  true  "Board ID"
// @Param        body  		 body      CreateListRequest  true  "List data"
// @Success      201   		 {object}  SingleListResponse
// @Failure      400   		 {object}  apierrors.ErrorResponse
// @Failure      401   		 {object}  apierrors.ErrorResponse
// @Failure      422   		 {object}  apierrors.ErrorResponse
// @Security     BearerAuth
// @Router       /boards/{board_id}/lists [post]
func (h *handler) Create(ctx fiber.Ctx) error {
	var req CreateListRequest

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

	boardID := ctx.Params("board_id")

	resp, err := h.service.CreateList(boardID, userID, &req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

// List godoc
// @Summary      List board lists
// @Tags         lists
// @Produce      json
// @Param        board_id  path      string  true  "Board ID"
// @Produce      json
// @Success      200  {object}  ListListsResponse
// @Failure      401  {object}  apierrors.ErrorResponse
// @Security     BearerAuth
// @Router       /boards/{board_id}/lists [get]
func (h *handler) List(ctx fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return apierrors.Unauthorized("missing user identity")
	}

	boardID := ctx.Params("board_id")

	resp, err := h.service.ListLists(boardID, userID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// Update godoc
// @Summary      Update a list
// @Tags         lists
// @Accept       json
// @Produce      json
// @Param        id    path      string             true  "List ID"
// @Param        body  body      UpdateListRequest  true  "List data"
// @Success      200   {object}  SingleListResponse
// @Failure      400   {object}  apierrors.ErrorResponse
// @Failure      401   {object}  apierrors.ErrorResponse
// @Failure      404   {object}  apierrors.ErrorResponse
// @Failure      422   {object}  apierrors.ErrorResponse
// @Security     BearerAuth
// @Router       /lists/{id} [patch]
func (h *handler) Update(ctx fiber.Ctx) error {
	var req UpdateListRequest

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

	resp, err := h.service.UpdateList(id, userID, &req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}
