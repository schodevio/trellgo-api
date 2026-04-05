package cards

import (
	"context"

	"github.com/schodevio/trellgo/internal/platform/apierrors"
)

type Service interface {
	CreateCard(listID, userID string, req *CreateCardRequest) (SingleCardResponse, error)
	ListCards(listID, userID string) (ListCardsResponse, error)
}

type listChecker interface {
	IsOwner(listID, userID string) (bool, error)
}

type service struct {
	repo        Repository
	listChecker listChecker
}

func newService(repo Repository, listChecker listChecker) Service {
	return &service{repo: repo, listChecker: listChecker}
}

func (s *service) CreateCard(listID, userID string, req *CreateCardRequest) (SingleCardResponse, error) {
	if _, err := s.listChecker.IsOwner(listID, userID); err != nil {
		return SingleCardResponse{}, err
	}

	card, err := s.repo.CreateCard(context.Background(), listID, req.Title, req.Description, req.Status, req.Position)
	if err != nil {
		return SingleCardResponse{}, apierrors.Internal("failed to create card")
	}

	return SingleCardResponse{
		Card: CardResponse{
			ID:          card.ID,
			ListID:      card.ListID,
			Title:       card.Title,
			Description: card.Description.String,
			Status:      card.Status,
			Position:    card.Position,
			CreatedAt:   card.CreatedAt.Time.String(),
			UpdatedAt:   card.UpdatedAt.Time.String(),
		},
	}, nil
}

func (s *service) ListCards(listID, userID string) (ListCardsResponse, error) {
	if _, err := s.listChecker.IsOwner(listID, userID); err != nil {
		return ListCardsResponse{}, err
	}

	cards, err := s.repo.GetListCards(context.Background(), listID)
	if err != nil {
		return ListCardsResponse{}, apierrors.Internal("failed to list cards")
	}

	resp := ListCardsResponse{Cards: make([]CardResponse, len(cards))}
	for i, card := range cards {
		resp.Cards[i] = CardResponse{
			ID:          card.ID,
			ListID:      card.ListID,
			Title:       card.Title,
			Description: card.Description.String,
			Status:      card.Status,
			Position:    card.Position,
			CreatedAt:   card.CreatedAt.Time.String(),
			UpdatedAt:   card.UpdatedAt.Time.String(),
		}
	}

	return resp, nil
}
