package lists

import (
	"context"

	"github.com/schodevio/trellgo/db/sqlc"
)

type Repository interface {
	CreateList(ctx context.Context, name string, boardID string, position int32) (sqlc.List, error)
}

type repository struct {
	queries *sqlc.Queries
}

func newRepository(queries *sqlc.Queries) Repository {
	return &repository{queries: queries}
}

func (r *repository) CreateList(ctx context.Context, name string, boardID string, position int32) (sqlc.List, error) {
	return r.queries.CreateList(ctx, sqlc.CreateListParams{
		Name:     name,
		BoardID:  boardID,
		Position: position,
	})
}
