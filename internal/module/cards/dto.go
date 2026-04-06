package cards

type CreateCardRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description,omitempty" validate:"max=1000"`
	Status      string `json:"status" validate:"required,oneof=todo in_progress done"`
	Position    int32  `json:"position" validate:"required,min=0"`
}

type UpdateCardRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description,omitempty" validate:"max=1000"`
	Status      string `json:"status" validate:"required,oneof=todo in_progress done"`
	Position    int32  `json:"position" validate:"required,min=0"`
}

type CardResponse struct {
	ID          string `json:"id"`
	ListID      string `json:"list_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
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
