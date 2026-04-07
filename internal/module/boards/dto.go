package boards

// Errors
const (
	BOARD_NOT_FOUND     string = "board not found"
	BOARDS_FETCH_FAILED string = "failed to fetch boards"
	BOARD_CREATE_FAILED string = "failed to create board"
	BOARD_UPDATE_FAILED string = "failed to update board"
	BOARD_DELETE_FAILED string = "failed to delete board"
)

type CreateBoardRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

type UpdateBoardRequest struct {
	Name string `json:"name"    validate:"required,min=1,max=255"`
}

type BoardResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type SingleBoardResponse struct {
	Board BoardResponse `json:"board"`
}

type ListBoardsResponse struct {
	Boards []BoardResponse `json:"boards"`
}
