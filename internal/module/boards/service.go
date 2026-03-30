package boards

import (
	"context"

	"github.com/schodevio/trellgo/internal/platform/apierrors"
)

type Service interface {
	CreateBoard(req *CreateBoardRequest) (SingleBoardResponse, error)
	ListBoards(userID string) (ListBoardsResponse, error)
	GetBoard(id, userID string) (SingleBoardResponse, error)
}

type service struct {
	repo Repository
}

func newService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateBoard(req *CreateBoardRequest) (SingleBoardResponse, error) {
	board, err := s.repo.CreateBoard(context.Background(), req.Name, req.UserID)
	if err != nil {
		return SingleBoardResponse{}, apierrors.Internal("failed to create board")
	}

	return SingleBoardResponse{
		Board: BoardResponse{
			ID:        board.ID,
			Name:      board.Name,
			UserID:    board.UserID,
			CreatedAt: board.CreatedAt.Time.String(),
			UpdatedAt: board.UpdatedAt.Time.String(),
		},
	}, nil
}

func (s *service) ListBoards(userID string) (ListBoardsResponse, error) {
	boards, err := s.repo.GetUserBoards(context.Background(), userID)
	if err != nil {
		return ListBoardsResponse{}, apierrors.Internal("failed to fetch boards")
	}

	items := make([]BoardResponse, 0, len(boards))
	for _, b := range boards {
		items = append(items, BoardResponse{
			ID:        b.ID,
			Name:      b.Name,
			UserID:    b.UserID,
			CreatedAt: b.CreatedAt.Time.String(),
			UpdatedAt: b.UpdatedAt.Time.String(),
		})
	}

	return ListBoardsResponse{Boards: items}, nil
}

func (s *service) GetBoard(id, userID string) (SingleBoardResponse, error) {
	board, err := s.repo.GetUserBoardByID(context.Background(), id, userID)
	if err != nil {
		return SingleBoardResponse{}, apierrors.NotFound("board not found")
	}

	return SingleBoardResponse{
		Board: BoardResponse{
			ID:        board.ID,
			Name:      board.Name,
			UserID:    board.UserID,
			CreatedAt: board.CreatedAt.Time.String(),
			UpdatedAt: board.UpdatedAt.Time.String(),
		},
	}, nil
}
