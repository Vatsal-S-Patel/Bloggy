package dto

import (
	"time"

	"github.com/google/uuid"
)

type AddBookmarkRequest struct {
	Name    string `json:"name" validate:"required,max=30"`
	Visible bool   `json:"visible"`
}

type BookmarkBlogs struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	FtImage   string    `json:"ft_image" db:"ft_image"`
	AuthorID  uuid.UUID `json:"author_id" db:"author_id"`
	Author    string    `json:"author" db:"author"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
