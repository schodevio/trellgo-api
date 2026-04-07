package lists

import (
	"context"

	"github.com/schodevio/trellgo/internal/platform/apierrors"
)

type Service interface {
	CreateList(boardID, userID string, req *CreateListRequest) (SingleListResponse, error)
	ListLists(boardID, userID string) (ListListsResponse, error)
	UpdateList(listID, userID string, req *UpdateListRequest) (SingleListResponse, error)
	DeleteList(listID, userID string) error
	IsOwner(listID, userID string) (bool, error)
}

type boardChecker interface {
	IsOwner(boardID, userID string) (bool, error)
}

type service struct {
	repo         Repository
	boardChecker boardChecker
}

func newService(repo Repository, boardChecker boardChecker) Service {
	return &service{repo: repo, boardChecker: boardChecker}
}

func (s *service) CreateList(boardID, userID string, req *CreateListRequest) (SingleListResponse, error) {
	if _, err := s.boardChecker.IsOwner(boardID, userID); err != nil {
		return SingleListResponse{}, err
	}

	list, err := s.repo.CreateList(context.Background(), boardID, req.Name, req.Position)
	if err != nil {
		return SingleListResponse{}, apierrors.Internal(LIST_CREATE_FAILED)
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

func (s *service) ListLists(boardID, userID string) (ListListsResponse, error) {
	if _, err := s.boardChecker.IsOwner(boardID, userID); err != nil {
		return ListListsResponse{}, err
	}

	lists, err := s.repo.GetBoardLists(context.Background(), boardID)
	if err != nil {
		return ListListsResponse{}, apierrors.Internal(LISTS_FETCH_FAILED)
	}

	items := make([]ListResponse, 0, len(lists))
	for _, l := range lists {
		items = append(items, ListResponse{
			ID:        l.ID,
			Name:      l.Name,
			BoardID:   l.BoardID,
			Position:  l.Position,
			CreatedAt: l.CreatedAt.Time.String(),
			UpdatedAt: l.UpdatedAt.Time.String(),
		})
	}

	return ListListsResponse{Lists: items}, nil
}

func (s *service) UpdateList(id, userID string, req *UpdateListRequest) (SingleListResponse, error) {
	list, err := s.repo.GetListByID(context.Background(), id)
	if err != nil {
		return SingleListResponse{}, apierrors.NotFound(LIST_NOT_FOUND)
	}

	if _, err := s.boardChecker.IsOwner(list.BoardID, userID); err != nil {
		return SingleListResponse{}, err
	}

	list, err = s.repo.UpdateListByID(context.Background(), id, req.Name)
	if err != nil {
		return SingleListResponse{}, apierrors.Internal(LIST_UPDATE_FAILED)
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

func (s *service) DeleteList(listID, userID string) error {
	list, err := s.repo.GetListByID(context.Background(), listID)
	if err != nil {
		return apierrors.NotFound(LIST_NOT_FOUND)
	}

	if _, err := s.boardChecker.IsOwner(list.BoardID, userID); err != nil {
		return err
	}

	if err := s.repo.DeleteListByID(context.Background(), listID); err != nil {
		return apierrors.Internal(LIST_DELETE_FAILED)
	}

	return nil
}

func (s *service) IsOwner(listID, userID string) (bool, error) {
	_, err := s.repo.GetUserListByID(context.Background(), listID, userID)
	if err != nil {
		return false, apierrors.NotFound(LIST_NOT_FOUND)
	}

	return true, nil
}
