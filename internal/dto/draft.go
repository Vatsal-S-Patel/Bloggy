package dto

type AddDraftRequest struct {
	Title    string `json:"title" validate:"required,max=130"`
	Subtitle string `json:"subtitle" validate:"required,max=170"`
	Content  string `json:"content" validate:"required"`
}
