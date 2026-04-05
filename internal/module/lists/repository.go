package lists

import (
	"context"

	"github.com/schodevio/trellgo/db/sqlc"
)

type Repository interface {
	CreateList(ctx context.Context, boardID, name string, position int32) (sqlc.List, error)
	GetBoardLists(ctx context.Context, boardID string) ([]sqlc.List, error)
	GetUserListByID(ctx context.Context, listID, userID string) (sqlc.List, error)
	GetListByID(ctx context.Context, listID string) (sqlc.List, error)
	UpdateListByID(ctx context.Context, listID, name string, position int32) (sqlc.List, error)
	DeleteListByID(ctx context.Context, listID string) error
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

func (r *repository) GetListByID(ctx context.Context, listID string) (sqlc.List, error) {
	return r.queries.GetList(ctx, listID)
}

func (r *repository) GetUserListByID(ctx context.Context, listID, userID string) (sqlc.List, error) {
	return r.queries.GetUserListByID(ctx, sqlc.GetUserListByIDParams{
		ID:     listID,
		UserID: userID,
	})
}

func (r *repository) UpdateListByID(ctx context.Context, listID, name string, position int32) (sqlc.List, error) {
	return r.queries.UpdateList(ctx, sqlc.UpdateListParams{
		ID:       listID,
		Name:     name,
		Position: position,
	})
}

func (r *repository) DeleteListByID(ctx context.Context, listID string) error {
	return r.queries.DeleteList(ctx, listID)
}
