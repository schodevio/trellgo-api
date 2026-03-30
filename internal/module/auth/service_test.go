package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	paseto "aidanwoods.dev/go-paseto"
	"github.com/schodevio/trellgo/db/sqlc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- mock ---

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) CreateRefreshToken(ctx context.Context, userID, token, userAgent, ipAddress string, expiresAt time.Time) (sqlc.RefreshToken, error) {
	args := m.Called(ctx, userID, token, userAgent, ipAddress, expiresAt)
	return args.Get(0).(sqlc.RefreshToken), args.Error(1)
}

func (m *mockRepository) GetRefreshTokenByRawToken(ctx context.Context, token string) (sqlc.RefreshToken, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(sqlc.RefreshToken), args.Error(1)
}

func (m *mockRepository) RevokeRefreshToken(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockRepository) CreateUser(ctx context.Context, email, passwordHash string) (sqlc.User, error) {
	args := m.Called(ctx, email, passwordHash)
	return args.Get(0).(sqlc.User), args.Error(1)
}

func (m *mockRepository) GetUserByEmail(ctx context.Context, email string) (sqlc.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(sqlc.User), args.Error(1)
}

// --- helpers ---

func newTestService(repo Repository) Service {
	return newService(repo, paseto.NewV4SymmetricKey())
}

func newTestServiceWithKey(repo Repository) (Service, paseto.V4SymmetricKey) {
	key := paseto.NewV4SymmetricKey()
	return newService(repo, key), key
}

func validRefreshToken(key paseto.V4SymmetricKey, userID string) string {
	now := time.Now()

	t := paseto.NewToken()
	t.SetSubject(userID)
	t.SetIssuedAt(now)
	t.SetNotBefore(now)
	t.SetExpiration(now.Add(refreshTokenTTL))
	t.SetString("typ", "refresh")

	return t.V4Encrypt(key, nil)
}

func hashedTestPassword() string {
	h, _ := hashPassword("password123")
	return h
}

// --- tests ---

func TestSignInUser_Success(t *testing.T) {
	repo := new(mockRepository)
	svc := newTestService(repo)

	user := sqlc.User{
		ID:           "user-1",
		Email:        "test@example.com",
		PasswordHash: hashedTestPassword(),
	}

	repo.
		On("GetUserByEmail", mock.Anything, "test@example.com").
		Return(user, nil)

	repo.
		On("CreateRefreshToken", mock.Anything, user.ID, mock.AnythingOfType("string"), "", "", mock.AnythingOfType("time.Time")).
		Return(sqlc.RefreshToken{}, nil)

	resp, err := svc.SignInUser(&SignInRequest{
		Email:    "test@example.com",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	repo.AssertExpectations(t)
}

func TestSignInUser_UserNotFound(t *testing.T) {
	repo := new(mockRepository)
	svc := newTestService(repo)

	repo.
		On("GetUserByEmail", mock.Anything, "notfound@example.com").
		Return(sqlc.User{}, errors.New("not found"))

	resp, err := svc.SignInUser(&SignInRequest{
		Email:    "notfound@example.com",
		Password: "password123",
	})

	assert.Error(t, err)
	assert.Empty(t, resp.AccessToken)
	assert.ErrorContains(t, err, "invalid credentials")
	repo.AssertExpectations(t)
}

func TestSignInUser_WrongPassword(t *testing.T) {
	repo := new(mockRepository)
	svc := newTestService(repo)

	user := sqlc.User{
		ID:           "user-1",
		Email:        "test@example.com",
		PasswordHash: hashedTestPassword(),
	}

	repo.
		On("GetUserByEmail", mock.Anything, "test@example.com").
		Return(user, nil)

	resp, err := svc.SignInUser(&SignInRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	})

	assert.Error(t, err)
	assert.Empty(t, resp.AccessToken)
	assert.ErrorContains(t, err, "invalid credentials")
	repo.AssertExpectations(t)
}

func TestSignInUser_SaveTokenFails(t *testing.T) {
	repo := new(mockRepository)
	svc := newTestService(repo)

	user := sqlc.User{
		ID:           "user-1",
		Email:        "test@example.com",
		PasswordHash: hashedTestPassword(),
	}

	repo.
		On("GetUserByEmail", mock.Anything, "test@example.com").
		Return(user, nil)

	repo.
		On("CreateRefreshToken", mock.Anything, user.ID, mock.AnythingOfType("string"), "", "", mock.AnythingOfType("time.Time")).
		Return(sqlc.RefreshToken{}, errors.New("db error"))

	resp, err := svc.SignInUser(&SignInRequest{
		Email:    "test@example.com",
		Password: "password123",
	})

	assert.Error(t, err)
	assert.Empty(t, resp.AccessToken)
	assert.ErrorContains(t, err, "failed to generate tokens")
	repo.AssertExpectations(t)
}

// --- RefreshToken ---

func TestRefreshToken_Success(t *testing.T) {
	repo := new(mockRepository)
	svc, key := newTestServiceWithKey(repo)

	token := validRefreshToken(key, "user-1")
	stored := sqlc.RefreshToken{ID: "rt-1", UserID: "user-1"}

	repo.
		On("GetRefreshTokenByRawToken", mock.Anything, token).
		Return(stored, nil)

	repo.
		On("RevokeRefreshToken", mock.Anything, stored.ID).
		Return(nil)

	repo.
		On("CreateRefreshToken", mock.Anything, stored.UserID, mock.AnythingOfType("string"), "", "", mock.AnythingOfType("time.Time")).
		Return(sqlc.RefreshToken{}, nil)

	resp, err := svc.RefreshToken(&RefreshRequest{Token: token})

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	assert.NotEqual(t, token, resp.RefreshToken)
	repo.AssertExpectations(t)
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	repo := new(mockRepository)
	svc := newTestService(repo)

	resp, err := svc.RefreshToken(&RefreshRequest{Token: "not-a-valid-token"})

	assert.Error(t, err)
	assert.Empty(t, resp.AccessToken)
	assert.ErrorContains(t, err, "invalid refresh token")
	repo.AssertExpectations(t)
}

func TestRefreshToken_TokenNotFound(t *testing.T) {
	repo := new(mockRepository)
	svc, key := newTestServiceWithKey(repo)

	token := validRefreshToken(key, "user-1")

	repo.
		On("GetRefreshTokenByRawToken", mock.Anything, token).
		Return(sqlc.RefreshToken{}, errors.New("not found"))

	resp, err := svc.RefreshToken(&RefreshRequest{Token: token})

	assert.Error(t, err)
	assert.Empty(t, resp.AccessToken)
	assert.ErrorContains(t, err, "invalid or expired refresh token")
	repo.AssertExpectations(t)
}

func TestRefreshToken_RevokeFails(t *testing.T) {
	repo := new(mockRepository)
	svc, key := newTestServiceWithKey(repo)

	token := validRefreshToken(key, "user-1")
	stored := sqlc.RefreshToken{ID: "rt-1", UserID: "user-1"}

	repo.
		On("GetRefreshTokenByRawToken", mock.Anything, token).
		Return(stored, nil)

	repo.
		On("RevokeRefreshToken", mock.Anything, stored.ID).
		Return(errors.New("db error"))

	resp, err := svc.RefreshToken(&RefreshRequest{Token: token})

	assert.Error(t, err)
	assert.Empty(t, resp.AccessToken)
	assert.ErrorContains(t, err, "failed to revoke refresh token")
	repo.AssertExpectations(t)
}

func TestRefreshToken_SaveNewTokenFails(t *testing.T) {
	repo := new(mockRepository)
	svc, key := newTestServiceWithKey(repo)

	token := validRefreshToken(key, "user-1")
	stored := sqlc.RefreshToken{ID: "rt-1", UserID: "user-1"}

	repo.
		On("GetRefreshTokenByRawToken", mock.Anything, token).
		Return(stored, nil)

	repo.
		On("RevokeRefreshToken", mock.Anything, stored.ID).
		Return(nil)

	repo.
		On("CreateRefreshToken", mock.Anything, stored.UserID, mock.AnythingOfType("string"), "", "", mock.AnythingOfType("time.Time")).
		Return(sqlc.RefreshToken{}, errors.New("db error"))

	resp, err := svc.RefreshToken(&RefreshRequest{Token: token})

	assert.Error(t, err)
	assert.Empty(t, resp.AccessToken)
	assert.ErrorContains(t, err, "failed to generate tokens")
	repo.AssertExpectations(t)
}

// --- SignUpUser ---

func TestSignUpUser_Success(t *testing.T) {
	repo := new(mockRepository)
	svc := newTestService(repo)

	req := &SignUpRequest{Email: "new@example.com", Password: "password123"}
	user := sqlc.User{ID: "user-1", Email: "new@example.com"}

	repo.
		On("GetUserByEmail", mock.Anything, req.Email).
		Return(sqlc.User{}, errors.New("not found"))

	repo.
		On("CreateUser", mock.Anything, req.Email, mock.Anything).
		Return(user, nil)

	resp, err := svc.SignUpUser(req)

	assert.NoError(t, err)
	assert.Equal(t, user.Email, resp.Email)
	repo.AssertExpectations(t)
}

func TestSignUpUser_EmailAlreadyExists(t *testing.T) {
	repo := new(mockRepository)
	svc := newTestService(repo)

	req := &SignUpRequest{Email: "existing@example.com", Password: "password123"}

	repo.
		On("GetUserByEmail", mock.Anything, req.Email).
		Return(sqlc.User{ID: "user-1", Email: req.Email}, nil)

	resp, err := svc.SignUpUser(req)

	assert.Error(t, err)
	assert.Empty(t, resp.Email)
	assert.ErrorContains(t, err, "user already exists")
	repo.AssertExpectations(t)
}

func TestSignUpUser_CreateFails(t *testing.T) {
	repo := new(mockRepository)
	svc := newTestService(repo)

	req := &SignUpRequest{Email: "new@example.com", Password: "password123"}

	repo.
		On("GetUserByEmail", mock.Anything, req.Email).
		Return(sqlc.User{}, errors.New("not found"))

	repo.
		On("CreateUser", mock.Anything, req.Email, mock.Anything).
		Return(sqlc.User{}, errors.New("db error"))

	resp, err := svc.SignUpUser(req)

	assert.Error(t, err)
	assert.Empty(t, resp.Email)
	assert.ErrorContains(t, err, "failed to create user")
	repo.AssertExpectations(t)
}

// --- SignOut ---

func TestSignOut_Success(t *testing.T) {
	repo := new(mockRepository)
	svc, key := newTestServiceWithKey(repo)

	token := validRefreshToken(key, "user-1")
	stored := sqlc.RefreshToken{ID: "rt-1", UserID: "user-1"}

	repo.
		On("GetRefreshTokenByRawToken", mock.Anything, token).
		Return(stored, nil)

	repo.
		On("RevokeRefreshToken", mock.Anything, stored.ID).
		Return(nil)

	err := svc.SignOut(token)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestSignOut_TokenNotFound(t *testing.T) {
	repo := new(mockRepository)
	svc, key := newTestServiceWithKey(repo)

	token := validRefreshToken(key, "user-1")

	repo.
		On("GetRefreshTokenByRawToken", mock.Anything, token).
		Return(sqlc.RefreshToken{}, errors.New("not found"))

	err := svc.SignOut(token)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "invalid or expired refresh token")
	repo.AssertExpectations(t)
}

func TestSignOut_RevokeFails(t *testing.T) {
	repo := new(mockRepository)
	svc, key := newTestServiceWithKey(repo)

	token := validRefreshToken(key, "user-1")
	stored := sqlc.RefreshToken{ID: "rt-1", UserID: "user-1"}

	repo.
		On("GetRefreshTokenByRawToken", mock.Anything, token).
		Return(stored, nil)

	repo.
		On("RevokeRefreshToken", mock.Anything, stored.ID).
		Return(errors.New("db error"))

	err := svc.SignOut(token)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "failed to sign out")
	repo.AssertExpectations(t)
}
