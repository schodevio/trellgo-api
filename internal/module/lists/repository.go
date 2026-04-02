package lists

import (
	"context"

	"github.com/schodevio/trellgo/db/sqlc"
)

type Repository interface {
	CreateList(ctx context.Context, boardID, name string, position int32) (sqlc.List, error)
	GetBoardLists(ctx context.Context, boardID string) ([]sqlc.List, error)
	GetListByID(ctx context.Context, id string) (sqlc.List, error)
	UpdateListByID(ctx context.Context, id, name string, position int32) (sqlc.List, error)
}

type repository struct {
	queries *sqlc.Queries
}

func newRepository(queries *sqlc.Queries) Repository {
	return &repository{queries: queries}
}

func (r *repository) CreateList(ctx context.Context, boardID, name string, position int32) (sqlc.List, error) {
	return r.queries.CreateList(ctx, sqlc.CreateListParams{
		Name:     name,
		BoardID:  boardID,
		Position: position,
	})
}

func (r *repository) GetBoardLists(ctx context.Context, boardID string) ([]sqlc.List, error) {
	return r.queries.GetBoardLists(ctx, boardID)
}

func (r *repository) GetListByID(ctx context.Context, id string) (sqlc.List, error) {
	return r.queries.GetList(ctx, id)
}

func (r *repository) UpdateListByID(ctx context.Context, id, name string, position int32) (sqlc.List, error) {
	return r.queries.UpdateList(ctx, sqlc.UpdateListParams{
		ID:       id,
		Name:     name,
		Position: position,
	})
}
