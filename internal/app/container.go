package app

import (
	"context"
	"log"

	"aidanwoods.dev/go-paseto"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/schodevio/trellgo/db/sqlc"
	"github.com/schodevio/trellgo/internal/module/auth"
	"github.com/schodevio/trellgo/internal/module/boards"
	"github.com/schodevio/trellgo/internal/module/health"
	"github.com/schodevio/trellgo/internal/module/lists"
	"github.com/schodevio/trellgo/internal/platform/config"
	"github.com/schodevio/trellgo/internal/platform/db"
)

type Container struct {
	DB     *pgxpool.Pool
	Health *health.Module
	Auth   *auth.Module
	Boards *boards.Module
	Lists  *lists.Module
}

func NewContainer(ctx context.Context, cfg *config.Config) *Container {
	authKey, err := paseto.V4SymmetricKeyFromHex(cfg.SecretKey)
	if err != nil {
		log.Fatalf("invalid SECRET_KEY: %v", err)
	}

	pool := db.NewPool(ctx, cfg.DBUrl)
	queries := sqlc.New(pool)

	healthModule := health.New(authKey)
	authModule := auth.New(queries, authKey)
	boardsModule := boards.New(queries, authKey)
	listsModule := lists.New(queries, authKey)

	return &Container{
		DB:     pool,
		Health: healthModule,
		Auth:   authModule,
		Boards: boardsModule,
		Lists:  listsModule,
	}
}
