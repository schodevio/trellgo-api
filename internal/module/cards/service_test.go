package cards

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/schodevio/trellgo/db/sqlc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- mock ---

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) CreateCard(ctx context.Context, listID, title, description, status string, position int32) (sqlc.Card, error) {
	args := m.Called(ctx, listID, title, description, status, position)
	return args.Get(0).(sqlc.Card), args.Error(1)
}

func (m *mockRepository) GetListCards(ctx context.Context, listID string) ([]sqlc.Card, error) {
	args := m.Called(ctx, listID)
	return args.Get(0).([]sqlc.Card), args.Error(1)
}

type mockListChecker struct {
	mock.Mock
}

func (m *mockListChecker) IsOwner(listID, userID string) (bool, error) {
	args := m.Called(listID, userID)
	return args.Bool(0), args.Error(1)
}

// --- tests ---

func TestCreateCard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mockRepository)
		listChecker := new(mockListChecker)
		svc := newService(repo, listChecker)

		req := &CreateCardRequest{
			Title:       "To Do",
			Description: "Task description",
			Status:      "open",
			Position:    1,
		}

		card := sqlc.Card{
			ID:          "card-1",
			Title:       req.Title,
			Description: pgtype.Text{String: req.Description, Valid: true},
			Status:      req.Status,
			Position:    req.Position,
		}

		listChecker.
			On("IsOwner", "list-1", "user-1").
			Return(true, nil)

		repo.
			On("CreateCard", mock.Anything, "list-1", req.Title, req.Description, req.Status, req.Position).
			Return(card, nil)

		resp, err := svc.CreateCard("list-1", "user-1", req)

		assert.NoError(t, err)
		assert.Equal(t, card.ID, resp.Card.ID)
		assert.Equal(t, card.Title, resp.Card.Title)
		assert.Equal(t, card.Description.String, resp.Card.Description)
		assert.Equal(t, card.Status, resp.Card.Status)
		assert.Equal(t, card.Position, resp.Card.Position)
		repo.AssertExpectations(t)
	})

	t.Run("not owner", func(t *testing.T) {
		repo := new(mockRepository)
		listChecker := new(mockListChecker)
		svc := newService(repo, listChecker)

		req := &CreateCardRequest{
			Title:       "To Do",
			Description: "Task description",
			Status:      "open",
			Position:    1,
		}

		listChecker.
			On("IsOwner", "list-1", "user-1").
			Return(false, assert.AnError)

		resp, err := svc.CreateCard("list-1", "user-1", req)

		assert.Error(t, err)
		assert.Equal(t, SingleCardResponse{}, resp)
		repo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(mockRepository)
		listChecker := new(mockListChecker)
		svc := newService(repo, listChecker)

		req := &CreateCardRequest{
			Title:       "To Do",
			Description: "Task description",
			Status:      "open",
			Position:    1,
		}

		listChecker.
			On("IsOwner", "list-1", "user-1").
			Return(true, nil)

		repo.
			On("CreateCard", mock.Anything, "list-1", req.Title, req.Description, req.Status, req.Position).
			Return(sqlc.Card{}, assert.AnError)

		resp, err := svc.CreateCard("list-1", "user-1", req)

		assert.Error(t, err)
		assert.Equal(t, SingleCardResponse{}, resp)
		repo.AssertExpectations(t)
	})
}

func TestListCards(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mockRepository)
		listChecker := new(mockListChecker)
		svc := newService(repo, listChecker)

		card := sqlc.Card{
			ID:          "card-1",
			ListID:      "list-1",
			Title:       "To Do",
			Description: pgtype.Text{String: "Task description", Valid: true},
			Status:      "open",
			Position:    1,
		}

		listChecker.
			On("IsOwner", "list-1", "user-1").
			Return(true, nil)

		repo.
			On("GetListCards", mock.Anything, "list-1").
			Return([]sqlc.Card{card}, nil)

		resp, err := svc.ListCards("list-1", "user-1")

		assert.NoError(t, err)
		assert.Len(t, resp.Cards, 1)
		assert.Equal(t, card.ID, resp.Cards[0].ID)
		assert.Equal(t, card.ListID, resp.Cards[0].ListID)
		assert.Equal(t, card.Title, resp.Cards[0].Title)
		assert.Equal(t, card.Description.String, resp.Cards[0].Description)
		assert.Equal(t, card.Status, resp.Cards[0].Status)
		assert.Equal(t, card.Position, resp.Cards[0].Position)
		repo.AssertExpectations(t)
	})

	t.Run("not owner", func(t *testing.T) {
		repo := new(mockRepository)
		listChecker := new(mockListChecker)
		svc := newService(repo, listChecker)

		listChecker.
			On("IsOwner", "list-1", "user-1").
			Return(false, assert.AnError)

		resp, err := svc.ListCards("list-1", "user-1")

		assert.Error(t, err)
		assert.Equal(t, ListCardsResponse{}, resp)
		repo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(mockRepository)
		listChecker := new(mockListChecker)
		svc := newService(repo, listChecker)

		listChecker.
			On("IsOwner", "list-1", "user-1").
			Return(true, nil)

		repo.
			On("GetListCards", mock.Anything, "list-1").
			Return([]sqlc.Card{}, assert.AnError)

		resp, err := svc.ListCards("list-1", "user-1")

		assert.Error(t, err)
		assert.Equal(t, ListCardsResponse{}, resp)
		repo.AssertExpectations(t)
	})
}
