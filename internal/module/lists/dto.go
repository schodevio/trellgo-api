package lists

// Errors
const (
	LIST_NOT_FOUND     string = "list not found"
	LISTS_FETCH_FAILED string = "failed to fetch lists"
	LIST_CREATE_FAILED string = "failed to create list"
	LIST_UPDATE_FAILED string = "failed to update list"
	LIST_MOVE_FAILED   string = "failed to move list"
	LIST_DELETE_FAILED string = "failed to delete list"
)

type CreateListRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=255"`
	Position int32  `json:"position" validate:"required,min=0"`
}

type UpdateListRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

type MoveListRequest struct {
	Position int32 `json:"position" validate:"required,min=0"`
}

type ListResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	BoardID   string `json:"board_id"`
	Position  int32  `json:"position"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type SingleListResponse struct {
	List ListResponse `json:"list"`
}

type ListListsResponse struct {
	Lists []ListResponse `json:"lists"`
}
