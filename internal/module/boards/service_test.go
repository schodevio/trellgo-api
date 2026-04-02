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

func (m *mockRepository) CreateBoard(ctx context.Context, userID, name string) (sqlc.Board, error) {
	args := m.Called(ctx, userID, name)
	return args.Get(0).(sqlc.Board), args.Error(1)
}

func (m *mockRepository) GetUserBoards(ctx context.Context, userID string) ([]sqlc.Board, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]sqlc.Board), args.Error(1)
}

func (m *mockRepository) GetUserBoardByID(ctx context.Context, boardID, userID string) (sqlc.Board, error) {
	args := m.Called(ctx, boardID, userID)
	return args.Get(0).(sqlc.Board), args.Error(1)
}

func (m *mockRepository) GetBoardByID(ctx context.Context, boardID string) (sqlc.Board, error) {
	args := m.Called(ctx, boardID)
	return args.Get(0).(sqlc.Board), args.Error(1)
}

func (m *mockRepository) UpdateBoardByID(ctx context.Context, boardID, name string) (sqlc.Board, error) {
	args := m.Called(ctx, boardID, name)
	return args.Get(0).(sqlc.Board), args.Error(1)
}

func (m *mockRepository) DeleteBoardByID(ctx context.Context, boardID string) error {
	args := m.Called(ctx, boardID)
	return args.Error(0)
}

// --- tests ---

func TestCreateBoard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mockRepository)
		svc := newService(repo)

		req := &CreateBoardRequest{Name: "My Board"}
		board := sqlc.Board{ID: "board-1", Name: "My Board", UserID: "user-1"}

		repo.
			On("CreateBoard", mock.Anything, "user-1", req.Name).
			Return(board, nil)

		resp, err := svc.CreateBoard("user-1", req)

		assert.NoError(t, err)
		assert.Equal(t, board.ID, resp.Board.ID)
		assert.Equal(t, board.Name, resp.Board.Name)
		assert.Equal(t, board.UserID, resp.Board.UserID)
		repo.AssertExpectations(t)
	})

	t.Run("repo fails", func(t *testing.T) {
		repo := new(mockRepository)
		svc := newService(repo)

		req := &CreateBoardRequest{Name: "My Board"}

		repo.
			On("CreateBoard", mock.Anything, "user-1", req.Name).
			Return(sqlc.Board{}, errors.New("db error"))

		resp, err := svc.CreateBoard("user-1", req)

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

		req := &UpdateBoardRequest{Name: "Renamed"}
		board := sqlc.Board{ID: "board-1", Name: "Renamed", UserID: "user-1"}

		repo.
			On("GetUserBoardByID", mock.Anything, "board-1", "user-1").
			Return(board, nil)

		repo.
			On("UpdateBoardByID", mock.Anything, "board-1", req.Name).
			Return(board, nil)

		resp, err := svc.UpdateBoard("board-1", "user-1", req)

		assert.NoError(t, err)
		assert.Equal(t, board.ID, resp.Board.ID)
		assert.Equal(t, board.Name, resp.Board.Name)
		assert.Equal(t, board.UserID, resp.Board.UserID)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(mockRepository)
		svc := newService(repo)

		req := &UpdateBoardRequest{Name: "Renamed"}

		repo.
			On("GetUserBoardByID", mock.Anything, "board-1", "user-1").
			Return(sqlc.Board{}, errors.New("not found"))

		resp, err := svc.UpdateBoard("board-1", "user-1", req)

		assert.Error(t, err)
		assert.Empty(t, resp.Board.ID)
		assert.ErrorContains(t, err, "board not found")
		repo.AssertExpectations(t)
	})

	t.Run("update fails", func(t *testing.T) {
		repo := new(mockRepository)
		svc := newService(repo)

		req := &UpdateBoardRequest{Name: "Renamed"}
		board := sqlc.Board{ID: "board-1", Name: "My Board", UserID: "user-1"}

		repo.
			On("GetUserBoardByID", mock.Anything, "board-1", "user-1").
			Return(board, nil)

		repo.
			On("UpdateBoardByID", mock.Anything, "board-1", req.Name).
			Return(sqlc.Board{}, errors.New("db error"))

		resp, err := svc.UpdateBoard("board-1", "user-1", req)

		assert.Error(t, err)
		assert.Empty(t, resp.Board.ID)
		assert.ErrorContains(t, err, "failed to update board")
		repo.AssertExpectations(t)
	})
}

// --- Delete ---

func TestDeleteBoard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mockRepository)
		svc := newService(repo)

		board := sqlc.Board{ID: "board-1", Name: "My Board", UserID: "user-1"}

		repo.
			On("GetUserBoardByID", mock.Anything, "board-1", "user-1").
			Return(board, nil)

		repo.
			On("DeleteBoardByID", mock.Anything, "board-1").
			Return(nil)

		err := svc.DeleteBoard("board-1", "user-1")

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(mockRepository)
		svc := newService(repo)

		repo.
			On("GetUserBoardByID", mock.Anything, "board-1", "user-2").
			Return(sqlc.Board{}, errors.New("not found"))

		err := svc.DeleteBoard("board-1", "user-2")

		assert.Error(t, err)
		assert.ErrorContains(t, err, "board not found")
		repo.AssertExpectations(t)
	})

	t.Run("delete fails", func(t *testing.T) {
		repo := new(mockRepository)
		svc := newService(repo)

		board := sqlc.Board{ID: "board-1", Name: "My Board", UserID: "user-1"}

		repo.
			On("GetUserBoardByID", mock.Anything, "board-1", "user-1").
			Return(board, nil)

		repo.
			On("DeleteBoardByID", mock.Anything, "board-1").
			Return(errors.New("db error"))

		err := svc.DeleteBoard("board-1", "user-1")

		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to delete board")
		repo.AssertExpectations(t)
	})
}
