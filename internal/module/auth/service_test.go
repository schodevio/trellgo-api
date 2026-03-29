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

func (m *mockRepository) CreateUser(ctx context.Context, data *SignUpRequest) (sqlc.User, error) {
	args := m.Called(ctx, data)
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
