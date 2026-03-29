package apierrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- constructors ---

func TestBadRequest(t *testing.T) {
	err := BadRequest("invalid body")
	assert.Equal(t, 400, err.Status)
	assert.Equal(t, "bad_request", err.Code)
	assert.Equal(t, "invalid body", err.Message)
}

func TestUnauthorized(t *testing.T) {
	err := Unauthorized("missing token")
	assert.Equal(t, 401, err.Status)
	assert.Equal(t, "unauthorized", err.Code)
	assert.Equal(t, "missing token", err.Message)
}

func TestForbidden(t *testing.T) {
	err := Forbidden("access denied")
	assert.Equal(t, 403, err.Status)
	assert.Equal(t, "forbidden", err.Code)
	assert.Equal(t, "access denied", err.Message)
}

func TestNotFound(t *testing.T) {
	err := NotFound("resource not found")
	assert.Equal(t, 404, err.Status)
	assert.Equal(t, "not_found", err.Code)
	assert.Equal(t, "resource not found", err.Message)
}

func TestConflict(t *testing.T) {
	err := Conflict("already exists")
	assert.Equal(t, 409, err.Status)
	assert.Equal(t, "conflict", err.Code)
	assert.Equal(t, "already exists", err.Message)
}

func TestUnprocessableEntity(t *testing.T) {
	details := ErrorDetails{"email": {"required", "email"}}
	err := UnprocessableEntity("validation failed", details)

	assert.Equal(t, 422, err.Status)
	assert.Equal(t, "unprocessable_entity", err.Code)
	assert.Equal(t, "validation failed", err.Message)
	assert.Equal(t, details, err.Details)
}

func TestUnprocessableEntity_NilDetails(t *testing.T) {
	err := UnprocessableEntity("validation failed", nil)
	assert.Equal(t, 422, err.Status)
	assert.Nil(t, err.Details)
}

func TestInternal(t *testing.T) {
	err := Internal("something went wrong")
	assert.Equal(t, 500, err.Status)
	assert.Equal(t, "internal_error", err.Code)
	assert.Equal(t, "something went wrong", err.Message)
}

// --- Error() ---

func TestAPIError_ErrorString(t *testing.T) {
	err := BadRequest("invalid body")
	assert.Equal(t, "(400) bad_request: invalid body", err.Error())
}

func TestAPIError_ErrorString_Internal(t *testing.T) {
	err := Internal("db failure")
	assert.Equal(t, "(500) internal_error: db failure", err.Error())
}

// --- implements error interface ---

func TestAPIError_ImplementsError(t *testing.T) {
	var err error = BadRequest("test")
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "bad_request")
}
