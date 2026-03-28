package auth

import (
	"context"

	paseto "aidanwoods.dev/go-paseto"
	"github.com/schodevio/trellgo/internal/platform/apierrors"
)

type Service interface {
	SignUpUser(data *SignUpRequest) (SignUpResponse, error)
	SignInUser(data *SignInRequest) (SignInResponse, error)
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
		return SignInResponse{}, apierrors.Unauthorized("invalid credentials")
	}

	if !checkPassword(data.Password, user.PasswordHash) {
		return SignInResponse{}, apierrors.Unauthorized("invalid credentials")
	}

	accessToken, refreshToken := generateTokens(user.ID, s.authKey)

	return SignInResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) SignUpUser(data *SignUpRequest) (SignUpResponse, error) {
	_, err := s.repo.GetUserByEmail(context.Background(), data.Email)
	if err == nil {
		return SignUpResponse{}, apierrors.Conflict("user already exists")
	}

	user, err := s.repo.CreateUser(context.Background(), data)
	if err != nil {
		return SignUpResponse{}, err
	}

	return SignUpResponse{Email: user.Email}, nil
}
