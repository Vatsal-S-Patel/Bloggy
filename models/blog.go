package models

import (
	"time"

	"github.com/google/uuid"
)

type Blog struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Subtitle  string    `json:"subtitle" db:"subtitle"`
	Content   string    `json:"content" db:"content"`
	FtImage   string    `json:"ft_image" db:"ft_image"`
	AuthorID  uuid.UUID `json:"author_id" db:"author_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Tag struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type BlogTag struct {
	BlogID uuid.UUID `json:"blog_id"`
	TagID  uuid.UUID `json:"tag_id"`
}

type Draft struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Subtitle  string    `json:"subtitle" db:"subtitle"`
	Content   string    `json:"content" db:"content"`
	AuthorID  uuid.UUID `json:"author_id" db:"author_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
