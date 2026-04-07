package boards

import (
	"context"

	"github.com/schodevio/trellgo/internal/platform/apierrors"
)

type Service interface {
	CreateBoard(userID string, req *CreateBoardRequest) (SingleBoardResponse, error)
	ListBoards(userID string) (ListBoardsResponse, error)
	GetBoard(boardID, userID string) (SingleBoardResponse, error)
	UpdateBoard(boardID, userID string, req *UpdateBoardRequest) (SingleBoardResponse, error)
	DeleteBoard(boardID, userID string) error
	IsOwner(boardID, userID string) (bool, error)
}

type service struct {
	repo Repository
}

func newService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateBoard(userID string, req *CreateBoardRequest) (SingleBoardResponse, error) {
	board, err := s.repo.CreateBoard(context.Background(), userID, req.Name)
	if err != nil {
		return SingleBoardResponse{}, apierrors.Internal(BOARD_CREATE_FAILED)
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
		return ListBoardsResponse{}, apierrors.Internal(BOARDS_FETCH_FAILED)
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

func (s *service) GetBoard(boardID, userID string) (SingleBoardResponse, error) {
	board, err := s.repo.GetUserBoardByID(context.Background(), boardID, userID)
	if err != nil {
		return SingleBoardResponse{}, apierrors.NotFound(BOARD_NOT_FOUND)
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

func (s *service) UpdateBoard(boardID, userID string, req *UpdateBoardRequest) (SingleBoardResponse, error) {
	board, err := s.repo.GetUserBoardByID(context.Background(), boardID, userID)
	if err != nil {
		return SingleBoardResponse{}, apierrors.NotFound(BOARD_NOT_FOUND)
	}

	board, err = s.repo.UpdateBoardByID(context.Background(), boardID, req.Name)
	if err != nil {
		return SingleBoardResponse{}, apierrors.Internal(BOARD_UPDATE_FAILED)
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

func (s *service) DeleteBoard(boardID, userID string) error {
	if _, err := s.repo.GetUserBoardByID(context.Background(), boardID, userID); err != nil {
		return apierrors.NotFound(BOARD_NOT_FOUND)
	}

	if err := s.repo.DeleteBoardByID(context.Background(), boardID); err != nil {
		return apierrors.Internal(BOARD_DELETE_FAILED)
	}

	return nil
}

func (s *service) IsOwner(boardID, userID string) (bool, error) {
	_, err := s.repo.GetUserBoardByID(context.Background(), boardID, userID)
	if err != nil {
		return false, apierrors.NotFound(BOARD_NOT_FOUND)
	}

	return true, nil
}
