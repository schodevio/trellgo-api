package lists

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

func (m *mockRepository) CreateList(ctx context.Context, name, boardID string, position int32) (sqlc.List, error) {
	args := m.Called(ctx, name, boardID, position)
	return args.Get(0).(sqlc.List), args.Error(1)
}

func (m *mockRepository) GetBoardLists(ctx context.Context, boardID string) ([]sqlc.List, error) {
	args := m.Called(ctx, boardID)
	return args.Get(0).([]sqlc.List), args.Error(1)
}

type mockBoardGuard struct {
	mock.Mock
}

func (m *mockBoardGuard) IsOwner(boardID, userID string) error {
	args := m.Called(boardID, userID)
	return args.Error(0)
}

// --- tests ---

func TestCreateList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mockRepository)
		boardGuard := new(mockBoardGuard)
		svc := newService(repo, boardGuard)

		req := &CreateListRequest{Name: "My List", Position: 1}
		list := sqlc.List{ID: "list-1", Name: "My List", BoardID: "board-1", Position: 1}

		boardGuard.
			On("IsOwner", "board-1", "user-1").
			Return(nil)

		repo.
			On("CreateList", mock.Anything, req.Name, "board-1", req.Position).
			Return(list, nil)

		resp, err := svc.CreateList("board-1", "user-1", req)

		assert.NoError(t, err)
		assert.Equal(t, list.ID, resp.List.ID)
		assert.Equal(t, list.Name, resp.List.Name)
		assert.Equal(t, list.BoardID, resp.List.BoardID)
		assert.Equal(t, list.Position, resp.List.Position)
		repo.AssertExpectations(t)
	})

	t.Run("not owner", func(t *testing.T) {
		repo := new(mockRepository)
		boardGuard := new(mockBoardGuard)
		svc := newService(repo, boardGuard)

		req := &CreateListRequest{Name: "My List", Position: 1}

		boardGuard.
			On("IsOwner", "board-1", "user-1").
			Return(errors.New("board not found"))

		resp, err := svc.CreateList("board-1", "user-1", req)

		assert.Error(t, err)
		assert.Empty(t, resp.List.ID)
		assert.ErrorContains(t, err, "board not found")
		boardGuard.AssertExpectations(t)
	})

	t.Run("repo fails", func(t *testing.T) {
		repo := new(mockRepository)
		boardGuard := new(mockBoardGuard)
		svc := newService(repo, boardGuard)

		req := &CreateListRequest{Name: "My List", Position: 1}

		boardGuard.
			On("IsOwner", "board-1", "user-1").
			Return(nil)

		repo.
			On("CreateList", mock.Anything, req.Name, "board-1", req.Position).
			Return(sqlc.List{}, errors.New("db error"))

		resp, err := svc.CreateList("board-1", "user-1", req)

		assert.Error(t, err)
		assert.Empty(t, resp.List.ID)
		assert.ErrorContains(t, err, "failed to create list")
		repo.AssertExpectations(t)
	})
}

func TestListLists(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mockRepository)
		boardGuard := new(mockBoardGuard)
		svc := newService(repo, boardGuard)

		lists := []sqlc.List{
			{ID: "list-1", Name: "First", BoardID: "board-1", Position: 1},
			{ID: "list-2", Name: "Second", BoardID: "board-1", Position: 2},
		}

		boardGuard.
			On("IsOwner", "board-1", "user-1").
			Return(nil)

		repo.
			On("GetBoardLists", mock.Anything, "board-1").
			Return(lists, nil)

		resp, err := svc.ListLists("board-1", "user-1")

		assert.NoError(t, err)
		assert.Len(t, resp.Lists, 2)
		assert.Equal(t, "list-1", resp.Lists[0].ID)
		assert.Equal(t, "list-2", resp.Lists[1].ID)
		repo.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		repo := new(mockRepository)
		boardGuard := new(mockBoardGuard)
		svc := newService(repo, boardGuard)

		boardGuard.
			On("IsOwner", "board-1", "user-1").
			Return(nil)

		repo.
			On("GetBoardLists", mock.Anything, "board-1").
			Return([]sqlc.List{}, nil)

		resp, err := svc.ListLists("board-1", "user-1")

		assert.NoError(t, err)
		assert.Empty(t, resp.Lists)
		repo.AssertExpectations(t)
	})

	t.Run("not owner", func(t *testing.T) {
		repo := new(mockRepository)
		boardGuard := new(mockBoardGuard)
		svc := newService(repo, boardGuard)

		boardGuard.
			On("IsOwner", "board-1", "user-1").
			Return(errors.New("board not found"))

		resp, err := svc.ListLists("board-1", "user-1")

		assert.Error(t, err)
		assert.Empty(t, resp.Lists)
		assert.ErrorContains(t, err, "board not found")
		boardGuard.AssertExpectations(t)
	})

	t.Run("repo fails", func(t *testing.T) {
		repo := new(mockRepository)
		boardGuard := new(mockBoardGuard)
		svc := newService(repo, boardGuard)

		boardGuard.
			On("IsOwner", "board-1", "user-1").
			Return(nil)

		repo.
			On("GetBoardLists", mock.Anything, "board-1").
			Return([]sqlc.List{}, errors.New("db error"))

		resp, err := svc.ListLists("board-1", "user-1")

		assert.Error(t, err)
		assert.Empty(t, resp.Lists)
		assert.ErrorContains(t, err, "failed to fetch lists")
		repo.AssertExpectations(t)
	})
}
