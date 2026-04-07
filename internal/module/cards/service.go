package cards

import (
	"context"

	"github.com/schodevio/trellgo/internal/platform/apierrors"
)

type Service interface {
	CreateCard(listID, userID string, req *CreateCardRequest) (SingleCardResponse, error)
	ListCards(listID, userID string) (ListCardsResponse, error)
	UpdateCard(cardID, userID string, req *UpdateCardRequest) (SingleCardResponse, error)
	MoveCard(cardID, userID string, req *MoveCardRequest) (SingleCardResponse, error)
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
	ctx := context.Background()

	// Check list ownership
	if _, err := s.listChecker.IsOwner(listID, userID); err != nil {
		return SingleCardResponse{}, err
	}

	// Create the card
	card, err := s.repo.CreateCard(ctx, listID, req.Title, req.Description, req.Position)
	if err != nil {
		return SingleCardResponse{}, apierrors.Internal(CARD_CREATE_FAILED)
	}

	// Reorder cards in the list
	if err := s.repo.ReorderCardsInList(ctx, listID); err != nil {
		return SingleCardResponse{}, apierrors.Internal(CARD_UPDATE_FAILED)
	}

	return SingleCardResponse{
		Card: CardResponse{
			ID:          card.ID,
			ListID:      card.ListID,
			Title:       card.Title,
			Description: card.Description.String,
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
		return ListCardsResponse{}, apierrors.Internal(CARDS_FETCH_FAILED)
	}

	resp := ListCardsResponse{Cards: make([]CardResponse, len(cards))}
	for i, card := range cards {
		resp.Cards[i] = CardResponse{
			ID:          card.ID,
			ListID:      card.ListID,
			Title:       card.Title,
			Description: card.Description.String,
			Position:    card.Position,
			CreatedAt:   card.CreatedAt.Time.String(),
			UpdatedAt:   card.UpdatedAt.Time.String(),
		}
	}

	return resp, nil
}

func (s *service) UpdateCard(cardID, userID string, req *UpdateCardRequest) (SingleCardResponse, error) {
	ctx := context.Background()

	// Fetch the card
	card, err := s.repo.GetCardByID(ctx, cardID)
	if err != nil {
		return SingleCardResponse{}, apierrors.NotFound(CARD_NOT_FOUND)
	}

	// Check card ownership
	if _, err := s.listChecker.IsOwner(card.ListID, userID); err != nil {
		return SingleCardResponse{}, err
	}

	// Update the card
	card, err = s.repo.UpdateCardByID(ctx, cardID, req.Title, req.Description)
	if err != nil {
		return SingleCardResponse{}, apierrors.Internal(CARD_UPDATE_FAILED)
	}

	return SingleCardResponse{
		Card: CardResponse{
			ID:          card.ID,
			ListID:      card.ListID,
			Title:       card.Title,
			Description: card.Description.String,
			Position:    card.Position,
			CreatedAt:   card.CreatedAt.Time.String(),
			UpdatedAt:   card.UpdatedAt.Time.String(),
		},
	}, nil
}

func (s *service) MoveCard(cardID, userID string, req *MoveCardRequest) (SingleCardResponse, error) {
	ctx := context.Background()

	// Fetch the card
	card, err := s.repo.GetCardByID(ctx, cardID)
	if err != nil {
		return SingleCardResponse{}, apierrors.NotFound(CARD_NOT_FOUND)
	}

	// Check card ownership
	if _, err := s.listChecker.IsOwner(card.ListID, userID); err != nil {
		return SingleCardResponse{}, err
	}

	oldListID := card.ListID

	if oldListID != req.ListID {
		// If moving to a different list, check ownership of the new list as well
		if _, err := s.listChecker.IsOwner(req.ListID, userID); err != nil {
			return SingleCardResponse{}, err
		}
	}

	// Update card's list and position
	card, err = s.repo.MoveCardByID(ctx, cardID, req.ListID, req.Position)
	if err != nil {
		return SingleCardResponse{}, apierrors.Internal(CARD_UPDATE_FAILED)
	}

	// Reorder cards in the new list
	if err := s.repo.ReorderCardsInList(ctx, req.ListID); err != nil {
		return SingleCardResponse{}, apierrors.Internal(CARD_UPDATE_FAILED)
	}

	// If the card was moved to a different list, reorder cards in the old list as well
	if oldListID != req.ListID {
		if err := s.repo.ReorderCardsInList(ctx, oldListID); err != nil {
			return SingleCardResponse{}, apierrors.Internal(CARD_UPDATE_FAILED)
		}
	}

	return SingleCardResponse{
		Card: CardResponse{
			ID:          card.ID,
			ListID:      card.ListID,
			Title:       card.Title,
			Description: card.Description.String,
			Position:    card.Position,
			CreatedAt:   card.CreatedAt.Time.String(),
			UpdatedAt:   card.UpdatedAt.Time.String(),
		},
	}, nil
}

func (s *service) DeleteCard(cardID, userID string) error {
	ctx := context.Background()

	// Fetch the card
	card, err := s.repo.GetCardByID(ctx, cardID)
	if err != nil {
		return apierrors.NotFound(CARD_NOT_FOUND)
	}

	// Check card ownership
	if _, err := s.listChecker.IsOwner(card.ListID, userID); err != nil {
		return err
	}

	// Delete the card
	if err := s.repo.DeleteCardByID(ctx, cardID); err != nil {
		return apierrors.Internal(CARD_DELETE_FAILED)
	}

	return nil
}
