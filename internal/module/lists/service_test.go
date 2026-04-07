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

func (m *mockRepository) CreateList(ctx context.Context, boardID, name string, position int32) (sqlc.List, error) {
	args := m.Called(ctx, boardID, name, position)
	return args.Get(0).(sqlc.List), args.Error(1)
}

func (m *mockRepository) GetBoardLists(ctx context.Context, boardID string) ([]sqlc.List, error) {
	args := m.Called(ctx, boardID)
	return args.Get(0).([]sqlc.List), args.Error(1)
}

func (m *mockRepository) GetUserListByID(ctx context.Context, listID, userID string) (sqlc.List, error) {
	args := m.Called(ctx, listID, userID)
	return args.Get(0).(sqlc.List), args.Error(1)
}

func (m *mockRepository) GetListByID(ctx context.Context, id string) (sqlc.List, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sqlc.List), args.Error(1)
}

func (m *mockRepository) UpdateListByID(ctx context.Context, id, name string) (sqlc.List, error) {
	args := m.Called(ctx, id, name)
	return args.Get(0).(sqlc.List), args.Error(1)
}

func (m *mockRepository) MoveListByID(ctx context.Context, id string, position int32) (sqlc.List, error) {
	args := m.Called(ctx, id, position)
	return args.Get(0).(sqlc.List), args.Error(1)
}

func (m *mockRepository) ReorderListsInBoard(ctx context.Context, boardID string) error {
	args := m.Called(ctx, boardID)
	return args.Error(0)
}

func (m *mockRepository) DeleteListByID(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockBoardChecker struct {
	mock.Mock
}

func (m *mockBoardChecker) IsOwner(boardID, userID string) (bool, error) {
	args := m.Called(boardID, userID)
	return args.Bool(0), args.Error(1)
}

// --- tests ---

func TestCreateList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mockRepository)
		boardChecker := new(mockBoardChecker)
		svc := newService(repo, boardChecker)

		req := &CreateListRequest{Name: "My List", Position: 1}
		list := sqlc.List{ID: "list-1", Name: "My List", BoardID: "board-1", Position: 1}

		boardChecker.
			On("IsOwner", "board-1", "user-1").
			Return(true, nil)

		repo.
			On("CreateList", mock.Anything, "board-1", req.Name, req.Position).
			Return(list, nil)

		repo.
			On("ReorderListsInBoard", mock.Anything, "board-1").
			Return(nil)

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
		boardChecker := new(mockBoardChecker)
		svc := newService(repo, boardChecker)

		req := &CreateListRequest{Name: "My List", Position: 1}

		boardChecker.
			On("IsOwner", "board-1", "user-1").
			Return(false, errors.New("board not found"))

		resp, err := svc.CreateList("board-1", "user-1", req)

		assert.Error(t, err)
		assert.Empty(t, resp.List.ID)
		assert.ErrorContains(t, err, "board not found")
		boardChecker.AssertExpectations(t)
	})

	t.Run("repo fails on create", func(t *testing.T) {
		repo := new(mockRepository)
		boardChecker := new(mockBoardChecker)
		svc := newService(repo, boardChecker)

		req := &CreateListRequest{Name: "My List", Position: 1}

		boardChecker.
			On("IsOwner", "board-1", "user-1").
			Return(true, nil)

		repo.
			On("CreateList", mock.Anything, "board-1", req.Name, req.Position).
			Return(sqlc.List{}, errors.New("db error"))

		resp, err := svc.CreateList("board-1", "user-1", req)

		assert.Error(t, err)
		assert.Empty(t, resp.List.ID)
		assert.ErrorContains(t, err, "failed to create list")
		repo.AssertExpectations(t)
	})

	t.Run("repo fails on reorder", func(t *testing.T) {
		repo := new(mockRepository)
		boardChecker := new(mockBoardChecker)
		svc := newService(repo, boardChecker)

		req := &CreateListRequest{Name: "My List", Position: 1}
		list := sqlc.List{ID: "list-1", Name: "My List", BoardID: "board-1", Position: 1}

		boardChecker.
			On("IsOwner", "board-1", "user-1").
			Return(true, nil)

		repo.
			On("CreateList", mock.Anything, "board-1", req.Name, req.Position).
			Return(list, nil)

		repo.
			On("ReorderListsInBoard", mock.Anything, "board-1").
			Return(assert.AnError)

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
		boardChecker := new(mockBoardChecker)
		svc := newService(repo, boardChecker)

		lists := []sqlc.List{
			{ID: "list-1", Name: "First", BoardID: "board-1", Position: 1},
			{ID: "list-2", Name: "Second", BoardID: "board-1", Position: 2},
		}

		boardChecker.
			On("IsOwner", "board-1", "user-1").
			Return(true, nil)

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
		boardChecker := new(mockBoardChecker)
		svc := newService(repo, boardChecker)

		boardChecker.
			On("IsOwner", "board-1", "user-1").
			Return(true, nil)

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
		boardChecker := new(mockBoardChecker)
		svc := newService(repo, boardChecker)

		boardChecker.
			On("IsOwner", "board-1", "user-1").
			Return(false, errors.New("board not found"))

		resp, err := svc.ListLists("board-1", "user-1")

		assert.Error(t, err)
		assert.Empty(t, resp.Lists)
		assert.ErrorContains(t, err, "board not found")
		boardChecker.AssertExpectations(t)
	})

	t.Run("repo fails", func(t *testing.T) {
		repo := new(mockRepository)
		boardChecker := new(mockBoardChecker)
		svc := newService(repo, boardChecker)

		boardChecker.
			On("IsOwner", "board-1", "user-1").
			Return(true, nil)

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

func TestUpdateList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mockRepository)
		boardChecker := new(mockBoardChecker)
		svc := newService(repo, boardChecker)

		req := &UpdateListRequest{Name: "Updated"}
		existing := sqlc.List{ID: "list-1", Name: "Old", BoardID: "board-1", Position: 1}
		updated := sqlc.List{ID: "list-1", Name: "Updated", BoardID: "board-1", Position: 1}

		repo.
			On("GetListByID", mock.Anything, "list-1").
			Return(existing, nil)

		boardChecker.
			On("IsOwner", "board-1", "user-1").
			Return(true, nil)

		repo.
			On("UpdateListByID", mock.Anything, "list-1", req.Name).
			Return(updated, nil)

		resp, err := svc.UpdateList("list-1", "user-1", req)

		assert.NoError(t, err)
		assert.Equal(t, updated.ID, resp.List.ID)
		assert.Equal(t, updated.Name, resp.List.Name)
		assert.Equal(t, updated.Position, resp.List.Position)
		repo.AssertExpectations(t)
		boardChecker.AssertExpectations(t)
	})

	t.Run("list not found", func(t *testing.T) {
		repo := new(mockRepository)
		boardChecker := new(mockBoardChecker)
		svc := newService(repo, boardChecker)

		req := &UpdateListRequest{Name: "Updated"}

		repo.
			On("GetListByID", mock.Anything, "list-1").
			Return(sqlc.List{}, errors.New("not found"))

		resp, err := svc.UpdateList("list-1", "user-1", req)

		assert.Error(t, err)
		assert.Empty(t, resp.List.ID)
		assert.ErrorContains(t, err, "list not found")
		repo.AssertExpectations(t)
	})

	t.Run("not owner", func(t *testing.T) {
		repo := new(mockRepository)
		boardChecker := new(mockBoardChecker)
		svc := newService(repo, boardChecker)

		req := &UpdateListRequest{Name: "Updated"}
		existing := sqlc.List{ID: "list-1", Name: "Old", BoardID: "board-1", Position: 1}

		repo.
			On("GetListByID", mock.Anything, "list-1").
			Return(existing, nil)

		boardChecker.
			On("IsOwner", "board-1", "user-1").
			Return(false, errors.New("board not found"))

		resp, err := svc.UpdateList("list-1", "user-1", req)

		assert.Error(t, err)
		assert.Empty(t, resp.List.ID)
		assert.ErrorContains(t, err, "board not found")
		repo.AssertExpectations(t)
		boardChecker.AssertExpectations(t)
	})

	t.Run("repo fails", func(t *testing.T) {
		repo := new(mockRepository)
		boardChecker := new(mockBoardChecker)
		svc := newService(repo, boardChecker)

		req := &UpdateListRequest{Name: "Updated"}
		existing := sqlc.List{ID: "list-1", Name: "Old", BoardID: "board-1", Position: 1}

		repo.
			On("GetListByID", mock.Anything, "list-1").
			Return(existing, nil)

		boardChecker.
			On("IsOwner", "board-1", "user-1").
			Return(true, nil)

		repo.
			On("UpdateListByID", mock.Anything, "list-1", req.Name).
			Return(sqlc.List{}, errors.New("db error"))

		resp, err := svc.UpdateList("list-1", "user-1", req)

		assert.Error(t, err)
		assert.Empty(t, resp.List.ID)
		assert.ErrorContains(t, err, "failed to update list")
		repo.AssertExpectations(t)
		boardChecker.AssertExpectations(t)
	})
}

func TestMoveList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mockRepository)
		boardChecker := new(mockBoardChecker)
		svc := newService(repo, boardChecker)

		req := &MoveListRequest{Position: 2}
		existing := sqlc.List{ID: "list-1", Name: "My List", BoardID: "board-1", Position: 1}
		moved := sqlc.List{ID: "list-1", Name: "My List", BoardID: "board-1", Position: 2}

		repo.
			On("GetListByID", mock.Anything, "list-1").
			Return(existing, nil)

		boardChecker.
			On("IsOwner", "board-1", "user-1").
			Return(true, nil)

		repo.
			On("MoveListByID", mock.Anything, "list-1", req.Position).
			Return(moved, nil)

		repo.
			On("ReorderListsInBoard", mock.Anything, "board-1").
			Return(nil)

		resp, err := svc.MoveList("list-1", "user-1", req)

		assert.NoError(t, err)
		assert.Equal(t, moved.ID, resp.List.ID)
		assert.Equal(t, moved.Position, resp.List.Position)
		repo.AssertExpectations(t)
		boardChecker.AssertExpectations(t)
	})

	t.Run("list not found", func(t *testing.T) {
		repo := new(mockRepository)
		boardChecker := new(mockBoardChecker)
		svc := newService(repo, boardChecker)

		req := &MoveListRequest{Position: 2}

		repo.
			On("GetListByID", mock.Anything, "list-1").
			Return(sqlc.List{}, assert.AnError)

		resp, err := svc.MoveList("list-1", "user-1", req)

		assert.Error(t, err)
		assert.Empty(t, resp.List.ID)
		assert.ErrorContains(t, err, "list not found")
		repo.AssertExpectations(t)
	})

	t.Run("not owner", func(t *testing.T) {
		repo := new(mockRepository)
		boardChecker := new(mockBoardChecker)
		svc := newService(repo, boardChecker)

		req := &MoveListRequest{Position: 2}
		existing := sqlc.List{ID: "list-1", Name: "My List", BoardID: "board-1", Position: 1}

		repo.
			On("GetListByID", mock.Anything, "list-1").
			Return(existing, nil)

		boardChecker.
			On("IsOwner", "board-1", "user-1").
			Return(false, errors.New("board not found"))

		resp, err := svc.MoveList("list-1", "user-1", req)

		assert.Error(t, err)
		assert.Empty(t, resp.List.ID)
		assert.ErrorContains(t, err, "board not found")
		repo.AssertExpectations(t)
		boardChecker.AssertExpectations(t)
	})

	t.Run("repo fails on move", func(t *testing.T) {
		repo := new(mockRepository)
		boardChecker := new(mockBoardChecker)
		svc := newService(repo, boardChecker)

		req := &MoveListRequest{Position: 2}
		existing := sqlc.List{ID: "list-1", Name: "My List", BoardID: "board-1", Position: 1}

		repo.
			On("GetListByID", mock.Anything, "list-1").
			Return(existing, nil)

		boardChecker.
			On("IsOwner", "board-1", "user-1").
			Return(true, nil)

		repo.
			On("MoveListByID", mock.Anything, "list-1", req.Position).
			Return(sqlc.List{}, errors.New("db error"))

		resp, err := svc.MoveList("list-1", "user-1", req)

		assert.Error(t, err)
		assert.Empty(t, resp.List.ID)
		assert.ErrorContains(t, err, "failed to move list")
		repo.AssertExpectations(t)
		boardChecker.AssertExpectations(t)
	})

	t.Run("repo fails on reorder", func(t *testing.T) {
		repo := new(mockRepository)
		boardChecker := new(mockBoardChecker)
		svc := newService(repo, boardChecker)

		req := &MoveListRequest{Position: 2}
		existing := sqlc.List{ID: "list-1", Name: "My List", BoardID: "board-1", Position: 1}
		moved := sqlc.List{ID: "list-1", Name: "My List", BoardID: "board-1", Position: 2}

		repo.
			On("GetListByID", mock.Anything, "list-1").
			Return(existing, nil)

		boardChecker.
			On("IsOwner", "board-1", "user-1").
			Return(true, nil)

		repo.
			On("MoveListByID", mock.Anything, "list-1", req.Position).
			Return(moved, nil)

		repo.
			On("ReorderListsInBoard", mock.Anything, "board-1").
			Return(assert.AnError)

		resp, err := svc.MoveList("list-1", "user-1", req)

		assert.Error(t, err)
		assert.Empty(t, resp.List.ID)
		assert.ErrorContains(t, err, "failed to move list")
		repo.AssertExpectations(t)
		boardChecker.AssertExpectations(t)
	})
}

func TestDeleteList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mockRepository)
		boardChecker := new(mockBoardChecker)
		svc := newService(repo, boardChecker)

		existing := sqlc.List{ID: "list-1", BoardID: "board-1"}

		repo.
			On("GetListByID", mock.Anything, "list-1").
			Return(existing, nil)

		boardChecker.
			On("IsOwner", "board-1", "user-1").
			Return(true, nil)

		repo.
			On("DeleteListByID", mock.Anything, "list-1").
			Return(nil)

		err := svc.DeleteList("list-1", "user-1")

		assert.NoError(t, err)
		repo.AssertExpectations(t)
		boardChecker.AssertExpectations(t)
	})

	t.Run("list not found", func(t *testing.T) {
		repo := new(mockRepository)
		boardChecker := new(mockBoardChecker)
		svc := newService(repo, boardChecker)

		repo.
			On("GetListByID", mock.Anything, "list-1").
			Return(sqlc.List{}, errors.New("not found"))

		err := svc.DeleteList("list-1", "user-1")

		assert.Error(t, err)
		assert.ErrorContains(t, err, "list not found")
		repo.AssertExpectations(t)
	})

	t.Run("not owner", func(t *testing.T) {
		repo := new(mockRepository)
		boardChecker := new(mockBoardChecker)
		svc := newService(repo, boardChecker)

		existing := sqlc.List{ID: "list-1", BoardID: "board-1"}

		repo.
			On("GetListByID", mock.Anything, "list-1").
			Return(existing, nil)

		boardChecker.
			On("IsOwner", "board-1", "user-1").
			Return(false, errors.New("board not found"))

		err := svc.DeleteList("list-1", "user-1")

		assert.Error(t, err)
		assert.ErrorContains(t, err, "board not found")
		repo.AssertExpectations(t)
		boardChecker.AssertExpectations(t)
	})

	t.Run("repo fails", func(t *testing.T) {
		repo := new(mockRepository)
		boardChecker := new(mockBoardChecker)
		svc := newService(repo, boardChecker)

		existing := sqlc.List{ID: "list-1", BoardID: "board-1"}

		repo.
			On("GetListByID", mock.Anything, "list-1").
			Return(existing, nil)

		boardChecker.
			On("IsOwner", "board-1", "user-1").
			Return(true, nil)

		repo.
			On("DeleteListByID", mock.Anything, "list-1").
			Return(errors.New("db error"))

		err := svc.DeleteList("list-1", "user-1")

		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to delete list")
		repo.AssertExpectations(t)
		boardChecker.AssertExpectations(t)
	})
}

func TestIsOwner(t *testing.T) {
	t.Run("is owner", func(t *testing.T) {
		repo := new(mockRepository)
		boardChecker := new(mockBoardChecker)
		svc := newService(repo, boardChecker)

		repo.
			On("GetUserListByID", mock.Anything, "list-1", "user-1").
			Return(sqlc.List{ID: "list-1"}, nil)

		isOwner, err := svc.IsOwner("list-1", "user-1")

		assert.NoError(t, err)
		assert.True(t, isOwner)
		repo.AssertExpectations(t)
	})

	t.Run("not owner", func(t *testing.T) {
		repo := new(mockRepository)
		boardChecker := new(mockBoardChecker)
		svc := newService(repo, boardChecker)

		repo.
			On("GetUserListByID", mock.Anything, "list-1", "user-1").
			Return(sqlc.List{}, errors.New("not found"))

		isOwner, err := svc.IsOwner("list-1", "user-1")

		assert.Error(t, err)
		assert.False(t, isOwner)
		assert.ErrorContains(t, err, "list not found")
		repo.AssertExpectations(t)
	})

	t.Run("repo fails", func(t *testing.T) {
		repo := new(mockRepository)
		boardChecker := new(mockBoardChecker)
		svc := newService(repo, boardChecker)

		repo.
			On("GetUserListByID", mock.Anything, "list-1", "user-1").
			Return(sqlc.List{}, errors.New("db error"))

		isOwner, err := svc.IsOwner("list-1", "user-1")

		assert.Error(t, err)
		assert.False(t, isOwner)
		assert.ErrorContains(t, err, "list not found")
		repo.AssertExpectations(t)
	})
}
