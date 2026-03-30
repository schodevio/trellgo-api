package boards

import (
	"context"

	"github.com/schodevio/trellgo/internal/platform/apierrors"
)

type Service interface {
	CreateBoard(req *CreateBoardRequest) (CreateBoardResponse, error)
}

type service struct {
	repo Repository
}

func newService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateBoard(req *CreateBoardRequest) (CreateBoardResponse, error) {
	board, err := s.repo.CreateBoard(context.Background(), req.Name, req.UserID)
	if err != nil {
		return CreateBoardResponse{}, apierrors.Internal("failed to create board")
	}

	return CreateBoardResponse{
		ID:        board.ID,
		Name:      board.Name,
		UserID:    board.UserID,
		CreatedAt: board.CreatedAt.Time.String(),
		UpdatedAt: board.UpdatedAt.Time.String(),
	}, nil
}
