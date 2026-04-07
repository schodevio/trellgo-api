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

func (m *mockRepository) CreateCard(ctx context.Context, listID, title, description string, position int32) (sqlc.Card, error) {
	args := m.Called(ctx, listID, title, description, position)
	return args.Get(0).(sqlc.Card), args.Error(1)
}

func (m *mockRepository) GetListCards(ctx context.Context, listID string) ([]sqlc.Card, error) {
	args := m.Called(ctx, listID)
	return args.Get(0).([]sqlc.Card), args.Error(1)
}

func (m *mockRepository) GetCardByID(ctx context.Context, cardID string) (sqlc.Card, error) {
	args := m.Called(ctx, cardID)
	return args.Get(0).(sqlc.Card), args.Error(1)
}

func (m *mockRepository) UpdateCardByID(ctx context.Context, cardID, title, description string) (sqlc.Card, error) {
	args := m.Called(ctx, cardID, title, description)
	return args.Get(0).(sqlc.Card), args.Error(1)
}

func (m *mockRepository) MoveCardByID(ctx context.Context, cardID, listID string, position int32) (sqlc.Card, error) {
	args := m.Called(ctx, cardID, listID, position)
	return args.Get(0).(sqlc.Card), args.Error(1)
}

func (m *mockRepository) ReorderCardsInList(ctx context.Context, listID string) error {
	args := m.Called(ctx, listID)
	return args.Error(0)
}

func (m *mockRepository) DeleteCardByID(ctx context.Context, cardID string) error {
	args := m.Called(ctx, cardID)
	return args.Error(0)
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
			Position:    1,
		}

		card := sqlc.Card{
			ID:          "card-1",
			Title:       req.Title,
			Description: pgtype.Text{String: req.Description, Valid: true},
			Position:    req.Position,
		}

		listChecker.
			On("IsOwner", "list-1", "user-1").
			Return(true, nil)

		repo.
			On("CreateCard", mock.Anything, "list-1", req.Title, req.Description, req.Position).
			Return(card, nil)

		repo.
			On("ReorderCardsInList", mock.Anything, "list-1").
			Return(nil)

		resp, err := svc.CreateCard("list-1", "user-1", req)

		assert.NoError(t, err)
		assert.Equal(t, card.ID, resp.Card.ID)
		assert.Equal(t, card.Title, resp.Card.Title)
		assert.Equal(t, card.Description.String, resp.Card.Description)
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

	t.Run("repo error on create", func(t *testing.T) {
		repo := new(mockRepository)
		listChecker := new(mockListChecker)
		svc := newService(repo, listChecker)

		req := &CreateCardRequest{
			Title:       "To Do",
			Description: "Task description",
			Position:    1,
		}

		listChecker.
			On("IsOwner", "list-1", "user-1").
			Return(true, nil)

		repo.
			On("CreateCard", mock.Anything, "list-1", req.Title, req.Description, req.Position).
			Return(sqlc.Card{}, assert.AnError)

		resp, err := svc.CreateCard("list-1", "user-1", req)

		assert.Error(t, err)
		assert.Equal(t, SingleCardResponse{}, resp)
		repo.AssertExpectations(t)
	})

	t.Run("repo error on reorder", func(t *testing.T) {
		repo := new(mockRepository)
		listChecker := new(mockListChecker)
		svc := newService(repo, listChecker)

		req := &CreateCardRequest{
			Title:       "To Do",
			Description: "Task description",
			Position:    1,
		}

		card := sqlc.Card{
			ID:          "card-1",
			ListID:      "list-1",
			Title:       req.Title,
			Description: pgtype.Text{String: req.Description, Valid: true},
			Position:    req.Position,
		}

		listChecker.
			On("IsOwner", "list-1", "user-1").
			Return(true, nil)

		repo.
			On("CreateCard", mock.Anything, "list-1", req.Title, req.Description, req.Position).
			Return(card, nil)

		repo.
			On("ReorderCardsInList", mock.Anything, "list-1").
			Return(assert.AnError)

		resp, err := svc.CreateCard("list-1", "user-1", req)

		assert.Error(t, err)
		assert.Equal(t, SingleCardResponse{}, resp)
		assert.ErrorContains(t, err, "failed to create card")
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

func TestUpdateCard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mockRepository)
		listChecker := new(mockListChecker)
		svc := newService(repo, listChecker)

		card := sqlc.Card{
			ID:          "card-1",
			ListID:      "list-1",
			Title:       "To Do",
			Description: pgtype.Text{String: "Task description", Valid: true},
			Position:    1,
		}

		req := &UpdateCardRequest{
			Title:       "In Progress",
			Description: "Updated task description",
		}

		listChecker.
			On("IsOwner", "list-1", "user-1").
			Return(true, nil)

		repo.
			On("GetCardByID", mock.Anything, "card-1").
			Return(card, nil)

		repo.
			On("UpdateCardByID", mock.Anything, "card-1", req.Title, req.Description).
			Return(sqlc.Card{
				ID:          "card-1",
				ListID:      "list-1",
				Title:       req.Title,
				Description: pgtype.Text{String: req.Description, Valid: true},
				Position:    card.Position,
			}, nil)

		resp, err := svc.UpdateCard("card-1", "user-1", req)

		assert.NoError(t, err)
		assert.Equal(t, req.Title, resp.Card.Title)
		assert.Equal(t, req.Description, resp.Card.Description)
		assert.Equal(t, card.Position, resp.Card.Position)
		repo.AssertExpectations(t)
	})

	t.Run("card not found", func(t *testing.T) {
		repo := new(mockRepository)
		listChecker := new(mockListChecker)
		svc := newService(repo, listChecker)

		req := &UpdateCardRequest{
			Title:       "In Progress",
			Description: "Updated task description",
		}

		repo.
			On("GetCardByID", mock.Anything, "card-1").
			Return(sqlc.Card{}, assert.AnError)

		resp, err := svc.UpdateCard("card-1", "user-1", req)

		assert.Error(t, err)
		assert.Equal(t, SingleCardResponse{}, resp)
		repo.AssertExpectations(t)
	})

	t.Run("not owner", func(t *testing.T) {
		repo := new(mockRepository)
		listChecker := new(mockListChecker)
		svc := newService(repo, listChecker)

		card := sqlc.Card{
			ID:          "card-1",
			ListID:      "list-1",
			Title:       "To Do",
			Description: pgtype.Text{String: "Task description", Valid: true},
			Position:    1,
		}

		req := &UpdateCardRequest{
			Title:       "In Progress",
			Description: "Updated task description",
		}

		repo.
			On("GetCardByID", mock.Anything, "card-1").
			Return(card, nil)

		listChecker.
			On("IsOwner", "list-1", "user-1").
			Return(false, assert.AnError)

		resp, err := svc.UpdateCard("card-1", "user-1", req)

		assert.Error(t, err)
		assert.Equal(t, SingleCardResponse{}, resp)
		repo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(mockRepository)
		listChecker := new(mockListChecker)
		svc := newService(repo, listChecker)

		card := sqlc.Card{
			ID:          "card-1",
			ListID:      "list-1",
			Title:       "To Do",
			Description: pgtype.Text{String: "Task description", Valid: true},
			Position:    1,
		}

		req := &UpdateCardRequest{
			Title:       "In Progress",
			Description: "Updated task description",
		}

		repo.
			On("GetCardByID", mock.Anything, "card-1").
			Return(card, nil)

		listChecker.
			On("IsOwner", "list-1", "user-1").
			Return(true, nil)

		repo.
			On("UpdateCardByID", mock.Anything, "card-1", req.Title, req.Description).
			Return(sqlc.Card{}, assert.AnError)

		resp, err := svc.UpdateCard("card-1", "user-1", req)

		assert.Error(t, err)
		assert.Equal(t, SingleCardResponse{}, resp)
		repo.AssertExpectations(t)
	})
}

func TestMoveCard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mockRepository)
		listChecker := new(mockListChecker)
		svc := newService(repo, listChecker)

		card := sqlc.Card{
			ID:          "card-1",
			ListID:      "list-1",
			Title:       "To Do",
			Description: pgtype.Text{String: "Task description", Valid: true},
			Position:    1,
		}

		req := &MoveCardRequest{
			ListID:   "list-2",
			Position: 0,
		}

		listChecker.
			On("IsOwner", "list-1", "user-1").
			Return(true, nil)

		listChecker.
			On("IsOwner", "list-2", "user-1").
			Return(true, nil)

		repo.
			On("GetCardByID", mock.Anything, "card-1").
			Return(card, nil)

		repo.
			On("MoveCardByID", mock.Anything, "card-1", req.ListID, req.Position).
			Return(sqlc.Card{
				ID:          "card-1",
				ListID:      req.ListID,
				Title:       card.Title,
				Description: card.Description,
				Position:    req.Position,
			}, nil)

		repo.
			On("ReorderCardsInList", mock.Anything, "list-2").
			Return(nil)

		repo.
			On("ReorderCardsInList", mock.Anything, "list-1").
			Return(nil)

		resp, err := svc.MoveCard("card-1", "user-1", req)

		assert.NoError(t, err)
		assert.Equal(t, req.ListID, resp.Card.ListID)
		assert.Equal(t, req.Position, resp.Card.Position)
		repo.AssertExpectations(t)
	})

	t.Run("card not found", func(t *testing.T) {
		repo := new(mockRepository)
		listChecker := new(mockListChecker)
		svc := newService(repo, listChecker)

		req := &MoveCardRequest{
			ListID:   "list-2",
			Position: 0,
		}

		repo.
			On("GetCardByID", mock.Anything, "card-1").
			Return(sqlc.Card{}, assert.AnError)

		resp, err := svc.MoveCard("card-1", "user-1", req)

		assert.Error(t, err)
		assert.Equal(t, SingleCardResponse{}, resp)
		repo.AssertExpectations(t)
	})

	t.Run("not owner for original list", func(t *testing.T) {
		repo := new(mockRepository)
		listChecker := new(mockListChecker)
		svc := newService(repo, listChecker)

		card := sqlc.Card{
			ID:          "card-1",
			ListID:      "list-1",
			Title:       "To Do",
			Description: pgtype.Text{String: "Task description", Valid: true},
			Position:    1,
		}

		req := &MoveCardRequest{
			ListID:   "list-2",
			Position: 0,
		}

		repo.
			On("GetCardByID", mock.Anything, "card-1").
			Return(card, nil)

		listChecker.
			On("IsOwner", "list-1", "user-1").
			Return(false, assert.AnError)

		resp, err := svc.MoveCard("card-1", "user-1", req)

		assert.Error(t, err)
		assert.Equal(t, SingleCardResponse{}, resp)
		repo.AssertExpectations(t)
	})

	t.Run("not owner for new list", func(t *testing.T) {
		repo := new(mockRepository)
		listChecker := new(mockListChecker)
		svc := newService(repo, listChecker)

		card := sqlc.Card{
			ID:          "card-1",
			ListID:      "list-1",
			Title:       "To Do",
			Description: pgtype.Text{String: "Task description", Valid: true},
			Position:    1,
		}

		req := &MoveCardRequest{
			ListID:   "list-2",
			Position: 0,
		}

		repo.
			On("GetCardByID", mock.Anything, "card-1").
			Return(card, nil)

		listChecker.
			On("IsOwner", "list-1", "user-1").
			Return(true, nil)

		listChecker.
			On("IsOwner", "list-2", "user-1").
			Return(false, assert.AnError)

		resp, err := svc.MoveCard("card-1", "user-1", req)

		assert.Error(t, err)
		assert.Equal(t, SingleCardResponse{}, resp)
		repo.AssertExpectations(t)
	})

	t.Run("repo error on move", func(t *testing.T) {
		repo := new(mockRepository)
		listChecker := new(mockListChecker)
		svc := newService(repo, listChecker)

		card := sqlc.Card{
			ID:          "card-1",
			ListID:      "list-1",
			Title:       "To Do",
			Description: pgtype.Text{String: "Task description", Valid: true},
			Position:    1,
		}

		req := &MoveCardRequest{
			ListID:   "list-2",
			Position: 0,
		}

		repo.
			On("GetCardByID", mock.Anything, "card-1").
			Return(card, nil)

		listChecker.
			On("IsOwner", "list-1", "user-1").
			Return(true, nil)

		listChecker.
			On("IsOwner", "list-2", "user-1").
			Return(true, nil)

		repo.
			On("MoveCardByID", mock.Anything, "card-1", req.ListID, req.Position).
			Return(sqlc.Card{}, assert.AnError)

		resp, err := svc.MoveCard("card-1", "user-1", req)

		assert.Error(t, err)
		assert.Equal(t, SingleCardResponse{}, resp)
		assert.ErrorContains(t, err, "failed to move card")
		repo.AssertExpectations(t)
	})

	t.Run("repo error on reorder", func(t *testing.T) {
		repo := new(mockRepository)
		listChecker := new(mockListChecker)
		svc := newService(repo, listChecker)

		card := sqlc.Card{
			ID:          "card-1",
			ListID:      "list-1",
			Title:       "To Do",
			Description: pgtype.Text{String: "Task description", Valid: true},
			Position:    1,
		}

		req := &MoveCardRequest{
			ListID:   "list-2",
			Position: 0,
		}

		repo.
			On("GetCardByID", mock.Anything, "card-1").
			Return(card, nil)

		listChecker.
			On("IsOwner", "list-1", "user-1").
			Return(true, nil)

		listChecker.
			On("IsOwner", "list-2", "user-1").
			Return(true, nil)

		repo.
			On("MoveCardByID", mock.Anything, "card-1", req.ListID, req.Position).
			Return(sqlc.Card{
				ID:          "card-1",
				ListID:      req.ListID,
				Title:       card.Title,
				Description: card.Description,
				Position:    req.Position,
			}, nil)

		repo.
			On("ReorderCardsInList", mock.Anything, "list-2").
			Return(assert.AnError)

		resp, err := svc.MoveCard("card-1", "user-1", req)

		assert.Error(t, err)
		assert.Equal(t, SingleCardResponse{}, resp)
		assert.ErrorContains(t, err, "failed to move card")
		repo.AssertExpectations(t)
	})
}

func TestDeleteCard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mockRepository)
		listChecker := new(mockListChecker)
		svc := newService(repo, listChecker)

		card := sqlc.Card{
			ID:          "card-1",
			ListID:      "list-1",
			Title:       "To Do",
			Description: pgtype.Text{String: "Task description", Valid: true},
			Position:    1,
		}

		repo.
			On("GetCardByID", mock.Anything, "card-1").
			Return(card, nil)

		listChecker.
			On("IsOwner", "list-1", "user-1").
			Return(true, nil)

		repo.
			On("DeleteCardByID", mock.Anything, "card-1").
			Return(nil)

		err := svc.DeleteCard("card-1", "user-1")

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("card not found", func(t *testing.T) {
		repo := new(mockRepository)
		listChecker := new(mockListChecker)
		svc := newService(repo, listChecker)

		repo.
			On("GetCardByID", mock.Anything, "card-1").
			Return(sqlc.Card{}, assert.AnError)

		err := svc.DeleteCard("card-1", "user-1")

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("not owner", func(t *testing.T) {
		repo := new(mockRepository)
		listChecker := new(mockListChecker)
		svc := newService(repo, listChecker)

		card := sqlc.Card{
			ID:          "card-1",
			ListID:      "list-1",
			Title:       "To Do",
			Description: pgtype.Text{String: "Task description", Valid: true},
			Position:    1,
		}

		repo.
			On("GetCardByID", mock.Anything, "card-1").
			Return(card, nil)

		listChecker.
			On("IsOwner", "list-1", "user-1").
			Return(false, assert.AnError)

		err := svc.DeleteCard("card-1", "user-1")

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(mockRepository)
		listChecker := new(mockListChecker)
		svc := newService(repo, listChecker)

		card := sqlc.Card{
			ID:          "card-1",
			ListID:      "list-1",
			Title:       "To Do",
			Description: pgtype.Text{String: "Task description", Valid: true},
			Position:    1,
		}

		repo.
			On("GetCardByID", mock.Anything, "card-1").
			Return(card, nil)

		listChecker.
			On("IsOwner", "list-1", "user-1").
			Return(true, nil)

		repo.
			On("DeleteCardByID", mock.Anything, "card-1").
			Return(assert.AnError)

		err := svc.DeleteCard("card-1", "user-1")

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}
