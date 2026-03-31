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

// --- tests ---

func TestCreateList_Success(t *testing.T) {
	repo := new(mockRepository)
	svc := newService(repo)

	req := &CreateListRequest{Name: "My List", Position: 1}
	list := sqlc.List{ID: "list-1", Name: "My List", BoardID: "board-1", Position: 1}

	repo.
		On("CreateList", mock.Anything, req.Name, "board-1", req.Position).
		Return(list, nil)

	resp, err := svc.CreateList("board-1", req)

	assert.NoError(t, err)
	assert.Equal(t, list.ID, resp.List.ID)
	assert.Equal(t, list.Name, resp.List.Name)
	assert.Equal(t, list.BoardID, resp.List.BoardID)
	assert.Equal(t, list.Position, resp.List.Position)
	repo.AssertExpectations(t)
}

func TestCreateList_RepoFails(t *testing.T) {
	repo := new(mockRepository)
	svc := newService(repo)

	req := &CreateListRequest{Name: "My List", Position: 1}

	repo.
		On("CreateList", mock.Anything, req.Name, "board-1", req.Position).
		Return(sqlc.List{}, errors.New("db error"))

	resp, err := svc.CreateList("board-1", req)

	assert.Error(t, err)
	assert.Empty(t, resp.List.ID)
	assert.ErrorContains(t, err, "failed to create list")
	repo.AssertExpectations(t)
}
