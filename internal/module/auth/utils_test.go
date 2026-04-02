package auth

import (
	"testing"
	"time"

	paseto "aidanwoods.dev/go-paseto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- hashToken ---

func TestHashToken(t *testing.T) {
	t.Run("deterministic", func(t *testing.T) {
		a := hashToken("some-token")
		b := hashToken("some-token")

		assert.Equal(t, a, b)
	})

	t.Run("differentTokens", func(t *testing.T) {
		assert.NotEqual(t, hashToken("token-a"), hashToken("token-b"))
	})

	t.Run("not empty", func(t *testing.T) {
		assert.NotEmpty(t, hashToken("some-token"))
	})
}

// --- checkPassword ---

func TestCheckPassword(t *testing.T) {
	t.Run("correct password", func(t *testing.T) {
		hash, err := hashPassword("secret123")

		require.NoError(t, err)
		assert.True(t, checkPassword("secret123", hash))

	})

	t.Run("wrong password", func(t *testing.T) {
		hash, err := hashPassword("secret123")

		require.NoError(t, err)
		assert.False(t, checkPassword("wrongpassword", hash))
	})

	t.Run("empty password", func(t *testing.T) {
		hash, err := hashPassword("secret123")

		require.NoError(t, err)
		assert.False(t, checkPassword("", hash))
	})
}

// --- parseRefreshToken ---

func makeToken(key paseto.V4SymmetricKey, typ string, expiration time.Time) string {
	now := time.Now()

	t := paseto.NewToken()
	t.SetSubject("user-42")
	t.SetIssuedAt(now)
	t.SetNotBefore(now)
	t.SetExpiration(expiration)
	t.SetString("typ", typ)

	return t.V4Encrypt(key, nil)
}

func TestParseREfreshToken(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		key := paseto.NewV4SymmetricKey()
		token := makeToken(key, "refresh", time.Now().Add(time.Hour))

		subject, err := parseRefreshToken(token, key)

		assert.NoError(t, err)
		assert.Equal(t, "user-42", subject)
	})

	t.Run("access token rejected", func(t *testing.T) {
		key := paseto.NewV4SymmetricKey()
		token := makeToken(key, "access", time.Now().Add(time.Hour))

		subject, err := parseRefreshToken(token, key)

		assert.Error(t, err)
		assert.Empty(t, subject)
		assert.ErrorContains(t, err, "invalid token type")
	})

	t.Run("expired token", func(t *testing.T) {
		key := paseto.NewV4SymmetricKey()
		token := makeToken(key, "refresh", time.Now().Add(-time.Minute))

		subject, err := parseRefreshToken(token, key)

		assert.Error(t, err)
		assert.Empty(t, subject)
		assert.ErrorContains(t, err, "this token has expired")
	})

	t.Run("wrong key", func(t *testing.T) {
		key := paseto.NewV4SymmetricKey()
		wrongKey := paseto.NewV4SymmetricKey()
		token := makeToken(key, "refresh", time.Now().Add(time.Hour))

		subject, err := parseRefreshToken(token, wrongKey)

		assert.Error(t, err)
		assert.Empty(t, subject)
	})

	t.Run("invalid token string", func(t *testing.T) {
		key := paseto.NewV4SymmetricKey()

		subject, err := parseRefreshToken("not-a-token", key)

		assert.Error(t, err)
		assert.Empty(t, subject)
	})
}
