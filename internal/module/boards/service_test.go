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
	assert.Equal(t, board.ID, resp.ID)
	assert.Equal(t, board.Name, resp.Name)
	assert.Equal(t, board.UserID, resp.UserID)
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
	assert.Empty(t, resp.ID)
	assert.ErrorContains(t, err, "failed to create board")
	repo.AssertExpectations(t)
}
