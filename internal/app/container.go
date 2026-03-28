package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/schodevio/trellgo/db/sqlc"
	"github.com/schodevio/trellgo/internal/module/auth"
	"github.com/schodevio/trellgo/internal/platform/config"
	"github.com/schodevio/trellgo/internal/platform/db"
)

type Container struct {
	DB   *pgxpool.Pool
	Auth *auth.Module
}

func NewContainer(ctx context.Context, cfg *config.Config) *Container {
	pool := db.NewPool(ctx, cfg.DBUrl)
	queries := sqlc.New(pool)

	authModule := auth.New(queries, cfg.SecretKey)

	return &Container{
		DB:   pool,
		Auth: authModule,
	}
}
