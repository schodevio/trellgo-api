package cards

import (
	"context"

	"github.com/schodevio/trellgo/internal/platform/apierrors"
)

type Service interface {
	CreateCard(listID, userID string, req *CreateCardRequest) (SingleCardResponse, error)
	ListCards(listID, userID string) (ListCardsResponse, error)
	UpdateCard(cardID, userID string, req *UpdateCardRequest) (SingleCardResponse, error)
	DeleteCard(cardID, userID string) error
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

func (s *service) UpdateCard(cardID, userID string, req *UpdateCardRequest) (SingleCardResponse, error) {
	card, err := s.repo.GetCardByID(context.Background(), cardID)
	if err != nil {
		return SingleCardResponse{}, apierrors.NotFound("card not found")
	}

	if _, err := s.listChecker.IsOwner(card.ListID, userID); err != nil {
		return SingleCardResponse{}, err
	}

	updatedCard, err := s.repo.UpdateCardByID(context.Background(), cardID, req.Title, req.Description, req.Status, req.Position)
	if err != nil {
		return SingleCardResponse{}, apierrors.Internal("failed to update card")
	}

	return SingleCardResponse{
		Card: CardResponse{
			ID:          updatedCard.ID,
			ListID:      updatedCard.ListID,
			Title:       updatedCard.Title,
			Description: updatedCard.Description.String,
			Status:      updatedCard.Status,
			Position:    updatedCard.Position,
			CreatedAt:   updatedCard.CreatedAt.Time.String(),
			UpdatedAt:   updatedCard.UpdatedAt.Time.String(),
		},
	}, nil
}

func (s *service) DeleteCard(cardID, userID string) error {
	card, err := s.repo.GetCardByID(context.Background(), cardID)
	if err != nil {
		return apierrors.NotFound("card not found")
	}

	if _, err := s.listChecker.IsOwner(card.ListID, userID); err != nil {
		return err
	}

	if err := s.repo.DeleteCardByID(context.Background(), cardID); err != nil {
		return apierrors.Internal("failed to delete card")
	}

	return nil
}
