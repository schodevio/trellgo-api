package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/schodevio/trellgo/internal/platform/config"
	"github.com/schodevio/trellgo/internal/platform/db"
)

type Container struct {
	DB *pgxpool.Pool
}

func NewContainer(ctx context.Context, cfg *config.Config) *Container {
	pool := db.NewPool(ctx, cfg.DBUrl)

	return &Container{
		DB: pool,
	}
}
