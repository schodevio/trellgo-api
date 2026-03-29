package auth

import (
	"testing"
	"time"

	paseto "aidanwoods.dev/go-paseto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- hashToken ---

func TestHashToken_Deterministic(t *testing.T) {
	a := hashToken("some-token")
	b := hashToken("some-token")

	assert.Equal(t, a, b)
}

func TestHashToken_DifferentInputs(t *testing.T) {
	assert.NotEqual(t, hashToken("token-a"), hashToken("token-b"))
}

func TestHashToken_NotEmpty(t *testing.T) {
	assert.NotEmpty(t, hashToken("x"))
}

// --- checkPassword ---

func TestCheckPassword_Correct(t *testing.T) {
	hash, err := hashPassword("secret123")

	require.NoError(t, err)
	assert.True(t, checkPassword("secret123", hash))
}

func TestCheckPassword_Wrong(t *testing.T) {
	hash, err := hashPassword("secret123")

	require.NoError(t, err)
	assert.False(t, checkPassword("wrongpassword", hash))
}

func TestCheckPassword_EmptyPassword(t *testing.T) {
	hash, err := hashPassword("secret123")

	require.NoError(t, err)
	assert.False(t, checkPassword("", hash))
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

func TestParseRefreshToken_Valid(t *testing.T) {
	key := paseto.NewV4SymmetricKey()
	token := makeToken(key, "refresh", time.Now().Add(time.Hour))

	subject, err := parseRefreshToken(token, key)

	assert.NoError(t, err)
	assert.Equal(t, "user-42", subject)
}

func TestParseRefreshToken_AccessTokenRejected(t *testing.T) {
	key := paseto.NewV4SymmetricKey()
	token := makeToken(key, "access", time.Now().Add(time.Hour))

	subject, err := parseRefreshToken(token, key)

	assert.Error(t, err)
	assert.Empty(t, subject)
	assert.ErrorContains(t, err, "invalid token type")
}

func TestParseRefreshToken_Expired(t *testing.T) {
	key := paseto.NewV4SymmetricKey()
	token := makeToken(key, "refresh", time.Now().Add(-time.Minute))

	subject, err := parseRefreshToken(token, key)

	assert.Error(t, err)
	assert.Empty(t, subject)
}

func TestParseRefreshToken_WrongKey(t *testing.T) {
	key := paseto.NewV4SymmetricKey()
	wrongKey := paseto.NewV4SymmetricKey()
	token := makeToken(key, "refresh", time.Now().Add(time.Hour))

	subject, err := parseRefreshToken(token, wrongKey)

	assert.Error(t, err)
	assert.Empty(t, subject)
}

func TestParseRefreshToken_InvalidString(t *testing.T) {
	key := paseto.NewV4SymmetricKey()

	subject, err := parseRefreshToken("not-a-token", key)

	assert.Error(t, err)
	assert.Empty(t, subject)
}
