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

func (m *mockRepository) GetBoardsByUserID(ctx context.Context, userID string) ([]sqlc.Board, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]sqlc.Board), args.Error(1)
}

// --- tests ---

func TestCreateBoard_Success(t *testing.T) {
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
}

func TestCreateBoard_RepoFails(t *testing.T) {
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
}

// --- List ---

func TestListBoards_Success(t *testing.T) {
	repo := new(mockRepository)
	svc := newService(repo)

	boards := []sqlc.Board{
		{ID: "board-1", Name: "First", UserID: "user-1"},
		{ID: "board-2", Name: "Second", UserID: "user-1"},
	}

	repo.
		On("GetBoardsByUserID", mock.Anything, "user-1").
		Return(boards, nil)

	resp, err := svc.ListBoards("user-1")

	assert.NoError(t, err)
	assert.Len(t, resp.Boards, 2)
	assert.Equal(t, "board-1", resp.Boards[0].ID)
	assert.Equal(t, "board-2", resp.Boards[1].ID)
	repo.AssertExpectations(t)
}

func TestListBoards_Empty(t *testing.T) {
	repo := new(mockRepository)
	svc := newService(repo)

	repo.
		On("GetBoardsByUserID", mock.Anything, "user-1").
		Return([]sqlc.Board{}, nil)

	resp, err := svc.ListBoards("user-1")

	assert.NoError(t, err)
	assert.Empty(t, resp.Boards)
	repo.AssertExpectations(t)
}

func TestListBoards_RepoFails(t *testing.T) {
	repo := new(mockRepository)
	svc := newService(repo)

	repo.
		On("GetBoardsByUserID", mock.Anything, "user-1").
		Return([]sqlc.Board{}, errors.New("db error"))

	resp, err := svc.ListBoards("user-1")

	assert.Error(t, err)
	assert.Empty(t, resp.Boards)
	assert.ErrorContains(t, err, "failed to fetch boards")
	repo.AssertExpectations(t)
}
