package auth

import (
	"context"
	"time"

	paseto "aidanwoods.dev/go-paseto"
	"github.com/schodevio/trellgo/internal/platform/apierrors"
)

type Service interface {
	RefreshToken(data *RefreshRequest) (SignInResponse, error)
	SignUpUser(data *SignUpRequest) (SignUpResponse, error)
	SignInUser(data *SignInRequest) (SignInResponse, error)
	SignOut(token string) error
}

type service struct {
	repo    Repository
	authKey paseto.V4SymmetricKey
}

func newService(repo Repository, authKey paseto.V4SymmetricKey) Service {
	return &service{repo: repo, authKey: authKey}
}

func (s *service) SignInUser(data *SignInRequest) (SignInResponse, error) {
	user, err := s.repo.GetUserByEmail(context.Background(), data.Email)
	if err != nil {
		return SignInResponse{}, apierrors.Unauthorized(INVALID_CREDENTIALS)
	}

	if !checkPassword(data.Password, user.PasswordHash) {
		return SignInResponse{}, apierrors.Unauthorized(INVALID_CREDENTIALS)
	}

	accessToken, refreshToken, err := s.generateTokens(&generateTokensData{
		UserID:    user.ID,
		UserAgent: data.UserAgent,
		IpAddress: data.IpAddress,
	})
	if err != nil {
		return SignInResponse{}, err
	}

	return SignInResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) RefreshToken(data *RefreshRequest) (SignInResponse, error) {
	_, err := parseRefreshToken(data.Token, s.authKey)
	if err != nil {
		return SignInResponse{}, apierrors.Unauthorized(INVALID_OR_EXPIRED_TOKEN)
	}

	stored, err := s.repo.GetRefreshTokenByRawToken(context.Background(), data.Token)
	if err != nil {
		return SignInResponse{}, apierrors.Unauthorized(INVALID_OR_EXPIRED_TOKEN)
	}

	if err := s.repo.RevokeRefreshToken(context.Background(), stored.ID); err != nil {
		return SignInResponse{}, apierrors.Internal(REFRESH_TOKEN_REVOKE_FAILED)
	}

	accessToken, newRefreshToken, err := s.generateTokens(&generateTokensData{
		UserID:    stored.UserID,
		UserAgent: data.UserAgent,
		IpAddress: data.IpAddress,
	})
	if err != nil {
		return SignInResponse{}, err
	}

	return SignInResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *service) SignUpUser(data *SignUpRequest) (SignUpResponse, error) {
	_, err := s.repo.GetUserByEmail(context.Background(), data.Email)
	if err == nil {
		return SignUpResponse{}, apierrors.Conflict(USER_ALREADY_EXISTS)
	}

	passwordHash, err := hashPassword(data.Password)
	if err != nil {
		return SignUpResponse{}, apierrors.Internal(PASSWORD_HASH_FAILED)
	}

	user, err := s.repo.CreateUser(context.Background(), data.Email, passwordHash)
	if err != nil {
		return SignUpResponse{}, apierrors.Internal(USER_CREATE_FAILED)
	}

	return SignUpResponse{Email: user.Email}, nil
}

func (s *service) SignOut(token string) error {
	stored, err := s.repo.GetRefreshTokenByRawToken(context.Background(), token)
	if err != nil {
		return apierrors.Unauthorized(INVALID_OR_EXPIRED_TOKEN)
	}

	if err := s.repo.RevokeRefreshToken(context.Background(), stored.ID); err != nil {
		return apierrors.Internal(SIGN_OUT_FAILED)
	}

	return nil
}

// private

func (s *service) generateTokens(data *generateTokensData) (string, string, error) {
	currentTime := time.Now()
	accessExpiration := currentTime.Add(accessTokenTTL)
	refreshExpiration := currentTime.Add(refreshTokenTTL)

	access := paseto.NewToken()
	access.SetSubject(data.UserID)
	access.SetIssuedAt(currentTime)
	access.SetNotBefore(currentTime)
	access.SetExpiration(accessExpiration)
	access.SetString("typ", "access")

	refresh := paseto.NewToken()
	refresh.SetSubject(data.UserID)
	refresh.SetIssuedAt(currentTime)
	refresh.SetNotBefore(currentTime)
	refresh.SetExpiration(refreshExpiration)
	refresh.SetString("typ", "refresh")

	accessToken := access.V4Encrypt(s.authKey, nil)
	refreshToken := refresh.V4Encrypt(s.authKey, nil)

	_, err := s.repo.CreateRefreshToken(
		context.Background(),
		data.UserID,
		refreshToken,
		data.UserAgent,
		data.IpAddress,
		refreshExpiration,
	)
	if err != nil {
		return "", "", apierrors.Internal(TOKEN_GENERATE_FAILED)
	}

	return accessToken, refreshToken, nil
}
