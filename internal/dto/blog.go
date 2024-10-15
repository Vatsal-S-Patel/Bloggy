package dto

import "github.com/google/uuid"

type PublishBlogRequest struct {
	Title    string      `json:"title" validate:"required,max=130"`
	Subtitle string      `json:"subtitle" validate:"required,max=170"`
	Content  string      `json:"content" validate:"required"`
	FtImage  string      `json:"ft_image" validate:"omitempty,url"`
	Tags     []uuid.UUID `json:"tags,omitempty" validate:"dive,required,uuid"`
}
