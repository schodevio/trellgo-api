package boards

import (
	"context"

	"github.com/schodevio/trellgo/db/sqlc"
)

type Repository interface {
	CreateBoard(ctx context.Context, name, userID string) (sqlc.Board, error)
	GetUserBoards(ctx context.Context, userID string) ([]sqlc.Board, error)
	GetUserBoardByID(ctx context.Context, id, userID string) (sqlc.Board, error)
	UpdateUserBoardByID(ctx context.Context, name, id, userID string) (sqlc.Board, error)
	DeleteUserBoardByID(ctx context.Context, id, userID string) error
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

func (r *repository) GetUserBoards(ctx context.Context, userID string) ([]sqlc.Board, error) {
	return r.queries.GetUserBoards(ctx, userID)
}

func (r *repository) GetUserBoardByID(ctx context.Context, id, userID string) (sqlc.Board, error) {
	return r.queries.GetUserBoardByID(ctx, sqlc.GetUserBoardByIDParams{
		ID:     id,
		UserID: userID,
	})
}

func (r *repository) UpdateUserBoardByID(ctx context.Context, name, id, userID string) (sqlc.Board, error) {
	return r.queries.UpdateUserBoardByID(ctx, sqlc.UpdateUserBoardByIDParams{
		Name:   name,
		ID:     id,
		UserID: userID,
	})
}

func (r *repository) DeleteUserBoardByID(ctx context.Context, id, userID string) error {
	return r.queries.DeleteUserBoardByID(ctx, sqlc.DeleteUserBoardByIDParams{
		ID:     id,
		UserID: userID,
	})
}
