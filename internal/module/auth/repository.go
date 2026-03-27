package auth

import (
	"context"

	"github.com/schodevio/trellgo/db/sqlc"
)

type Repository interface {
	CreateUser(ctx context.Context, user *CreateUserRequest) (sqlc.User, error)
}

type repository struct {
	queries *sqlc.Queries
}

func newRepository(queries *sqlc.Queries) Repository {
	return &repository{queries: queries}
}

func (r *repository) CreateUser(ctx context.Context, data *CreateUserRequest) (sqlc.User, error) {
	passwordHash, err := hashPassword(data.Password)
	if err != nil {
		return sqlc.User{}, err
	}

	return r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        data.Email,
		PasswordHash: passwordHash,
	})
}
