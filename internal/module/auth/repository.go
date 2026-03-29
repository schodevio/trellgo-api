package auth

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/schodevio/trellgo/db/sqlc"
)

type Repository interface {
	// Refresh tokens
	CreateRefreshToken(ctx context.Context, userID, token, userAgent, ipAddress string, expiresAt time.Time) (sqlc.RefreshToken, error)
	GetRefreshTokenByRawToken(ctx context.Context, token string) (sqlc.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, id string) error

	// Users
	CreateUser(ctx context.Context, data *SignUpRequest) (sqlc.User, error)
	GetUserByEmail(ctx context.Context, email string) (sqlc.User, error)
}

type repository struct {
	queries *sqlc.Queries
}

func newRepository(queries *sqlc.Queries) Repository {
	return &repository{queries: queries}
}

func (r *repository) CreateRefreshToken(ctx context.Context, userID, token, userAgent, ipAddress string, expiresAt time.Time) (sqlc.RefreshToken, error) {
	return r.queries.CreateRefreshToken(ctx, sqlc.CreateRefreshTokenParams{
		UserID:    userID,
		TokenHash: hashToken(token),
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
		UserAgent: pgtype.Text{String: userAgent, Valid: userAgent != ""},
		IpAddress: pgtype.Text{String: ipAddress, Valid: ipAddress != ""},
	})
}

func (r *repository) GetRefreshTokenByRawToken(ctx context.Context, token string) (sqlc.RefreshToken, error) {
	return r.queries.GetRefreshTokenByHash(ctx, hashToken(token))
}

func (r *repository) RevokeRefreshToken(ctx context.Context, id string) error {
	_, err := r.queries.RevokeRefreshToken(ctx, id)
	return err
}

func (r *repository) CreateUser(ctx context.Context, data *SignUpRequest) (sqlc.User, error) {
	passwordHash, err := hashPassword(data.Password)
	if err != nil {
		return sqlc.User{}, err
	}

	return r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        data.Email,
		PasswordHash: passwordHash,
	})
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (sqlc.User, error) {
	return r.queries.GetUserByEmail(ctx, email)
}
