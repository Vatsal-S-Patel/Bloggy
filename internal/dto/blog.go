package dto

import (
	"time"

	"github.com/Vatsal-S-Patel/Bloggy/models"
	"github.com/google/uuid"
)

type PublishBlogRequest struct {
	Title    string      `json:"title" validate:"required,max=130"`
	Subtitle string      `json:"subtitle" validate:"required,max=170"`
	Content  string      `json:"content" validate:"required"`
	FtImage  string      `json:"ft_image" validate:"omitempty,url"`
	Tags     []uuid.UUID `json:"tags,omitempty" validate:"dive,required,uuid"`
}

type Blog struct {
	ID        uuid.UUID     `json:"id" db:"id"`
	Title     string        `json:"title" db:"title"`
	Subtitle  string        `json:"subtitle" db:"subtitle"`
	Content   string        `json:"content" db:"content"`
	FtImage   string        `json:"ft_image" db:"ft_image"`
	Tags      []*models.Tag `json:"tags"`
	AuthorID  uuid.UUID     `json:"author_id" db:"author_id"`
	Author    string        `json:"author" db:"author"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" db:"updated_at"`
}
