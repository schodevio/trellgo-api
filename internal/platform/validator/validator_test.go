package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- test structs ---

type loginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type rangeRequest struct {
	Age   int    `json:"age"   validate:"gte=18,lte=120"`
	Score int    `json:"score" validate:"min=0,max=100"`
	Role  string `json:"role"  validate:"oneof=admin user guest"`
}

type fieldNameRequest struct {
	FirstName string `json:"first_name" validate:"required"`
}

// --- Validate ---

func TestValidate_Valid(t *testing.T) {
	err := Validate(&loginRequest{Email: "user@example.com", Password: "secret123"})
	assert.Nil(t, err)
}

func TestValidate_RequiredField(t *testing.T) {
	err := Validate(&loginRequest{Email: "", Password: "secret123"})

	require.NotNil(t, err)
	assert.Equal(t, 422, err.Status)
	assert.Contains(t, err.Details, "email")
	assert.Contains(t, err.Details["email"], "required")
}

func TestValidate_MultipleErrors(t *testing.T) {
	err := Validate(&loginRequest{Email: "", Password: ""})

	require.NotNil(t, err)
	assert.Contains(t, err.Details, "email")
	assert.Contains(t, err.Details, "password")
}

func TestValidate_InvalidEmail(t *testing.T) {
	err := Validate(&loginRequest{Email: "not-an-email", Password: "secret123"})

	require.NotNil(t, err)
	assert.Contains(t, err.Details, "email")
	assert.Contains(t, err.Details["email"], "email")
}

func TestValidate_MinLength(t *testing.T) {
	err := Validate(&loginRequest{Email: "user@example.com", Password: "short"})

	require.NotNil(t, err)
	assert.Contains(t, err.Details, "password")
	assert.Contains(t, err.Details["password"], "min:8")
}

func TestValidate_UsesJsonTagAsFieldKey(t *testing.T) {
	err := Validate(&fieldNameRequest{FirstName: ""})

	require.NotNil(t, err)
	// key must be "first_name" (json tag), not "FirstName" (Go field name)
	assert.Contains(t, err.Details, "first_name")
	assert.NotContains(t, err.Details, "FirstName")
}

func TestValidate_GteAndLte(t *testing.T) {
	err := Validate(&rangeRequest{Age: 17, Score: 50, Role: "admin"})

	require.NotNil(t, err)
	assert.Contains(t, err.Details, "age")
	assert.Contains(t, err.Details["age"], "gte:18")
}

func TestValidate_MaxValue(t *testing.T) {
	err := Validate(&rangeRequest{Age: 30, Score: 101, Role: "admin"})

	require.NotNil(t, err)
	assert.Contains(t, err.Details, "score")
	assert.Contains(t, err.Details["score"], "max:100")
}

func TestValidate_OneOf(t *testing.T) {
	err := Validate(&rangeRequest{Age: 30, Score: 50, Role: "superadmin"})

	require.NotNil(t, err)
	assert.Contains(t, err.Details, "role")
	assert.Contains(t, err.Details["role"], "oneof:admin user guest")
}

func TestValidate_Returns422(t *testing.T) {
	err := Validate(&loginRequest{})

	require.NotNil(t, err)
	assert.Equal(t, 422, err.Status)
	assert.Equal(t, "unprocessable_entity", err.Code)
}
