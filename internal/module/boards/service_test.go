package boards

import (
	"context"
	"errors"
	"testing"

	"github.com/schodevio/trellgo/db/sqlc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- mock ---

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) CreateBoard(ctx context.Context, name, userID string) (sqlc.Board, error) {
	args := m.Called(ctx, name, userID)
	return args.Get(0).(sqlc.Board), args.Error(1)
}

func (m *mockRepository) GetUserBoards(ctx context.Context, userID string) ([]sqlc.Board, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]sqlc.Board), args.Error(1)
}

func (m *mockRepository) GetUserBoardByID(ctx context.Context, id, userID string) (sqlc.Board, error) {
	args := m.Called(ctx, id, userID)
	return args.Get(0).(sqlc.Board), args.Error(1)
}

func (m *mockRepository) UpdateUserBoardByID(ctx context.Context, name, id, userID string) (sqlc.Board, error) {
	args := m.Called(ctx, name, id, userID)
	return args.Get(0).(sqlc.Board), args.Error(1)
}

func (m *mockRepository) DeleteUserBoardByID(ctx context.Context, id, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

// --- tests ---

func TestCreateBoard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mockRepository)
		svc := newService(repo)

		req := &CreateBoardRequest{Name: "My Board", UserID: "user-1"}
		board := sqlc.Board{ID: "board-1", Name: "My Board", UserID: "user-1"}

		repo.
			On("CreateBoard", mock.Anything, req.Name, req.UserID).
			Return(board, nil)

		resp, err := svc.CreateBoard(req)

		assert.NoError(t, err)
		assert.Equal(t, board.ID, resp.Board.ID)
		assert.Equal(t, board.Name, resp.Board.Name)
		assert.Equal(t, board.UserID, resp.Board.UserID)
		repo.AssertExpectations(t)
	})

	t.Run("repo fails", func(t *testing.T) {
		repo := new(mockRepository)
		svc := newService(repo)

		req := &CreateBoardRequest{Name: "My Board", UserID: "user-1"}

		repo.
			On("CreateBoard", mock.Anything, req.Name, req.UserID).
			Return(sqlc.Board{}, errors.New("db error"))

		resp, err := svc.CreateBoard(req)

		assert.Error(t, err)
		assert.Empty(t, resp.Board.ID)
		assert.ErrorContains(t, err, "failed to create board")
		repo.AssertExpectations(t)
	})
}

// --- List ---

func TestListBoards(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mockRepository)
		svc := newService(repo)

		boards := []sqlc.Board{
			{ID: "board-1", Name: "First", UserID: "user-1"},
			{ID: "board-2", Name: "Second", UserID: "user-1"},
		}

		repo.
			On("GetUserBoards", mock.Anything, "user-1").
			Return(boards, nil)

		resp, err := svc.ListBoards("user-1")

		assert.NoError(t, err)
		assert.Len(t, resp.Boards, 2)
		assert.Equal(t, "board-1", resp.Boards[0].ID)
		assert.Equal(t, "board-2", resp.Boards[1].ID)
		repo.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		repo := new(mockRepository)
		svc := newService(repo)

		repo.
			On("GetUserBoards", mock.Anything, "user-1").
			Return([]sqlc.Board{}, nil)

		resp, err := svc.ListBoards("user-1")

		assert.NoError(t, err)
		assert.Empty(t, resp.Boards)
		repo.AssertExpectations(t)
	})

	t.Run("repo fails", func(t *testing.T) {
		repo := new(mockRepository)
		svc := newService(repo)

		repo.
			On("GetUserBoards", mock.Anything, "user-1").
			Return([]sqlc.Board{}, errors.New("db error"))

		resp, err := svc.ListBoards("user-1")

		assert.Error(t, err)
		assert.Empty(t, resp.Boards)
		assert.ErrorContains(t, err, "failed to fetch boards")
		repo.AssertExpectations(t)
	})
}

func TestGetBoard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mockRepository)
		svc := newService(repo)

		board := sqlc.Board{ID: "board-1", Name: "My Board", UserID: "user-1"}

		repo.
			On("GetUserBoardByID", mock.Anything, "board-1", "user-1").
			Return(board, nil)

		resp, err := svc.GetBoard("board-1", "user-1")

		assert.NoError(t, err)
		assert.Equal(t, board.ID, resp.Board.ID)
		assert.Equal(t, board.Name, resp.Board.Name)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(mockRepository)
		svc := newService(repo)

		repo.
			On("GetUserBoardByID", mock.Anything, "board-1", "user-2").
			Return(sqlc.Board{}, errors.New("not found"))

		resp, err := svc.GetBoard("board-1", "user-2")

		assert.Error(t, err)
		assert.Empty(t, resp.Board.ID)
		assert.ErrorContains(t, err, "board not found")
		repo.AssertExpectations(t)
	})
}

// --- Update ---

func TestUpdateBoard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mockRepository)
		svc := newService(repo)

		req := &UpdateBoardRequest{Name: "Renamed", ID: "board-1", UserID: "user-1"}
		board := sqlc.Board{ID: "board-1", Name: "Renamed", UserID: "user-1"}

		repo.
			On("UpdateUserBoardByID", mock.Anything, req.Name, req.ID, req.UserID).
			Return(board, nil)

		resp, err := svc.UpdateBoard(req)

		assert.NoError(t, err)
		assert.Equal(t, board.ID, resp.Board.ID)
		assert.Equal(t, board.Name, resp.Board.Name)
		assert.Equal(t, board.UserID, resp.Board.UserID)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(mockRepository)
		svc := newService(repo)

		req := &UpdateBoardRequest{Name: "Renamed", ID: "board-1", UserID: "user-2"}

		repo.
			On("UpdateUserBoardByID", mock.Anything, req.Name, req.ID, req.UserID).
			Return(sqlc.Board{}, errors.New("not found"))

		resp, err := svc.UpdateBoard(req)

		assert.Error(t, err)
		assert.Empty(t, resp.Board.ID)
		assert.ErrorContains(t, err, "board not found")
		repo.AssertExpectations(t)
	})
}

// --- Delete ---

func TestDeleteBoard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mockRepository)
		svc := newService(repo)

		repo.
			On("DeleteUserBoardByID", mock.Anything, "board-1", "user-1").
			Return(nil)

		err := svc.DeleteBoard("board-1", "user-1")

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(mockRepository)
		svc := newService(repo)

		repo.
			On("DeleteUserBoardByID", mock.Anything, "board-1", "user-2").
			Return(errors.New("not found"))

		err := svc.DeleteBoard("board-1", "user-2")

		assert.Error(t, err)
		assert.ErrorContains(t, err, "board not found")
		repo.AssertExpectations(t)
	})
}

// --- IsOwner ---

func TestIsOwner(t *testing.T) {
	t.Run("is owner", func(t *testing.T) {
		repo := new(mockRepository)
		svc := newService(repo)

		repo.
			On("GetUserBoardByID", mock.Anything, "board-1", "user-1").
			Return(sqlc.Board{ID: "board-1", UserID: "user-1"}, nil)

		err := svc.IsOwner("board-1", "user-1")

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("not owner", func(t *testing.T) {
		repo := new(mockRepository)
		svc := newService(repo)

		repo.
			On("GetUserBoardByID", mock.Anything, "board-1", "user-2").
			Return(sqlc.Board{}, errors.New("board not found"))

		err := svc.IsOwner("board-1", "user-2")

		assert.Error(t, err)
		assert.ErrorContains(t, err, "board not found")
		repo.AssertExpectations(t)
	})
}
