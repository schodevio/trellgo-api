package boards

import (
	"context"

	"github.com/schodevio/trellgo/db/sqlc"
)

type Repository interface {
	CreateBoard(ctx context.Context, name, userID string) (sqlc.Board, error)
}

type repository struct {
	queries *sqlc.Queries
}

func newRepository(queries *sqlc.Queries) Repository {
	return &repository{queries: queries}
}

func (r *repository) CreateBoard(ctx context.Context, name, userID string) (sqlc.Board, error) {
	return r.queries.CreateBoard(ctx, sqlc.CreateBoardParams{
		Name:   name,
		UserID: userID,
	})
}
