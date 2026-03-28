package auth

import (
	"context"

	"github.com/schodevio/trellgo/internal/platform/apierrors"
)

type Service interface {
	CreateUser(data *SignUpRequest) (SignUpResponse, error)
}

type service struct {
	repo Repository
}

func newService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateUser(data *SignUpRequest) (SignUpResponse, error) {
	user, err := s.repo.GetUserByEmail(context.Background(), data.Email)
	if err == nil {
		return SignUpResponse{}, apierrors.Conflict("user already exists")
	}

	user, err = s.repo.CreateUser(context.Background(), data)
	if err != nil {
		return SignUpResponse{}, err
	}

	return SignUpResponse{Email: user.Email}, nil
}
