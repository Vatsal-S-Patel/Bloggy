package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Body      string    `json:"body" db:"body"`
	Claps     int       `json:"claps" db:"claps"`
	Replies   int       `json:"replies" db:"replies"`
	AuthorID  uuid.UUID `json:"author_id" db:"author_id"`
	BlogID    uuid.UUID `json:"blog_id" db:"blog_id"`
	ParentID  uuid.UUID `json:"parent_id" db:"parent_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
