package cards

// Errors
const (
	CARD_NOT_FOUND     string = "card not found"
	CARDS_FETCH_FAILED string = "failed to fetch cards"
	CARD_CREATE_FAILED string = "failed to create card"
	CARD_UPDATE_FAILED string = "failed to update card"
	CARD_MOVE_FAILED   string = "failed to move card"
	CARD_DELETE_FAILED string = "failed to delete card"
)

type CreateCardRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description,omitempty" validate:"max=1000"`
	Position    int32  `json:"position" validate:"required,min=0"`
}

type UpdateCardRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description,omitempty" validate:"max=1000"`
}

type MoveCardRequest struct {
	ListID   string `json:"list_id" validate:"required"`
	Position int32  `json:"position" validate:"required,min=0"`
}

type CardResponse struct {
	ID          string `json:"id"`
	ListID      string `json:"list_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Position    int32  `json:"position"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type SingleCardResponse struct {
	Card CardResponse `json:"card"`
}

type ListCardsResponse struct {
	Cards []CardResponse `json:"cards"`
}
