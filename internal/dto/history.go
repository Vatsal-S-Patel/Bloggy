package dto

import (
	"time"

	"github.com/google/uuid"
)

type History struct {
	BlogID    uuid.UUID `json:"blog_id" db:"id"`
	Title     string    `json:"title" db:"title"`
	FtImage   string    `json:"ft_image" db:"ft_image"`
	AuthorID  uuid.UUID `json:"author_id" db:"author_id"`
	Author    string    `json:"author" db:"author"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
