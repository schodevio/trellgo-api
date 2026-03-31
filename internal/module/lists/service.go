package lists

import (
	"context"

	"github.com/schodevio/trellgo/internal/platform/apierrors"
)

type Service interface {
	CreateList(boardID string, req *CreateListRequest) (SingleListResponse, error)
}

type service struct {
	repo Repository
}

func newService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateList(boardID string, req *CreateListRequest) (SingleListResponse, error) {
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
