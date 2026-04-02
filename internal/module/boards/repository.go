package boards

import (
	"context"

	"github.com/schodevio/trellgo/db/sqlc"
)

type Repository interface {
	CreateBoard(ctx context.Context, userID, name string) (sqlc.Board, error)
	GetUserBoards(ctx context.Context, userID string) ([]sqlc.Board, error)
	GetUserBoardByID(ctx context.Context, boardID, userID string) (sqlc.Board, error)
	GetBoardByID(ctx context.Context, boardID string) (sqlc.Board, error)
	UpdateBoardByID(ctx context.Context, boardID, name string) (sqlc.Board, error)
	DeleteBoardByID(ctx context.Context, boardID string) error
}

type repository struct {
	queries *sqlc.Queries
}

func newRepository(queries *sqlc.Queries) Repository {
	return &repository{queries: queries}
}

func (r *repository) CreateBoard(ctx context.Context, userID, name string) (sqlc.Board, error) {
	return r.queries.CreateBoard(ctx, sqlc.CreateBoardParams{
		Name:   name,
		UserID: userID,
	})
}

func (r *repository) GetUserBoards(ctx context.Context, userID string) ([]sqlc.Board, error) {
	return r.queries.GetUserBoards(ctx, userID)
}

func (r *repository) GetUserBoardByID(ctx context.Context, boardID, userID string) (sqlc.Board, error) {
	return r.queries.GetUserBoardByID(ctx, sqlc.GetUserBoardByIDParams{
		ID:     boardID,
		UserID: userID,
	})
}

func (r *repository) GetBoardByID(ctx context.Context, boardID string) (sqlc.Board, error) {
	return r.queries.GetBoardByID(ctx, boardID)
}

func (r *repository) UpdateBoardByID(ctx context.Context, boardID, name string) (sqlc.Board, error) {
	return r.queries.UpdateBoardByID(ctx, sqlc.UpdateBoardByIDParams{
		ID:   boardID,
		Name: name,
	})
}

func (r *repository) DeleteBoardByID(ctx context.Context, boardID string) error {
	return r.queries.DeleteBoardByID(ctx, boardID)
}
