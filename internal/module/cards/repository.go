package cards

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/schodevio/trellgo/db/sqlc"
)

type Repository interface {
	CreateCard(ctx context.Context, listID, title, description, status string, position int32) (sqlc.Card, error)
	GetListCards(ctx context.Context, listID string) ([]sqlc.Card, error)
	GetCardByID(ctx context.Context, cardID string) (sqlc.Card, error)
	UpdateCardByID(ctx context.Context, cardID, title, description, status string, position int32) (sqlc.Card, error)
}

type repository struct {
	queries *sqlc.Queries
}

func newRepository(queries *sqlc.Queries) Repository {
	return &repository{queries: queries}
}

func (r *repository) CreateCard(ctx context.Context, listID, title, description, status string, position int32) (sqlc.Card, error) {
	return r.queries.CreateCard(ctx, sqlc.CreateCardParams{
		ListID:      listID,
		Title:       title,
		Description: pgtype.Text{String: description, Valid: description != ""},
		Status:      status,
		Position:    position,
	})
}

func (r *repository) GetListCards(ctx context.Context, listID string) ([]sqlc.Card, error) {
	return r.queries.GetListCards(ctx, listID)
}

func (r *repository) GetCardByID(ctx context.Context, cardID string) (sqlc.Card, error) {
	return r.queries.GetCardByID(ctx, cardID)
}

func (r *repository) UpdateCardByID(ctx context.Context, cardID, title, description, status string, position int32) (sqlc.Card, error) {
	return r.queries.UpdateCardByID(ctx, sqlc.UpdateCardByIDParams{
		ID:          cardID,
		Title:       title,
		Description: pgtype.Text{String: description, Valid: description != ""},
		Status:      status,
		Position:    position,
	})
}
