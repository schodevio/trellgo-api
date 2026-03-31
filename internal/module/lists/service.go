package lists

import (
	"context"

	"github.com/schodevio/trellgo/internal/platform/apierrors"
)

type Service interface {
	CreateList(boardID, userID string, req *CreateListRequest) (SingleListResponse, error)
}

type boardGuard interface {
	IsOwner(boardID, userID string) error
}

type service struct {
	repo       Repository
	boardGuard boardGuard
}

func newService(repo Repository, boardGuard boardGuard) Service {
	return &service{repo: repo, boardGuard: boardGuard}
}

func (s *service) CreateList(boardID, userID string, req *CreateListRequest) (SingleListResponse, error) {
	if err := s.boardGuard.IsOwner(boardID, userID); err != nil {
		return SingleListResponse{}, err
	}

	list, err := s.repo.CreateList(context.Background(), req.Name, boardID, req.Position)
	if err != nil {
		return SingleListResponse{}, apierrors.Internal("failed to create list")
	}

	return SingleListResponse{
		List: ListResponse{
			ID:        list.ID,
			Name:      list.Name,
			BoardID:   list.BoardID,
			Position:  list.Position,
			CreatedAt: list.CreatedAt.Time.String(),
			UpdatedAt: list.UpdatedAt.Time.String(),
		},
	}, nil
}
