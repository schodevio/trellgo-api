package lists

type CreateListRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=255"`
	Position int32  `json:"position" validate:"required,min=0"`
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
