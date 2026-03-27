package auth

import (
	"context"
)

type Service interface {
	CreateUser(user *CreateUserRequest) (UserResponse, error)
}

type service struct {
	repo Repository
}

func newService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateUser(data *CreateUserRequest) (UserResponse, error) {
	user, err := s.repo.CreateUser(context.Background(), data)
	if err != nil {
		return UserResponse{}, err
	}

	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}, nil
}
