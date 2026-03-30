package boards

type CreateBoardRequest struct {
	Name   string `json:"name" validate:"required,min=1,max=255"`
	UserID string `json:"user_id" validate:"required"`
}

type CreateBoardResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
