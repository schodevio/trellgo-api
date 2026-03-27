package apierrors

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
)

// Structs

type ErrorResponse struct {
	Error APIError `json:"error"`
}

type APIError struct {
	Status  int          `json:"-"`
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Details ErrorDetails `json:"details,omitempty"`
}

type ErrorDetails map[string][]string

func New(status int, code, message string) *APIError {
	return &APIError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

func (e *APIError) Error() string {
	return fmt.Sprintf("(%d) %s: %s", e.Status, e.Code, e.Message)
}

// Errors

func BadRequest(message string) *APIError {
	return New(fiber.StatusBadRequest, "bad_request", message)
}

func Unauthorized(message string) *APIError {
	return New(fiber.StatusUnauthorized, "unauthorized", message)
}

func Forbidden(message string) *APIError {
	return New(fiber.StatusForbidden, "forbidden", message)
}

func NotFound(message string) *APIError {
	return New(fiber.StatusNotFound, "not_found", message)
}

func Conflict(message string) *APIError {
	return New(fiber.StatusConflict, "conflict", message)
}

func UnprocessableEntity(message string, details ErrorDetails) *APIError {
	return &APIError{
		Status:  fiber.StatusUnprocessableEntity,
		Code:    "unprocessable_entity",
		Message: message,
		Details: details,
	}
}

func Internal(message string) *APIError {
	return New(fiber.StatusInternalServerError, "internal_error", message)
}
