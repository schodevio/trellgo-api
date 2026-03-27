package app

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/schodevio/trellgo/internal/platform/apierrors"
)

// Handler

func ErrorHandler(c fiber.Ctx, err error) error {
	var apiErr *apierrors.APIError

	// 1. Custom error
	if errors.As(err, &apiErr) {
		return c.Status(apiErr.Status).JSON(apierrors.ErrorResponse{Error: *apiErr})
	}

	var fiberErr *fiber.Error

	// 2. Fiber error
	if errors.As(err, &fiberErr) {
		return c.Status(fiberErr.Code).JSON(apierrors.ErrorResponse{
			Error: apierrors.APIError{
				Status:  fiberErr.Code,
				Code:    "http_error",
				Message: fiberErr.Message,
			},
		})
	}

	// 3. Error fallback
	return c.Status(fiber.StatusInternalServerError).JSON(apierrors.ErrorResponse{
		Error: apierrors.APIError{
			Status:  fiber.StatusInternalServerError,
			Code:    "internal_error",
			Message: "An unexpected error occurred.",
		},
	})
}
