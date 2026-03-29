package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	paseto "aidanwoods.dev/go-paseto"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/schodevio/trellgo/internal/platform/apierrors"
)

func errorHandler(c fiber.Ctx, err error) error {
	var apiErr *apierrors.APIError

	if errors.As(err, &apiErr) {
		return c.Status(apiErr.Status).JSON(fiber.Map{"error": apiErr.Message})
	}

	return c.Status(fiber.StatusInternalServerError).SendString("internal error")
}

func newTestApp(key paseto.V4SymmetricKey) *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: errorHandler})

	app.Get("/protected", Authenticate(key), func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString(c.Locals("user_id").(string))
	})

	return app
}

func makeAccessToken(key paseto.V4SymmetricKey, userID string, expiration time.Time) string {
	now := time.Now()

	t := paseto.NewToken()
	t.SetSubject(userID)
	t.SetIssuedAt(now)
	t.SetNotBefore(now)
	t.SetExpiration(expiration)
	t.SetString("typ", "access")

	return t.V4Encrypt(key, nil)
}

func makeRefreshToken(key paseto.V4SymmetricKey, userID string) string {
	now := time.Now()

	t := paseto.NewToken()
	t.SetSubject(userID)
	t.SetIssuedAt(now)
	t.SetNotBefore(now)
	t.SetExpiration(now.Add(time.Hour))
	t.SetString("typ", "refresh")

	return t.V4Encrypt(key, nil)
}

func doRequest(app *fiber.App, authHeader string) *http.Response {
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	resp, _ := app.Test(req)
	return resp
}

// --- tests ---

func TestAuthenticate_Success(t *testing.T) {
	key := paseto.NewV4SymmetricKey()
	app := newTestApp(key)
	token := makeAccessToken(key, "user-42", time.Now().Add(time.Hour))

	resp := doRequest(app, "Bearer "+token)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuthenticate_SetsUserIDInLocals(t *testing.T) {
	key := paseto.NewV4SymmetricKey()
	app := newTestApp(key)
	token := makeAccessToken(key, "user-42", time.Now().Add(time.Hour))

	resp := doRequest(app, "Bearer "+token)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body := make([]byte, 32)
	n, _ := resp.Body.Read(body)
	assert.Equal(t, "user-42", string(body[:n]))
}

func TestAuthenticate_MissingHeader(t *testing.T) {
	key := paseto.NewV4SymmetricKey()
	app := newTestApp(key)

	resp := doRequest(app, "")

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuthenticate_InvalidHeaderFormat(t *testing.T) {
	key := paseto.NewV4SymmetricKey()
	app := newTestApp(key)

	resp := doRequest(app, "Token some-value")

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuthenticate_MissingBearerPrefix(t *testing.T) {
	key := paseto.NewV4SymmetricKey()
	app := newTestApp(key)
	token := makeAccessToken(key, "user-42", time.Now().Add(time.Hour))

	resp := doRequest(app, token)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuthenticate_ExpiredToken(t *testing.T) {
	key := paseto.NewV4SymmetricKey()
	app := newTestApp(key)
	token := makeAccessToken(key, "user-42", time.Now().Add(-time.Minute))

	resp := doRequest(app, "Bearer "+token)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuthenticate_RefreshTokenRejected(t *testing.T) {
	key := paseto.NewV4SymmetricKey()
	app := newTestApp(key)
	token := makeRefreshToken(key, "user-42")

	resp := doRequest(app, "Bearer "+token)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuthenticate_WrongKey(t *testing.T) {
	key := paseto.NewV4SymmetricKey()
	wrongKey := paseto.NewV4SymmetricKey()
	app := newTestApp(wrongKey)
	token := makeAccessToken(key, "user-42", time.Now().Add(time.Hour))

	resp := doRequest(app, "Bearer "+token)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuthenticate_InvalidToken(t *testing.T) {
	key := paseto.NewV4SymmetricKey()
	app := newTestApp(key)

	resp := doRequest(app, "Bearer not-a-real-token")

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuthenticate_BearerCaseInsensitive(t *testing.T) {
	key := paseto.NewV4SymmetricKey()
	app := newTestApp(key)
	token := makeAccessToken(key, "user-42", time.Now().Add(time.Hour))

	resp := doRequest(app, "BEARER "+token)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
